# Machine hardware info

Each subdirectory represents a single machine.
This directory contains files with raw hardware info data from the respective machine.
It optionally contains a `hardware-info.json` file which is the output of the `cmd/hardware-info` application.

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
LC_ALL=POSIX df --portability --block-size=1 / /var/lib/snapd/snaps
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

## hardware-info.json (optional)

If we have access to the machine, and we can run the hardware-info application on it, add the output to this file.
It is used for cross-validation of the output between the actual hardware-info vs the hardware info from test data.

```
stack-utils.hardware-info --pretty --friendly
```
