(configure-openwebui)=
# Configure Open WebUI for use with inference snaps

[Open WebUI][open-webui-docs] provides a convenient web interface for interacting with inference snaps that have an OpenAI endpoint.
This guide will help you set up and use an inference snap with Open WebUI, hosted on the same machine.

You must have Open WebUI and an inference snap installed to continue.

````{tip}

Open WebUI is also available as a snap.
To install [the snap](https://snapcraft.io/open-webui), run:

```shell
sudo snap install open-webui
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
