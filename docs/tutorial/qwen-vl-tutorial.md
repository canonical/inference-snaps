(qwen-vl-tutorial)=
# Get started with Inference Snaps using Qwen VL

This tutorial walks you through the process of installing and using an inference snap.
Specifically, you will install and run *Qwen 2.5 VL*, a [multi-modal large language model (LLM)](https://www.nvidia.com/en-us/glossary/multimodal-large-language-models/).
This model is also referred to as a Vision Language Model (VLM) because it can ingest and interpret images, along with text prompts.

## Set up your computer

To complete this tutorial you will need:
- A computer running an operating system that supports snaps, e.g., Ubuntu 24.04. Older Ubuntu versions aren't supported.
- Hardware drivers. These may be pre-installed in your OS, but some devices require {ref}`additional drivers <install-drivers>`.
- A Docker installation

## Install the snap

Install the Qwen VL snap from the `2.5` track (which refers to the version of the model).
It is currently released with `beta` stability as indicated by the channel:

```{terminal}
:input: sudo snap install qwen-vl --channel 2.5/beta

qwen-vl (2.5/beta) 2.5 from Canonical IoT Labs✓ installed
```

During installation, the snap detects information about your system such as installed hardware and available memory.
It uses this information to select an appropriate engine for your setup. 

In this context, {ref}`engine <engines>` refers to a combination of model weights and a runtime, along with some utilities to make it work, like a server exposing a standard API.

## Check the status

Once the installation is complete, the server used to inference with the model starts automatically.
This server comes from the engine.

Run the status command to see the selected engine and the status of the server.
The output should include the server's status and the URL where the API can be accessed:

```{terminal}
:input: qwen-vl status
engine: intel-cpu
status: online
endpoints:
    openai: http://localhost:8326/v3
```

In this case, the `intel-cpu` engine was automatically selected during installation.
You may end up with a different engine.

## Chat with the LLM

The `qwen-vl` snap includes a basic chat client that can be used to chat with the LLM.
Use the `chat` command to start a chat session:

```{terminal}
:input: qwen-vl chat
Connected to http://localhost:8326/v3
Type your prompt, then ENTER to submit. CTRL-C to quit.
» Hi there
Hello! How can I assist you today?

»  
```

This chat client only supports text input.
For queries that include images, you'll need to set up a more advanced client.

## Run and configure Open WebUI

Open WebUI is an advanced client interface for working with LLMs.
Run the Open WebUI as a Docker container:

```shell
docker run --network=host --env PORT=9099 ghcr.io/open-webui/open-webui:0.6
```

* `--network=host` allows the Docker container to access the URL of the inference snap's HTTP server on the host machine
* `--env PORT=9099` changes the port that Open WebUI listens on from the default `8080` to `9099` to prevent collisions with other services on your host machine.

Confirm that you can access the web interface at [http://localhost:9099](http://localhost:9099), and {ref}`configure a new connection in Open WebUI <configure-openwebui>`.
You will need the URL from the status command.

## Prompt Qwen VL with an image

Once the new connection is set, start a *New Chat* on the Open WebUI interface.
In the dropdown menu at the top listing available models, select Qwen 2.5 VL.

``` {note}
If the model is not listed, check if the configuration in the previous step was done incorrectly.
```

The chat screen has a text field, but you can add images and other attachments for additional context by clicking on the `+` button.

```{image} attach-image.png
:alt: Selecting an image attachment
:align: center
```

Download and add this image of the [{spellexception}`Dom Luís I Bridge`](https://commons.wikimedia.org/wiki/File:Portugal,_Porto,_Dominic_Lu%C3%ADs_l_Bridge_(52594422458).jpg) as an attachment, and enter a suitable question in the chat screen.
The model will take some time to process the query before it starts printing a response.

In this example, the prompt used was "What is the significance of the triangles?". 
The expectation was that the LLM would explain the architectural reason for the triangular sections in the bridge.
However, the model spotted the triangular buoy and responded according to that instead of the bridge.

```{image} response.png
:alt: Response of the model
:align: center
```

## Next steps

To become more familiar with inference snaps:
- Learn how to {ref}`switch to a different engine <switch-between-engines>` after installation
- Configure an inference snap to {ref}`work with your IDE <use-in-ide>`
