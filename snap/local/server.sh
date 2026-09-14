#!/bin/bash -eu

modelctl export-status

port="$(modelctl get http.port)"
host="$(modelctl get http.host)"

echo "Starting mock OpenAI server on $host:$port"
export PYTHONUNBUFFERED=1
exec modelctl run -- "$SNAP"/bin/mock-openai-server/server.py --port "$port" --host "$host" --reasoning --delay 0.05
