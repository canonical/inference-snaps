(use-ai-model-snap-with-open-webui)=
# Use an AI model snap with Open WebUI

Open WebUI provides a convenient web interface that can be used for interact with AI models packaged as snaps.
This guide will help you set up and use an AI model snap with Open WebUI.
This guide assumes Open WebUI and an AI model snap are installed.
If not, see the [Open WebUI installation guide][open-webui-install] and the [AI model snap installation guide](install-deepseek-r1-snap.md).

## Configure Open WebUI

First open the Open WebUI interface on your browser.
The default address is `http://localhost:8081`.
After logging in:

1. Click on your account user icon.
2. Select **Settings** from the dropdown menu.
3. Select **Connections** from the left sidebar.
4. Click on the **+** button to add a new connection.
5. On the **URL** field, enter the URL of the AI model snap you want to use.
  For example, if you installed the Deepseek R1 model snap, the URL would be `http://localhost:8080/v1`.
  To identify the URL of the AI model snap, see the [Use an AI model snap via its OpenAI API](use-network-api.md) guide.
6. The API key field is demanded by the tool, although it is not used by the AI model snaps.
  You can enter any value here, such as `dummy-key`.
7. Click on **Save** to save the connection, and **Save** again to save the settings.

```{note}
If you have Ollama installed, a few addtional steps are required after configuring the connection:

1. Go to the **Admin Panel** > **Settings** > **Connections**.
2. Then disable Ollama API.

In case you're using [OpenVINO™ Model Server][open-vino-model-server] also disable OpenAI API option.
Also, for OpenVINO™ Model Server, model ID needs to be set.
Check out [Use an AI model snap via its OpenAI API](use-network-api.md) to know how to get the model ID.
```


On the left top corner, click on the **Select a model** and select the one described as `/snap/<ai-model-name>/components...`.

The model is now ready to be used.

<!-- links -->
[open-webui-install]: https://github.com/open-webui/open-webui?tab=readme-ov-file#how-to-install-
[open-vino-model-server]: https://docs.openvino.ai/2025/model-server/ovms_what_is_openvino_model_server.html
