(run-in-wsl)=
# Run an inference snap in WSL

Inference snaps can be installed in WSL just like on a native Ubuntu installation.
Because of the way GPUs are exposed on WSL, the hardware detection in inference snaps is not able detect the GPUs.
A CPU engine will therefore be automatically selected.

If your Windows machine has a GPU that is supported by {spellexception}`WSL's` GPU passthrough mechanism, some inference snaps can also use it.
Refer to {ref}`available snaps <available-snaps>` to see which inference snaps are supported on WSL.

## Prerequisites

The GPU drivers in Windows need to be updated to a recent version that supports exposing the GPU to WSL.
Refer to your GPU vendor's documentation for instructions on how to do this.

The rest of this how-to assumes you have followed Microsoft's guide to [install WSL](https://learn.microsoft.com/en-us/windows/wsl/install) and create the Ubuntu Linux distribution.

Fully updated and perform a restart of the Ubuntu distribution before continuing:
```shell
sudo apt update
sudo apt upgrade
sudo reboot
```

## Verify GPU availability in WSL

For NVIDIA GPUs, you can verify access using the `nvidia-smi` tool in WSL.

For Intel GPUs you can use the `clinfo` tool to check if the GPU is detected.
This tool is not available by default in WSL, but you can install it along with the Intel {spellexception}`OpenCL` drivers:
```shell
sudo apt update
sudo apt install clinfo intel-opencl-icd
clinfo
```

## Run inference snap

Install a supported inference snap.

From the list of engines, find the engine with a description that matches your GPU:
```shell
<inference-snap> list-engines
```

Switch to that engine and restart the snap:
```shell
sudo <inference-snap> use-engine <engine-name>
sudo snap restart <inference-snap>
```

You can verify that the engine is running by looking at the snap logs:
```shell
sudo snap logs <inference-snap> -f
```

If the snap starts successfully, the GPU is detected and the inference snap is ready to use.
