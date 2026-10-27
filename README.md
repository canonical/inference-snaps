# Inference Snaps CLI

Configuration and engine management for inference snaps.

## Resources

📚 **[Documentation](https://documentation.ubuntu.com/inference-snaps/)**, learn how to use inference snaps and find the [CLI reference](https://documentation.ubuntu.com/inference-snaps/reference/models-cli/)

💬 **[Discussions](https://github.com/canonical/inference-snaps/discussions)**, ask questions and share ideas

🐛 **[Issues](https://github.com/canonical/inference-snaps/issues)**, report bugs and request features

## Development

### Prerequisites

- [Go](https://go.dev/doc/install)
- [Snapcraft](https://ubuntu.com/docs/snapcraft/9/how-to/set-up-snapcraft/), including LXD for building the snap

### Run tests

```bash
make test
```

### Build binaries

The CLIs included in this repo can be built using the following command:

```bash
make build
```

### Test snap

The [modelctl-test](test_data/snap/) snap is used to develop and test the CLI in a snap environment.
It bundles the CLI with test engines, runtimes, models, and a mock OpenAI server.

The CLI is built on the host and staged into the snap project, together with the test data, before packing.

To build the snap and its components, run:

```bash
make snap
```

To install the latest build and connect the required interfaces, run:

```bash
make install
```

Note: this does not rebuild the snap. Run `make snap` first after making changes.

To remove build outputs and the Snapcraft build container, run `make clean`.

## Usage

The following assumes use of the modelctl-test snap to use the CLI.

### Machine Info

A summary of the current host machine can be obtained by running:

```
modelctl-test machine
```

This prints a machine-readable summary of the host system. 

Errors and warnings are printed as standard errors.
This allows piping the output to another application.

### Select Engine

This command can be used to perform engine selection using static data.
It is useful for testing purposes.

To use, pipe the machine info in **JSON format** into `select-engine`.
You also need to provide the location of the engine manifests from which the selection should be made.

The result is printed as **JSON** to the standard output by default (use `--format=yaml` for YAML output), while any other log messages are written as standard errors.

Example:

```bash
$ modelctl-test machine --format=json | modelctl-test debug select-engine --engines test_data/engines/
❌ ampere - not compatible: required cpu device not found
❌ ampere-altra - not compatible: required cpu device not found
❌ arm-neon - not compatible: required device not found
✅ cpu-avx1 - compatible, score = 14
✅ cpu-avx2 - compatible, score = 17
❌ cpu-avx512 - not compatible: required cpu device not found
🟠 cpu-exptl - score = 12, experimental
✅ cuda-generic - compatible, score = 107
✅ example-memory - compatible, score = 18
✅ intel-cpu - compatible, score = 18
✅ intel-gpu - compatible, score = 72
❌ intel-npu - not compatible: required device not found
Selected engine for your hardware configuration: cuda-generic

{
  "engines": [
    {
      "name": "ampere",
      "description": "Test ampere selection",
      ...
    }
  ],
  "top-engine": "cuda-generic"
}
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
