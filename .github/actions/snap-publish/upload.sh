#!/bin/bash -eu

function check_file() {
    local file=$1
    if [ ! -f "$file" ]; then
        echo "Error: file not found: $file"
        exit 1
    fi
}

channel=$1
arch=${2:-$(dpkg --print-architecture)} # if not set, take the current architecture

if [[ "$(yq --version)" != *v4* ]]; then
    echo "Please install yq v4."
    exit 1
fi

# validate channel
if [[ ! "$channel" =~ ^[a-z0-9-]+/[a-z0-9-]+(/[a-z0-9-]+)?$ ]]; then
    echo "Invalid Snap channel: $channel"
    exit 1
fi

# load snapcraft.yaml into variable, explode to evaluate aliases
snapcraft_yaml=$(yq '. | explode(.)' snap/snapcraft.yaml)

snap_name=$(echo "$snapcraft_yaml" | yq '.name')
snap_version=$(echo "$snapcraft_yaml" | yq '.version')
snap_file="${snap_name}_${snap_version}_${arch}.snap"
check_file "$snap_file"
snap_size=$(du -h "$snap_file" | cut -f1)

echo -e "Snap file:\n\t$snap_file $snap_size"

# Extract components from snapcraft.yaml
components=$(echo "$snapcraft_yaml" | yq '.components | to_entries | .[].key')

# Build components argument list
component_args=()
echo "Snap components:"
for comp_name in $components; do
    comp_ver=$(echo "$snapcraft_yaml" | yq ".components.$comp_name.version")
    comp_file="${snap_name}+${comp_name}_${comp_ver}.comp"
    check_file "$comp_file"
    comp_size=$(du -h "$comp_file" | cut -f1)
    echo -e "\t$comp_file $comp_size"

    component_args+=(--component "$comp_name=$comp_file")
done

echo -e "Channel:\n\t$channel"

snapcraft upload "$snap_file" "${component_args[@]}" --release="$channel"
