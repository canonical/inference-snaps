(use-ai-model-snap-with-open-webui)=
# Use an AI model snap with Open WebUI

Open WebUI provides a convenient web interface for interacting with AI models packaged as snaps.
This guide will help you set up and use an AI model snap with Open WebUI.
You must have Open WebUI and an AI model snap installed to continue.
If you don't, check out the [Open WebUI installation guide][open-webui-install] and the [AI model snap installation guide](install-deepseek-r1-snap.md).

## Configure Open WebUI

First open the Open WebUI interface on your browser.
The default address is `http://localhost:8081`.
After logging in:

1. Click on your user account icon.
2. Select **Settings** from the dropdown menu.
3. Select **Connections** from the left sidebar.
4. Click on the **+** button to add a new connection.
5. In the **URL** field, enter the URL of the AI model snap you want to use.
  To identify the URL of the AI model snap, see the [Use an AI model snap via its OpenAI API](use-network-api.md) guide.
6. The API key field is required, but it is not used by the AI model snaps.
  You can enter any value here, such as `dummy-key`.
7. On the **Model IDs** field add the model name.
  Refer to the same guide mentioned in step 5 to identify the correct model name.
8. Click on **Save** to save the connection, and **Save** again to save the settings.

```{note}
By default, Open WebUI will first try to use a locally installed Ollama instance, before it sends the prompt to the selected model.
This can cause an unnecessary delay.
We recommend that you disable the Ollama and OpenAI connections in the Open WebUI admin settings if you are not using them:
1. Click on your user account icon, then select **Settings** and choose **Admin Settings**.
2. The **Admin Panel** will open. Go to **Settings** > **Connections**.
3. Then disable **OpenAI API** and **Ollama API**. Click **Save**.
```


On the left top corner, click on the **Select a model**.
Select the model name you configured under connections.

The model is now ready to be used.

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
