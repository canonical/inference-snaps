(install-drivers)=
# Install drivers
Inference snaps require certain drivers available on the host in order to function correctly.
You must install these drivers before installing inference snaps to enable hardware detection and deployment of a matching {ref}`engine <engines>`. 


In case you installed a driver after installing an inference snap, repeat the engine selection:
```shell
sudo <inference-snap> use-engine --auto
```

## Intel GPU

The user space drivers for Intel GPUs (integrated and discrete) are included in the snaps.

Lunar Lake and Battlemage GPUs may require a hardware enablement (HWE) kernel.
Please refer [here](https://dgpu-docs.intel.com/driver/client/overview.html) for details.


## Intel NPU

If your system has an Intel NPU and you want to make use of it, you need to install the user space driver for it.
It is available as a snap:

```
sudo snap install intel-npu-driver
```

## NVIDIA GPU

You need to install the NVIDIA driver on the host machine to use an NVIDIA GPU with inference snaps.

Inference snaps use CUDA 12 and depend on NVIDIA driver version **525 or higher**.
Refer to [CUDA Minor Version Compatibility](https://docs.nvidia.com/deploy/cuda-compatibility/minor-version-compatibility.html) for more details.

The driver may be pre-installed on your system. You can check the installed version with:
```shell
$ apt list --installed nvidia-driver* 
Listing... Done
nvidia-driver-580/noble-updates,noble-security,now 580.95.05-0ubuntu0.24.04.2 amd64 [installed]
```

If not installed, or if the installed version is lower than 525, you can install or upgrade the driver using the following commands:
```shell
sudo apt update
sudo apt install nvidia-driver-xxx
```
Replace `xxx` with the driver version (e.g., `525`, etc.). 
The latest version is usually recommended and found by running: `apt list --installed nvidia-driver*`

```{tip}
On Ubuntu Desktop, you can also install the driver using the "Additional Drivers" tool.
```

Reboot after installing or upgrading the NVIDIA driver:
```shell
sudo reboot
```
