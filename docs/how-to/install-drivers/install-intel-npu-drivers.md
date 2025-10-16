(install-intel-npu-drivers)=
# Install Intel NPU drivers

If your system has an Intel NPU, and you want to make use of it, you need to install the userspace driver for it.
It is available as a snap:

```
sudo snap install intel-npu-driver
```

After installing the Inference Snap, connect the NPU driver to it:

```
sudo snap connect <inference-snap>:npu-libs intel-npu-driver
``` 
