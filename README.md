# ML Snap Utils

This repo contains utilities used in snapping machine learning (AI) workloads.

## Build

The CLIs included in this repo can be built using the following commands.

Hardware Info:

```bash
go build github.com/canonical/ml-snap-utils/cmd/hardware-info
```

Select Stack:

```bash
go build github.com/canonical/ml-snap-utils/cmd/select-stack
```

To build a snap for these applications, run:

```bash
snapcraft -v
```

Then install the snap and connect the required interfaces:

```bash
sudo snap install --dangerous ./ml-snap-utils_*.snap
sudo snap connect ml-snap-utils:hardware-observe 
```

## Usage

### Hardware Info

A help message is printed out when providing the `-h` or `--help` flags.

```bash
$ ml-snap-utils.hardware-info -h
Usage of hardware-info:
  -file string
        Output json to this file. Default output is to stdout.
  -pretty
        Output pretty json. Default is compact json.
```

By default, the `hardware-info` application will print out a summary of the host system to `STDOUT` in compact JSON
format.
By specifying the `--pretty` flag, the JSON will be formatted for easier readability.
The `--file` argument allows writing the JSON data to a file, rather than to `STDOUT`.

Errors and warnings are printed to STDERR.

### Select Stack

The output from `hardware-info` can be piped into `select-stack`.
You need to provide the location of the stack definitions from which the selection should be made.

The result is written as json to STDOUT, while any other log messages are available on STDERR.

Example:

```bash
$ ml-snap-utils.hardware-info | ml-snap-utils.select-stack --stacks=test_data/stacks/
2024/12/10 11:28:03 Vendor specific info for Intel GPU not implemented
2024/12/10 11:28:03 Stack cpu-f32 not selected: not enough memory
2024/12/10 11:28:03 Stack fallback-cpu matches. Score = 4.000000
2024/12/10 11:28:03 Stack fallback-gpu not selected: any: could not find a required device
2024/12/10 11:28:03 Stack llamacpp-avx2 matches. Score = 3.200000
2024/12/10 11:28:03 Stack llamacpp-avx512 not selected: any: could not find a required device
{"name":"fallback-cpu","components":["llamacpp","model-q4-k-m-gguf"],"score":4}
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
