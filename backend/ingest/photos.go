package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"pitwall/backend/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Driver photos come from Wikipedia's page summaries (images hosted on Wikimedia
// Commons). OCB has no driver images. Each driver is looked up once; the result
// (found or not) is stored so later polls skip them.

const (
	wikiSummaryURL = "https://en.wikipedia.org/api/rest_v1/page/summary/"
	// Wikimedia only serves resized images at its standard widths (others return 400);
	// 960 is the smallest standard step that stays sharp on a high-DPI card.
	photoWidth = 960
)

var wikiClient = &http.Client{Timeout: 10 * time.Second}

// errRateLimited means Wikipedia asked us to slow down; the rest of the batch
// waits for the next poll instead of continuing to hit the API.
var errRateLimited = errors.New("wikipedia rate limit")

// thumbSize matches the "/320px-" width segment in a Wikimedia thumbnail URL.
var thumbSize = regexp.MustCompile(`/\d+px-`)

type wikiSummary struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Thumbnail   *struct {
		Source string `json:"source"`
	} `json:"thumbnail"`
	OriginalImage *struct {
		Source string `json:"source"`
		Width  int    `json:"width"`
	} `json:"originalimage"`
	ContentURLs struct {
		Desktop struct {
			Page string `json:"page"`
		} `json:"desktop"`
	} `json:"content_urls"`
}

// SyncPhotos looks up a photo for every driver that hasn't been checked yet.
func SyncPhotos(ctx context.Context, pool *pgxpool.Pool) error {
	drivers, err := db.DriversNeedingPhoto(ctx, pool)
	if err != nil {
		return err
	}
	if len(drivers) == 0 {
		return nil
	}

	checked, found := 0, 0
	for i, d := range drivers {
		// pause between every lookup, including failed ones, to stay well under Wikipedia's limits
		if i > 0 {
			time.Sleep(time.Second)
		}

		photo, page, err := findPhoto(ctx, d.FirstName+" "+d.LastName)
		if errors.Is(err, errRateLimited) {
			log.Printf("photos: rate limited, stopping; %d drivers left for next poll", len(drivers)-i)
			break
		}
		if err != nil {
			// leave unchecked so a temporary failure gets retried next poll
			log.Printf("photos: %s %s: %v", d.FirstName, d.LastName, err)
			continue
		}
		if err := db.SetDriverPhoto(ctx, pool, d.ID, photo, page); err != nil {
			return err
		}
		checked++
		if photo != nil {
			found++
		}
	}

	log.Printf("photos: checked %d drivers, found %d photos", checked, found)
	return nil
}

// findPhoto tries "<name> (racing driver)" first, because plain names like
// "George Russell" land on disambiguation pages or on someone else entirely.
// Returns nil, nil when no suitable page or image exists.
func findPhoto(ctx context.Context, name string) (photo, page *string, err error) {
	for _, title := range []string{name + " (racing driver)", name} {
		s, err := fetchSummary(ctx, title)
		if err != nil {
			return nil, nil, err
		}
		if s == nil || s.Type != "standard" {
			continue
		}
		// guard against a same-named person who isn't a driver
		if !strings.Contains(strings.ToLower(s.Description), "driver") {
			continue
		}
		img := imageURL(s)
		if img == "" {
			return nil, nil, nil
		}
		p := s.ContentURLs.Desktop.Page
		return &img, &p, nil
	}
	return nil, nil, nil
}

func fetchSummary(ctx context.Context, title string) (*wikiSummary, error) {
	u := wikiSummaryURL + url.PathEscape(strings.ReplaceAll(title, " ", "_"))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	// Wikimedia asks API clients to identify themselves
	req.Header.Set("User-Agent", "pitwall/0.1 (motorsport dashboard portfolio project)")

	resp, err := wikiClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errRateLimited
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia %s: status %d", title, resp.StatusCode)
	}

	var s wikiSummary
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

// imageURL picks a card-sized image: a resized thumbnail when the original is
// larger than photoWidth (originals can be several MB), otherwise the original
// itself, since Wikimedia won't upscale a thumbnail past the original's width.
func imageURL(s *wikiSummary) string {
	var u string
	switch {
	case s.OriginalImage != nil && s.OriginalImage.Width > 0 && s.OriginalImage.Width <= photoWidth:
		u = s.OriginalImage.Source
	case s.Thumbnail != nil:
		u = thumbSize.ReplaceAllString(s.Thumbnail.Source, fmt.Sprintf("/%dpx-", photoWidth))
	case s.OriginalImage != nil:
		u = s.OriginalImage.Source
	}
	// drop Wikipedia's tracking query string
	if i := strings.IndexByte(u, '?'); i >= 0 {
		u = u[:i]
	}
	return u
}
