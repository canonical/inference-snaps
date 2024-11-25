# Hardware Detect

This program detects system hardware and provides a summary in JSON format.

## Build

To build the CLI for hardware-info, run the following command in the root of this repository:

```bash
go build github.com/canonical/hardware-info/cmd/hardware-info
```

To build a snap for this application, run:

```bash
snapcraft -v
```

Then install the snap and connect the required interfaces:

```bash
sudo snap install --dangerous ./hardware-info_*.snap
sudo snap connect hardware-info:hardware-observe
```

## Usage

A help message is printed out when providing the `-h` or `--help` flags.

```bash
$ hardware-info -h
Usage of hardware-info:
  -file string
        Output json to this file. Default output is to stdout.
  -pretty
        Output pretty json. Default is compact json.
```

By default, the `hardware-info` application will print out a summary of the host system to `STDOUT` in compact JSON format.
By specifying the `--pretty` flag, the JSON will be formatted for easier readability.
The `--file` argument allows writing the JSON data to a file, rather than to `STDOUT`.

Errors and warnings are printed to STDERR.

## Detecting NVIDIA GPU

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