#!/bin/bash -e

usage="Usage: ./dev/install.sh [--engine=<engine>] [--clean]"

engine=""
clean=false

for arg in "$@"; do
    case "$arg" in
        --engine=*)
            engine="${arg#*=}"
            ;;
        --clean)
            clean=true
            ;;
        --help)
            echo $usage
            exit 0
            ;;
        *)
            echo "Error: unknown argument '$arg'"
            echo $usage
            exit 1
            ;;
    esac
done

if [[ ! -f snap/snapcraft.yaml ]]; then
	echo "Error: snap/snapcraft.yaml not found. Run this script from the project's root." >&2
	exit 1
fi

# Make sure there is just one snap file
snap_file="$(ls *.snap)"
if [ "$(echo "$snap_file" | wc -l)" -ne 1 ]; then
    echo -e "Error: expected 1 snap file, found multiple: \n$snap_file"
    exit 1
fi

snap_name="$(echo "$snap_file" | cut -d'_' -f1)" # the part before the first underscore

# Make sure there is one .comp file per component
component_list=()
for comp_file in *.comp; do
    # Component file name patterns:
    # <snap_name>+<comp_name>_<comp_version>.comp
    # <snap_name>+<comp_name>.comp
    comp_name=$(echo "$comp_file" | 
        cut -d'+' -f2 | # drop snap name
        cut -d'.' -f1 | # drop file extension
        cut -d'_' -f1) # split by _, take 1st part
    
    # Check for duplicate components
    for existing_comp in "${component_list[@]}"; do
        if [[ "$existing_comp" == "$comp_name" ]]; then
            echo "Error: more than one component is named '$comp_name'"
            exit 1
        fi
    done
    component_list+=("$comp_name")
done

# Validate engine name only when provided
if [[ -n "$engine" ]]; then
    engine_file=./engines/$engine/engine.yaml
    if [[ ! -f "$engine_file" ]]; then
        echo "Unknown engine: $engine"
        exit 1
    fi
fi

if [[ "$clean" == true ]]; then
    echo "➤ Removing existing snap installation..."
    sudo snap remove "$snap_name" --purge
fi


echo "➤ Installing..."
sudo snap install --dangerous *.snap *.comp

echo "➤ Stop failing services..."
# On a fresh install, the snap is unable to auto-select an engine without hardware access.
# Stop since the service is going to fail without an engine
sudo snap stop "$snap_name"

echo "➤ Connecting interfaces [home, hardware-observe, process-control]..."
sudo snap connect $snap_name:home
sudo snap connect $snap_name:hardware-observe
sudo snap connect $snap_name:process-control || true # process-control is not available on all snaps

if [[ -n "$engine" ]]; then
    echo "➤ Setting engine to $engine..."
    sudo "$snap_name" use-engine "$engine" --assume-yes
fi

echo "➤ Starting services..."
sudo snap start "$snap_name"
