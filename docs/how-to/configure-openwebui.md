(configure-openwebui)=
# Configure Open WebUI for use with inference snaps

[Open WebUI][open-webui-docs] provides a convenient web interface for interacting with inference snaps that have an OpenAI endpoint.
This guide will help you set up and use an inference snap with Open WebUI.

You must have Open WebUI and an inference snap installed to continue.

````{note}
Open WebUI uses `8080` as its default HTTP port. 
This port is commonly used by other services like web servers and inference engines.
To avoid conflicts, use a different port for Open WebUI during installation.

In case you are using docker to run Open WebUI, you can set a different port using the `PORT` environment variable.
You must use `--network=host` option to allow Open WebUI to access the inference snap's API.

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

## Create a connection to the inference snap

Open settings in Open WebUI, select `Connections`, and click on the **+** icon to add a new connection.

In the `URL` field, enter the URL of the endpoint you want to use, then save it.
You can find the URL using the inference snap's `status` command.
Refer to {ref}`using OpenAI API <use-openai-api>` for details.

## Use the model

To learn how to use the model using Open WebUI, refer to the [Open WebUI documentation][open-webui-docs].

<!-- links -->
[open-webui-docs]: https://docs.openwebui.com/
