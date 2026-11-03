#!/bin/bash -eux

# This script can be used to locally run the engine-tps test on a single tesflinger machine.
# Export the environment variables listed in the README, and then run this script.

envsubst < testflinger.yaml > testflinger.temp.yaml

testflinger submit --poll testflinger.temp.yaml
