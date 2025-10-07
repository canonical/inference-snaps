# Famous Models CLI

This repo contains the CLI implementation for configuration and engine management.

## Development

### Run tests

```bash
go test -count 1 -failfast ./...
```

### Build binaries

The CLIs included in this repo can be built using the following commands:

```bash
go build ./cmd/cli
```

### Build snap

To build a snap for these applications, run:

```bash
snapcraft -v
```

Then install the snap and connect the required interfaces:

```bash
sudo snap install --dangerous ./stack-utils_*.snap
sudo snap connect stack-utils:hardware-observe 
```

> [!NOTE]
> This snap is named `stack-utils` for historical reasons. It will be refactored in the future.

## Installation

```bash
sudo snap install stack-utils --devmode
```

To build and install from source, refer to [here](#build-snap).

## Usage

The following assumes use of the stack-utils snap to use the CLI.

### Machine Info

A summary of the current host machine can be obtained by running:

```
stack-utils show-machine
```

This prints a machine-readable summary of the host system. 

Errors and warnings are printed as standard errors.
This allows piping the output to another application.

### Select Engine

This command can be used to perform engine selection using static data.
It is useful for testing purposes.

To use, pipe the machine info in JSON format into `select-engine`.
You also need to provide the location of the engine manifests from which the selection should be made.

The result is printed as JSON to the standard output, while any other log messages are written as standard errors.

Example:

```bash
$ stack-utils show-machine --format=json | stack-utils debug select-engine --engines test_data/engines/
❌ ampere - not compatible: devices all: required cpu device not found
❌ ampere-altra - not compatible: devices all: required cpu device not found
❌ arm-neon - not compatible: devices any: required device not found
✅ cpu-avx1 - compatible, score = 14
✅ cpu-avx2 - compatible, score = 17
❌ cpu-avx512 - not compatible: devices all: required cpu device not found
🟠 cpu-devel - score = 12, grade = devel
❌ cuda-generic - not compatible: devices all: required pci device not found
✅ example-memory - compatible, score = 18
✅ intel-cpu - compatible, score = 18
❌ intel-gpu - not compatible: devices any: required device not found
❌ intel-npu - not compatible: devices any: required device not found
Selected engine for your hardware configuration: example-memory

{"engines":[{"name":"ampere","description":"Test ampere selection" ...
```

## Notes

### Detecting NVIDIA GPU

On a clean 24.04 installation, you need to install the NVIDIA drivers and utils:

```
sudo apt install nvidia-driver-550-server nvidia-utils-550-server
sudo reboot
```

After a reboot run `nvidia-smi` to verify it is working:

```
$ nvidia-smi    
+-----------------------------------------------------------------------------------------+
| NVIDIA-SMI 550.127.05             Driver Version: 550.127.05     CUDA Version: 12.4     |
|-----------------------------------------+------------------------+----------------------+
| GPU  Name                 Persistence-M | Bus-Id          Disp.A | Volatile Uncorr. ECC |
| Fan  Temp   Perf          Pwr:Usage/Cap |           Memory-Usage | GPU-Util  Compute M. |
|                                         |                        |               MIG M. |
|=========================================+========================+======================|
|   0  Quadro T2000 with Max-Q ...    Off |   00000000:01:00.0 Off |                  N/A |
| N/A   49C    P0              8W /   35W |       1MiB /   4096MiB |      0%      Default |
|                                         |                        |                  N/A |
+-----------------------------------------+------------------------+----------------------+
...
```
