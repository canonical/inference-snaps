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
Open the Open WebUI interface on your browser.
The default address is [http://localhost:8081](http://localhost:8081).
After logging in:

1. Click on your user account icon and select **Settings** from the dropdown menu.
2. Under **Settings**, select **Admin Settings** on the bottom left.
3. Select **Connections** from the left sidebar.
4. Under OpenAI API section, click on the **+** button to add a new connection.
5. In the **URL** field, enter the URL of the AI model snap you want to use.
  To identify the URL of the AI model snap, checkout this guide on [using an AI model snap via its OpenAI API](use-network-api.md).
6. Save the connection, and click **Save** again to save the settings.

```{tip}
By default, Open WebUI will first try to use a locally installed Ollama instance before it sends the prompt to the selected model.
This can cause an unnecessary delay.
We recommend disabling the Ollama and OpenAI connections in the Open WebUI admin settings if you are not using them:
1. Open settings in Open WebUI and choose **Admin Settings**.
2. In the **Admin Panel**, go to **Settings** > **Connections**.
3. Disable **OpenAI API** and **Ollama API** and save the changes.
```


On the left top corner, click on the **Select a model**.
Select the model name you configured under connections.

## Use the model

The model is now ready to be used.
Refer to the [Open WebUI documentation][open-webui-docs] to learn more about how to use it.

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
[open-webui-docs]: https://docs.openwebui.com/
