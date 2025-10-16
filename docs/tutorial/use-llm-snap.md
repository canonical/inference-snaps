(use-llm-snap)=
# Get started with Qwen VL Inference Snap

This tutorial will guide you through the steps to install and use an Inference Snap.
We will specifically be using *Qwen 2.5 VL*, which is a multi-modal large language model.
This model is also known as a Vision Language Model (VLM), as it can ingest and interpret images, along with text prompt.

## Installing the snap

We install Qwen VL from the `2.5` track, which indicates the version of the model.
The snap is currently released with `beta` stability, indicated by the channel.

```{terminal}
:input: sudo snap install qwen-vl --channel 2.5/beta

qwen-vl (2.5/beta) 2.5 from Canonical IoT Labs✓ installed
```

During installation, the snap detects information about your system like installed hardware and available memory.
Based on this, an appropriate {ref}`engine <engines>` is chosen.
An *engine* in the context of Inference Snaps is a combination of model weights and a runtime, 
along with some utilities to make it work, like a server exposing a standard API.
Hardware vendors publish optimizations for some of these parts.
These optimization are packaged in silicon-specific engines.


## Check the status

After installation of the snap, a server is started that can be used for inferencing with the model.
This server comes from the selected engine.
To see the selected engine and the status of this server, run the status command:

```{terminal}
:input: qwen-vl status
engine: intel-cpu
status: online
endpoints:
    openai: http://localhost:8326/v3
```

In this case we see the `intel-cpu` engine was selected automatically.
Your system might get a different engine.
If the server started successfully, the status command includes the URL where the model can be accessed.


## Chat with the LLM

The `qwen-vl` snap includes a simple chat client that can be used to chat with the LLM.
A chat session can be started with the `chat` command:

```{terminal}
:input: qwen-vl chat
Connected to http://localhost:8326/v3
Type your prompt, then ENTER to submit. CTRL-C to quit.
» Hi there
Hello! How can I assist you today?

»  
```

This chat client is an example that only supports text input.
For more advanced queries, like ones that include images, we need to use a more advanced client.

## Run and configure Open WebUI

Open WebUI is an advanced client interface for working with LLMs.
It can be run as a Docker container:

```shell
docker run --network=host --env PORT=9099 ghcr.io/open-webui/open-webui:0.6
```

This command has a few customizations:
* `--network=host` is required so that the Docker container can access the URL of the Inference Snap on the host machine
* `--env PORT=9099` is to change the port which Open WebUI listens on, preventing collisions with other services on your host machine using the default `8080` port

After running the Open WebUI Docker container with the command above, the web interface can be accessed at [http://localhost:9099](http://localhost:9099).

Use the URL you obtained from the `status` command to configure a new connection in the settings of Open WebUI.
This is described in {ref}`configure-openwebui`.

## Prompt Qwen VL with an image

On the Open WebUI interface, start a *New Chat*.
At the top there is a dropdown menu listing available models.
If your configuration in the previous step was done correctly, you should see Qwen 2.5 VL being listed.
Select it.

The chat screen has a text field to type your question.
Attachments like images can be added for additional context by clicking on the `+` button.

```{image} attach-image.png
:alt: Selecting an image attachment
:align: center
```

As an example we attach [this image](https://commons.wikimedia.org/wiki/File:Portugal,_Porto,_Dominic_Lu%C3%ADs_l_Bridge_(52594422458).jpg) of a bridge.
We ask the question "What is the significance of the triangles?", hoping that the model would explain the architectural reason structures are tessellated.

After submitting the image and text prompt, it takes a bit of time to process the query, and then starts printing a response.
In this example we can see that the model spotted something in the image that we did not necessarily anticipate.

```{image} response.png
:alt: Response of the model
:align: center
```

## Conclusion

In this tutorial you learned how to install an Inference Snap.
A suitable engine was automatically installed, and its status and functionality were verified.
A more advanced client was installed that allows you to perform advanced queries.
As an example an image was attached to the text prompt, allowing you to ask questions related to the image.
