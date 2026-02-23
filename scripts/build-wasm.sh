#!/usr/bin/env bash
# Build the Jumbleberry Fields Solver as WebAssembly for GitHub Pages.
# Outputs go to docs/ which is the GitHub Pages root.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
DOCS_DIR="$ROOT_DIR/docs"

echo "==> Building WASM binary..."
cd "$ROOT_DIR"
GOOS=js GOARCH=wasm go build -o "$DOCS_DIR/solver.wasm" ./cmd/wasm

echo "==> Copying wasm_exec.js from Go SDK..."
GOROOT="$(go env GOROOT)"
cp "$GOROOT/lib/wasm/wasm_exec.js" "$DOCS_DIR/wasm_exec.js"

echo "==> Copying EV table..."
if [ -f "$ROOT_DIR/ev_table.json" ]; then
  cp "$ROOT_DIR/ev_table.json" "$DOCS_DIR/ev_table.json"
else
  echo "    WARNING: ev_table.json not found. The page will compute it on first load."
fi

echo ""
echo "Done! Files in docs/:"
ls -lh "$DOCS_DIR"/*.{wasm,js,html,json} 2>/dev/null | awk '{print "  " $5 "\t" $9}'
echo ""
echo "To test locally: cd docs && python3 -m http.server 8080"
