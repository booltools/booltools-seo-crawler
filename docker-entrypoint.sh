#!/bin/sh
set -e

cd /app

echo "Starting Go backend on port ${PORT:-8080}..."
./server &
GO_PID=$!

echo "Starting Astro frontend on port 4321..."
cd /app/web/dist
node server/entry.mjs &
WEB_PID=$!

trap "kill $GO_PID $WEB_PID 2>/dev/null; exit 0" SIGTERM SIGINT

wait $GO_PID $WEB_PID
