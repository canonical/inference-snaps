(install-deepseek-r1-snap)=
# Deepseek-R1 snap

The Deepseek-R1 snap installs a hardware-optimized runtime for inference with the [Deepseek-R1](https://github.com/deepseek-ai/DeepSeek-R1) Large Language Model.
It supports a range of hardware, in many cases with the help of drivers installed on the host.

This snap works with amd64 and some arm64 CPU architectures and has been tested on Intel, AMD, and Ampere CPUs.
It is also compatible with Intel and Nvidia GPUs, as well as Intel NPUs.

```{note}
This snap has mainly been tested on Ubuntu 24.04 and newer systems.
For the best experience, make sure the drivers are installed on the host before installing the snap.
```

## Install the snap

If you have a GPU or NPU, follow the appropriate guide to install its driver:
* {ref}`Intel GPU<install-intel-gpu-drivers>`
* {ref}`Nvidia GPU<install-nvidia-gpu-drivers>`
* {ref}`Intel NPU<install-intel-npu-drivers>`

After installing the drivers, install the snap from the specified channel:

```bash
sudo snap install deepseek-r1 --channel=<channel>
```

During installation, the most appropriate combination of inference engine and model is selected for your system.
The combination of engine and model, along with any required utilities or configs is called a *stack*.

## Run the server

All stacks expose a server application.
To conserve computing resources, it is only run on demand by the user.

To start the server:

```bash
sudo snap start deepseek-r1
```

When the server is running, it exposes an [OpenAI compatible endpoint](https://github.com/openai/openai-openapi) served via HTTP. 


## Chat with Deepseek-R1

The snap ships an example application that allows basic prompting from the command line.
To run this built-in chat application:

```bash
deepseek-r1.chat
```

You can also use any other OpenAI API compatible client to interact with this model via the network API. <!-- todo add link to network doc -->


## Checking the logs

To debug possible issues, you can query the server logs:

```bash
sudo snap logs deepseek-r1
```
To query more lines and follow the logs, use `-n -100 -f`.


## Change the stack

Query the available stacks:

```bash
sudo snap get deepseek-r1 stacks
```

Query details about a specific stack:

```{terminal}
:input: sudo snap get deepseek-r1 stacks.<name>

Key                           Value
stacks.<name>.compatible      true
...
```

This will print many fields describing what this stack is consisted of.
It is important to look at the `compatible` field to know if this stack will work on your hardware.

To change the stack that is used, set the stack option to the name of the stack you want.
Use this command **without** the `stack.` prefix:

```bash
sudo snap set deepseek-r1 stack=<name>
```

## Configure the snap

Configurations can be explored and changed using the `snap get` and `snap set` commands.

For example, to query the top level config:

```bash
sudo snap get deepseek-r1
```

Not all these configs can be changed, as they are defined by the selected stack.

### HTTP server configurations

The HTTP server's bind host and port are user changeable and can be seen under the `http` key:

```{terminal} 
:input: sudo snap get deepseek-r1 http

Key        Value
http.host  127.0.0.1
http.port  8080
```

To change the HTTP port to, for example, 8999:

```bash
sudo snap set deepseek-r1 http.port=8999
```

If you need other machines on the network to use this server, you need to change the bind host to 0.0.0.0:

```bash
sudo snap set deepseek-r1 http.host=0.0.0.0
```

Once changed, restart the server:

```bash
sudo snap restart deepseek-r1
```

### Model layers on GPU

CUDA-based stacks allow limiting the number of model layers that are loaded onto the GPU, by using the `n-gpu-layers` snap option.
This is useful if the GPU does not have enough vRAM.
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
