#!/bin/bash -eu

engine="$(modelctl show-engine --format=json | jq -r .name)"
modelctl run "$SNAP/engines/$engine/server" --wait-for-components
