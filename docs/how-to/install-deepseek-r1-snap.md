(install-deepseek-r1-snap)=
# Install Deepseek R1 snap

The Deepseek R1 snap installs a hardware-optimized runtime for inference with the Deepseek R1 LLM.
It supports a range of hardware, in many cases with the help of drivers installed on the host.

This snap works with amd64 and some arm64 CPU architectures and has been tested on Intel, AMD, and Ampere CPUs.
It is also compatible with Intel and Nvidia GPUs, as well as Intel NPUs.

```{note}
    This snap has mainly been tested on Ubuntu 24.04 and newer systems.
    For the best experience, make sure the drivers are installed on the host 
    before installing the snap.
```
## Installation

Open a terminal then set the right channel and install the snap by executing:

```bash
    sudo snap install deepseek-r1 --channel=<channel>
```

During installation, the snap detects the hardware and picks a suitable stack.

## Using the snap
Two use cases have currently been tested for the Deepseek R1 snap.
You can run a server on the snap or use it with a chat client to prompt the LLM.

### Run server

Each stack that is installed consists of an inference engine and a model.
The inference engine is a server application.

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

### Chat with Deepseek R1

You can use a range of OpenAI-compatible chat clients to interact with the server.
The snap also ships an application that allows basic prompting from the command line.

To run the built-in chat application:

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

For CUDA-based stacks, the number of layers that are loaded on to the GPU can be configured as well.
By default, all layers are loaded into VRAM. 
This requires enough VRAM to fit the entire model. 
To only load a limited number of layers onto the GPU use the n-gpu-layers snap option:

```bash
    sudo snap set deepseek-r1 n-gpu-layers=20
```
This is useful if your GPU does not have enough VRAM to fit the entire model.
The remaining layers will run on the CPU.

To reset to the default option, which is to load the entire model onto the GPU, unset the value:

```bash
    sudo snap unset deepseek-r1 n-gpu-layers
```

## Manage components

This snap uses snap components to deploy optional artifacts.
You can check the installed components with:

```bash
    sudo snap components deepseek-r1
```

To remove a component, use:

```bash
    sudo snap remove <snap+component>
```