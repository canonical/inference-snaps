(install-deepseek-r1-snap)=
# Install Deepseek R1 snap

The Deepseek R1 snap installs a hardware-optimized runtime for inference with the Deepseek R1 LLM.
It supports a range of hardware, in many cases with the help of drivers installed on the host.

This snap works with amd64 and some arm64 CPU architectures and has been tested on Intel, AMD, and Ampere CPUs.
It is also compatible with Intel and Nvidia GPUs, as well as Intel NPUs.
Check if your hardware is supported in the the full list of supported platforms. 

Follow these guides to install {ref}`Intel<install-intel-gpu-drivers>` and {ref}`Nvidia<install-nvidia-gpu-drivers>` GPU drivers, or {ref}`Intel NPU drivers<install-intel-npu-drivers>`.

```{note}
This snap has mainly been tested on Ubuntu 24.04 and newer systems.
For the best experience, make sure the drivers are installed on the host before installing the snap.
```
## Install the snap

Open a terminal, then set the right channel and install the snap:

```bash
sudo snap install deepseek-r1 --channel=<channel>
```

## Run server

Each stack installed consists of an inference engine and a model.
The inference engine is a server application.
To conserve computing resources, it is only run on demand by the user.

To start the server:

```bash
sudo snap start deepseek-r1
```

When the server is running, it exposes an an [OpenAI compatible endpoint](https://github.com/openai/openai-openapi) served via HTTP. 
The HTTP server's bind host and port have the following default values:

```{terminal} 
:input: sudo snap get deepseek-r1 http

Key        Value
http.host  127.0.0.1
http.port  8080
```

To change, for example, the HTTP port to 8999:

```bash
sudo snap set deepseek-r1 http.port=8999
```

Once changed, restart the server:

```bash
sudo snap restart deepseek-r1
```
Check out the {ref}`configuration section <configure-the-snap>` for details on configuring the snap.

To debug possible issues, you can query the server logs:

```bash
sudo snap logs deepseek-r1
```
To query more lines and follow the logs, use `-n -100 -f`.

## Chat with Deepseek R1

The snap ships an application that allows basic prompting from the command line.
To run this built-in chat application:

```bash
deepseek-r1.chat
```
(configure-the-snap)=
## Configure the snap

The Deepseek R1 snap can be configured in a few ways. 
Configurations can be explored and changed using the `snap get` and `snap set` commands.

For example, to query the top level config:

```bash
sudo snap get deepseek-r1
```

To query the available stacks:

```bash
snap get deepseek-r1 stacks
```

To change the stack(only possible if the snap was installed from the store):

```bash
sudo snap set deepseek-r1 stack=<stack>
```
To load a limited number of layers onto the GPU (only possible for CUDA-based stacks) use the n-gpu-layers snap option.
The remaining layers will run on the CPU.
```bash
sudo snap set deepseek-r1 n-gpu-layers=20
```


To reset to the default option, which is to load the entire model onto the GPU, unset the value:

```bash
sudo snap unset deepseek-r1 n-gpu-layers
```

## Uninstall the snap

To remove the Deepseek R1 snap:

```bash
sudo snap remove deepseek-r1
```
