# Install NVIDIA GPU drivers

You need a working CUDA setup on the host machine to use these snaps with an NVIDIA GPU.

The version of the driver and utilities might be different depending on your setup.
On Ubuntu Server 24.04, the following commands install the appropriate packages:

```
sudo apt update
sudo apt install nvidia-driver-550 nvidia-utils-550 nvidia-cuda-toolkit
sudo reboot
```
