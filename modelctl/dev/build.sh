#!/bin/bash -eu

if [[ ! -f snap/snapcraft.yaml ]]; then
	echo "Error: snap/snapcraft.yaml not found. Run this script from the project's root." >&2
	exit 1
fi

echo "➤ Removing old snap and component files..."
rm -fv *.snap *.comp

echo "➤ Packing snap and components..."
snapcraft -v pack
