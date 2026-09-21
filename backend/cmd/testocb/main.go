package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"pitwall/backend/ocb"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	drivers, err := ocb.FetchDrivers("formula1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got %d F1 drivers\n", len(drivers))

	for _, d := range drivers {
		fmt.Printf("%s %s\n", d.FirstName, d.LastName)
	}
}
