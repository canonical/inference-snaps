#!/bin/bash -eu

git_branch=$(git rev-parse --abbrev-ref HEAD)
default_channel="latest/edge/$git_branch"

channel=${1-$default_channel}

snap_file="$(ls *.snap)"

if [ "$(echo "$snap_file" | wc -l)" -ne 1 ]; then
    echo -e "Error: expected 1 snap file, found multiple: \n$snap_file"
    exit 1
fi

snap_size=$(du --apparent-size -b "$snap_file" | cut -f1)
total_size=$snap_size

snap_size_human=$(numfmt --to=iec "$snap_size")
echo -e "Snap file:\n\t$snap_file $snap_size_human"

# Build components argument list
component_args=()
component_list=()
max_comp_size=$((5 * 1024 * 1024 * 1024)) # 5G
echo "Snap components:"
for comp_file in *.comp; do

    # Component file name patterns:
    # <snap_name>+<comp_name>_<comp_version>.comp
    # <snap_name>+<comp_name>.comp
    comp_name=$(echo "$comp_file" | 
        cut -d'+' -f2 | # drop snap name
        cut -d'.' -f1 | # drop file extension
        cut -d'_' -f1) # split by _, take 1st part

    comp_size=$(du --apparent-size -b "$comp_file" | cut -f1)
    total_size=$((total_size + comp_size))

    comp_size_human=$(numfmt --to=iec "$comp_size")
    
    # Red warning for components smaller than 20KB
    empty_comp_warning=""
    if [[ $comp_size -lt 20000 ]]; then
        empty_comp_warning="\033[31m(empty?)\033[0m"
    fi

    echo -e "\t$comp_file $comp_size_human $empty_comp_warning"

    if [[ $comp_size -gt $max_comp_size ]]; then
        echo "Error: component '$comp_name' ($comp_file) is $comp_size_human, which exceeds the 5G limit. Aborting upload."
        exit 1
    fi

    # Check for duplicate components
    for existing_comp in "${component_list[@]}"; do
        if [[ "$existing_comp" == "$comp_name" ]]; then
            echo "Error: more than one component is named '$comp_name'"
            exit 1
        fi
    done
    component_list+=("$comp_name")

    component_args+=(--component "$comp_name=$comp_file")
done

echo -e "Channel:\n\t$channel"

total_size_human=$(numfmt --to=iec "$total_size")
echo -ne "\nUpload $total_size_human and release? [y/N] "
read confirmation
if [[ "$confirmation" != "y" && "$confirmation" != "Y" ]]; then
    exit 1
fi

set -x
snapcraft upload "$snap_file" "${component_args[@]}" --release="$channel"
