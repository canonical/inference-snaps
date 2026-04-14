(configure-openwebui)=
# Configure Open WebUI for use with inference snaps

[Open WebUI][open-webui-docs] provides a convenient web interface for interacting with inference snaps that have an OpenAI endpoint.
This guide will help you set up and use an inference snap with Open WebUI, hosted on the same machine.

You must have Open WebUI and an inference snap installed to continue.

``````{note}
Open WebUI uses `8080` as its default HTTP port.
This port is commonly used by other services like web servers and inference engines.
To avoid conflicts, use a different port for Open WebUI.

`````{tabs}

````{group-tab} Snap

If you installed the Open WebUI snap, you can change its network `host` or `port`.
For example:

```shell
snap set open-webui port=9099
```

To view the current Open WebUI URL:

```shell
open-webui
```

````

````{group-tab} Docker

In case you are using docker to run Open WebUI, you can set a different port using the `PORT` environment variable.
You must use `--network=host` option to allow Open WebUI to access the inference snap's API.

For example:
```shell
docker run --network=host --env PORT=9099 ghcr.io/open-webui/open-webui:latest
```

````

`````
``````

## Create a connection to the inference snap

Open the Open WebUI interface on your browser: `http://localhost:<port>`.

Open **Settings** in Open WebUI and enter the **Admin Settings** page, then select the **Connections** tab. Ensure that **OpenAI API** is enabled, and click the **+** button next to _Manage OpenAI API Connections_ to add a new connection.

In the `URL` field, enter the URL of the endpoint you want to use, then save it.
You can find the URL using the inference snap's `status` command.
Refer to {ref}`using OpenAI API <use-openai-api>` for details.

## Use the model

To learn how to use the model using Open WebUI, refer to the [Open WebUI documentation][open-webui-docs].

<!-- links -->
[open-webui-docs]: https://docs.openwebui.com/
