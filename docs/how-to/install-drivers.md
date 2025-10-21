(install-drivers)=
# Install drivers
Inference snaps require certain drivers available on the host.
The drivers should be installed before installing inference snaps to allow hardware detection and deployment of a matching {ref}`engine <engines>`. 

Refer below to find the specific instruction.

In case drivers were installed after the inference snap, perhaps by mistake or after installing new hardware, run the following command to repeat the engine selection:
```shell
sudo <inference-snap> use-engine --auto
```

(install-intel-gpu-drivers)=
## Install Intel GPU drivers

The user space drivers for Intel GPUs (integrated and discrete) are included in the snaps.

Lunar Lake and Battlemage GPUs may require a hardware enablement (HWE) kernel.
Please refer [here](https://dgpu-docs.intel.com/driver/client/overview.html) for details.


(install-intel-npu-drivers)=
## Install Intel NPU drivers

If your system has an Intel NPU, and you want to make use of it, you need to install the user space driver for it.
It is available as a snap:

```
sudo snap install intel-npu-driver
```


(install-nvidia-gpu-drivers)=
## Install NVIDIA GPU drivers

You need a working CUDA setup on the host machine to use an NVIDIA GPU.

The version of the driver and utilities might be different depending on your setup.
On Ubuntu Server 24.04, the following commands install the appropriate packages:

```
sudo apt update
sudo apt install nvidia-driver-550 nvidia-utils-550 nvidia-cuda-toolkit
sudo reboot
```
