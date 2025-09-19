(use-ai-model-snap-with-open-webui)=
# Use an AI model snap with Open WebUI

Open WebUI provides a convenient web interface for interacting with AI models packaged as snaps.
This guide will help you set up and use an AI model snap with Open WebUI.
You must have Open WebUI and an AI model snap installed to continue.
If you don't, check out the [Open WebUI installation guide][open-webui-install] and the [AI model snap installation guide](install-deepseek-r1-snap.md).

```{warning}
The default http port for Open WebUI is `8080`, which is also used by some AI model snaps.
To avoid conflicts, we recommend setting a different port for Open WebUI during installation.
```

## Configure Open WebUI

````{tip}
In case of using docker to run Open WebUI, you can set a different port by using the `PORT` environment variable.
Also it's necessary to use `--network=host` to allow Open WebUI to access the AI model snap's API.

For example:
```shell
docker run --network=host --env PORT=9099 ghcr.io/open-webui/open-webui:0.6
```
````
Open the Open WebUI interface on your browser: `http://localhost:<port>`
Register and logging in, then:

1. Click your account icon, select **Settings**, and then choose **Admin Settings** at the bottom left.
2. In the left sidebar, select **Connections**, and disable **OpenAI API** and **Ollama API** if they are not in use.
3. Enable **Direct Connections** and save your changes.

Next, create a new connection to the AI model snap:

1. Open the settings again and select **Connections**.
2. Click **+** to add a new connection.
3. In the **URL** field, enter the URL of the AI model snap you want to use, then save it.
  To find the correct URL, see the guide on [using an AI model snap via its OpenAI API](use-network-api.md).

On the left top corner, click on the **Select a model**.
Select the model name you configured under connections.

## Use the model

The model is now ready to be used.
Refer to the [Open WebUI documentation][open-webui-docs] to learn more about how to use it.

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
[open-webui-docs]: https://docs.openwebui.com/
