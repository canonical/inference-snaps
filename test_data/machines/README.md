# Machine hardware info

Each subdirectory represents a single machine.
This directory contains files with raw hardware info data from the respective machine.

## cpuinfo.txt

```
cat /proc/cpuinfo
```

## lspci.txt

```
lspci -vmmnD
```

## uname-m.txt

```
uname -m
```

## disk.txt

```
LC_ALL=POSIX df --portability --block-size=1 /var/lib/snapd/snaps
```

## meminfo.txt

```
cat /proc/meminfo
```

## additional-properties.json (optional)

Normally the additional properties are looked up on the host, using vendor specific tools.
This is not possible during testing when we do not have access to the host.
If any additional properties need to be added to PCI devices, add it to this file, based on the device slot address.

```
{
  "0000:00:02.0": {
    "vram": "14482374656",
    "compute_capability": "12.4"
  }
}
```

## machine-root

### AMD GPU testing
Navigate to the desired machine directory and create a `machine-root` directory.
```
mkdir -p machine-root/sys/class/kfd/kfd/topology/nodes
cp -ra /sys/class/kfd/kfd/topology/nodes/* machine-root/sys/class/kfd/kfd/topology/nodes/
mkdir -p machine-root/sys/class/drm
cp -a --parent /sys/class/drm/renderD* machine-root
```
Now find the file pointed by the symlink(s) `machine-root/sys/class/drm/renderD*`:
`ls -lah machine-root/sys/class/drm/renderD*`
And create the same symlink(s) in the test environment:
`mkdir -p machine-root/sys/class/drm/{PATH_RETRIEVED_FROM_ABOVE_COMMAND}`
e.g. `mkdir -p machine-root/sys/class/drm/../../devices/pci0000:00/0000:00:08.1/0000:c4:00.0/drm/renderD128`
To complete the setup, copy the file `mem_info_vram_total`:
`cp -a --parent /sys/bus/pci/devices/0000\:c4\:00.0/mem_info_vram_total machine-root/`
Substitute the pciSlot in the above command with the one corresponding to the GPU being tested, which can be found in the `lspci.txt` file.
In case the directory `machine-root/sys/class/drm/renderD128` is empty:
```
cd machine-root/sys/class/drm/renderD128
ln -s ../.. device
```
