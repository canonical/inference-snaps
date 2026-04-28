#!/bin/bash
set -euo pipefail

port="$(modelctl get webui.http.port)"
host="$(modelctl get webui.http.host)"

# The capabilities depend on the model and engine size
capabilities="text, text:markdown, vision"

exec modelctl serve-webui "$SNAP/webui" --port "$port" --host "$host" --capabilities "$capabilities"
