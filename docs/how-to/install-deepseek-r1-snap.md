(install-deepseek-r1-snap)=
# Install DeepSeek-R1 snap

The DeepSeek-R1 snap installs a hardware-optimized engine for inference with the [DeepSeek-R1](https://github.com/deepseek-ai/DeepSeek-R1) Large Language Model.
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

After installing the drivers, install the snap:

```bash
sudo snap install deepseek-r1
```

During installation, the most appropriate combination of inference runtime and model weights is selected for your system.
The combination of runtime and model weights, along with any required utilities or configurations is called an *engine*.

## Checking the status

After installation, the selected engine is started.
You can see which engine is selected, and its current state by using the status command:

```
deepseek-r1 status
```

## Chat with DeepSeek-R1

The snap ships an example application that allows basic prompting from the command line.
To run this built-in chat application:

```bash
deepseek-r1 chat
```

You can also use any other OpenAI API compatible client to interact with this model via the network API.
To do this, configure your client with the URL shown by the status command.


## Change the engine

It is possible to manually change to another engine.
To do this, first check which engines are available and compatible with your machine:

```bash
deepseek-r1 list-engines
```

Check the `COMAPT` column to know if this engine will work on your hardware.

You can also get more details about a specific engine:

```bash
deepseek-r1 show-engine <name>
```

This will print many fields describing what this engine consists of.

The engine that is used can be changed:

```bash
sudo deepseek-r1 use-engine <name>
```


## Start and stop the server

All engines contain a server application.
It is started automatically after installation.
To conserve computing resources, it can be stopped and started on demand by the user.

To stop the server:

```bash
sudo snap stop deepseek-r1
```

To start the server:

```bash
sudo snap start deepseek-r1
```

When the server is running, it exposes an [OpenAI compatible endpoint](https://github.com/openai/openai-openapi) served via HTTP. 


## Checking the logs

To debug possible issues, you can query the server logs:

```bash
sudo snap logs deepseek-r1
```
To query more lines and follow the logs, use `-n -100 -f`.

## Configure the snap

Configurations can be explored and changed using the `deepseek-r1 get` and `deepseek-r1 set` commands.

All configurations can be seen by running:

```bash
deepseek-r1 get
```

Not all configurations can be changed, as they are defined by the selected engine.

### HTTP server configurations

The HTTP server's bind host and port are user changeable and can be seen under the `http` key:

```{terminal} 
:input: deepseek-r1 get http

{
	"base-path": "v3",
	"host": "127.0.0.1",
	"port": 8080
}
```

To change the HTTP port to, for example, 8999:

```bash
sudo deepseek-r1 set http.port=8999
```

If you need other machines on the network to use this server, you need to change the bind host to 0.0.0.0:

```bash
sudo deepseek-r1 set http.host=0.0.0.0
```

Once changed, restart the server:

```bash
sudo snap restart deepseek-r1
```

### Model layers on GPU

CUDA-based stacks allow you to limit the number of model layers that are loaded onto the GPU by using the `n-gpu-layers` snap option.
This is useful if the GPU does not have enough vRAM.
The remaining layers will run on the CPU.
```bash
sudo deepseek-r1 set n-gpu-layers=20
```

[//]: # (To reset to the default option, which is to load the entire model onto the GPU, unset the value)

## Uninstall the snap

To remove the DeepSeek R1 snap:

```bash
sudo snap remove deepseek-r1
```
