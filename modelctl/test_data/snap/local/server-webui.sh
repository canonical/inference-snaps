#!/bin/bash
set -euo pipefail

port="$(modelctl get webui.http.port)"
host="$(modelctl get webui.http.host)"

exec modelctl serve-webui "$SNAP/webui" --port "$port" --host "$host"
