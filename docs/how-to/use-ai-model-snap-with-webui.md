(use-ai-model-snap-with-open-webui)=
# Use an AI model snap with Open WebUI

Open WebUI provides a convenient web interface for interacting with AI models packaged as snaps.
This guide will help you set up and use an AI model snap with Open WebUI.
You must have Open WebUI and an AI model snap installed to continue.
If you don't, check out the [Open WebUI installation guide][open-webui-install] and the [AI model snap installation guide](install-deepseek-r1-snap.md).

## Configure Open WebUI

Open the Open WebUI interface on your browser.
The default address is `http://localhost:8081`.
After logging in:

1. Click on your user account icon and select **Settings** from the dropdown menu.
2. Select **Connections** from the left sidebar.
3. Click on the **+** button to add a new connection.
4. In the **URL** field, enter the URL of the AI model snap you want to use.
  To identify the URL of the AI model snap, checkout this guide on [using an AI model snap via its OpenAI API](use-network-api.md).
5. The API key field is required but it's not used by the AI model snaps.
  You can enter any value here, such as, `dummy-key`.
6. Add the model name in the **Model IDs** field.
  Refer to the guide mentioned in step 4 to identify the correct model name.
7. Save the connection, and click **Save** again to save the settings.

```{tip}
By default, Open WebUI will first try to use a locally installed Ollama instance before it sends the prompt to the selected model.
This can cause an unnecessary delay.
We recommend disabling the Ollama and OpenAI connections in the Open WebUI admin settings if you are not using them:
1. Open settings in Open WebUI and choose **Admin Settings**.
2. The **Admin Panel** will open. Go to **Settings** > **Connections**.
3. Then disable **OpenAI API** and **Ollama API**. Click **Save**.
```


On the left top corner, click on the **Select a model**.
Select the model name you configured under connections.

The model is now ready to be used.

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
