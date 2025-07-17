(install-intel-npu-drivers)=
# Install Intel NPU drivers

If your system has an Intel NPU, and you want to make use of it, you need to install the userspace driver for it.
It is available as a snap:

```
sudo snap install intel-npu-driver
```

After installing the AI model snap, connect the NPU driver to it:

```
sudo snap connect <ai-model>:npu-libs intel-npu-driver
``` 
