#!/usr/bin/env bash

# Regenerates the diagrams in docs/generated from the code
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
out="$root/docs/generated"
mkdir -p "$out"

# Package import graph: which backend packages import which.
# Programs (package main) are blue, libraries are green.
cd "$root/backend"
go list -f '{{.Name}} {{.ImportPath}}{{range .Imports}} {{.}}{{end}}' ./... |
  awk '
    function short(p) { sub(/^pitwall\/backend\/?/, "", p); return p == "" ? "main" : p }
    BEGIN {
      print "digraph packages {"
      print "  rankdir=TB; bgcolor=\"transparent\"; nodesep=0.5; ranksep=0.7;"
      print "  node [shape=box, style=\"rounded,filled\", fontname=\"Helvetica\", fontsize=12, margin=\"0.25,0.12\", penwidth=0];"
      print "  edge [color=\"#8a8f98\", arrowsize=0.7];"
    }
    {
      pkg = short($2)
      color = ($1 == "main") ? "#3b82f6" : "#10b981"
      if ($1 == "main") progs = progs "\"" pkg "\"; "
      print "  \"" pkg "\" [fillcolor=\"" color "\", fontcolor=white];"
      for (i = 3; i <= NF; i++)
        if ($i ~ /^pitwall\/backend/) print "  \"" pkg "\" -> \"" short($i) "\";"
    }
    END { 
      print "  { rank=same; " progs "}"
      print "}" 
    }
  ' | dot -Tsvg -o "$out/packages.svg"

echo "wrote $out/packages.svg"
