(configure-ai-model-snap-for-open-webui)=
# Configure an AI model snap for Open WebUI

Open WebUI provides a convenient web interface for interacting with AI models packaged as snaps.
This guide will help you set up and use an AI model snap with Open WebUI.

You must have Open WebUI and an AI model snap installed to continue.
For installation instructions, check out the [Open WebUI][open-webui-install] and the [AI model snap](install-deepseek-r1-snap.md) installation guides.

````{note}
Open WebUI uses `8080` as its default HTTP port. 
This port is also used by some AI model snaps.
To avoid conflicts, use a different port for Open WebUI during installation.

In case you are using docker to run Open WebUI, you can set a different port using the `PORT` environment variable.
You must use `--network=host` option to allow Open WebUI to access the AI model snap's API.

For example:
```shell
docker run --network=host --env PORT=9099 ghcr.io/open-webui/open-webui:0.6
```
````

## Enable direct connections in Open WebUI

Open the Open WebUI interface on your browser: `http://localhost:<port>`.
Register and log in, then click on your account icon and open `Settings`. 

Go to `Admin Settings` and select `Connections`.
If they are not in use, disable the OpenAI and Ollama APIs.
Enable `Direct Connections` and save your changes.

## Create a connection to the AI model snap

Open settings in Open WebUI, select `Connections`, and click on the **+** icon to add a new connection.

In the `URL` field, enter the URL of the AI model snap you want to use, then save it.
To find the correct URL, check out this guide on [using an AI model snap via its OpenAI API](use-network-api.md).

Click `Select a model` and select the model name you configured under connections.

## Use the model

To learn how to use the model using Open WebUI, refer to the [Open WebUI documentation][open-webui-docs].

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
[open-webui-docs]: https://docs.openwebui.com/
