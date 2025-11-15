#!/usr/bin/env bash

# Demo script untuk mencoba endpoint secara berurutan.
# Pastikan server berjalan: `go run ./cmd/server` (port 8080).

BASE="http://localhost:8080"

echo "== Health =="
curl -s "$BASE/health" | sed -e 's/^/  /'
echo

echo "== Create item (Apple) =="
CREATE_RESP=$(curl -s -H "Content-Type: application/json" -d '{"name":"Apple"}' "$BASE/items")
echo "$CREATE_RESP" | sed -e 's/^/  /'
ID=$(echo "$CREATE_RESP" | sed -n 's/.*"id":"\([^"]*\)".*/\1/p')
echo "  Parsed ID: $ID"
echo

echo "== List items =="
curl -s "$BASE/items" | sed -e 's/^/  /'
echo

echo "== Get by ID =="
curl -s "$BASE/items/$ID" | sed -e 's/^/  /'
echo

echo "== Update (name=Banana, done=true) =="
curl -s -X PUT -H "Content-Type: application/json" -d '{"name":"Banana","done":true}' "$BASE/items/$ID" | sed -e 's/^/  /'
echo

echo "== Delete =="
curl -i -s -X DELETE "$BASE/items/$ID" | sed -e 's/^/  /'
echo

echo "== Get after delete (should be 404) =="
curl -i -s "$BASE/items/$ID" | sed -e 's/^/  /'
echo

echo "Done."