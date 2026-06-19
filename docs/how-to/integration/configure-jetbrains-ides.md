(configure-jetbrains-ides)=

# Use inference snaps from JetBrains IDEs

Many inference snaps provide an OpenAI-compliant API that can be integrated with JetBrains IDEs (e.g. IntelliJ IDEA, GoLand, PyCharm, RustRover).
This guide explains how to use the JetBrains AI Assistant plugin and configure an inference snap as a provider.
It assumes that the snap has already been installed and configured.

## Get the inference snap's endpoint

Use the inference snap's `status` command to retrieve the OpenAI endpoint URL.

For example, with the `gemma4` snap, run:

```shell
gemma4 status
```

The output will show the endpoint and model information:

```yaml
engine: nvidia-gpu
services:
    server: active
    server-webui: active
endpoints:
    openai: http://localhost:8336/v1
    webui: http://localhost:8337/
model:
    name: gemma4-e4b-q4-k-m
```

Note down the `openai` endpoint URL as you'll need it for configuration.

## Configure AI Assistant for local models

JetBrains IDEs support _AI Assistant_, which can be configured to use locally hosted models.
Open your JetBrains IDE and follow these steps to connect the inference snap:

### Ensure AI Assistant is installed

AI Assistant comes as a plugin for JetBrains IDEs and may already be installed in your IDE. If not, follow this procedure:

1. Navigate to **Settings** › **Plugins**.

2. Search for and select _JetBrains AI Assistant_.

3. Click the **Install** button.

4. You might need to restart your IDE for the installation to take effect.

For more information, please see JetBrains' page on how to [install AI Assistant][jetbrains-install-ai-assistant].

### Add the OpenAI-compatible endpoint

Once AI Assistant is installed, the API endpoint needs to be configured:

1. Navigate to **Settings** › **Tools** › **AI Assistant** › **Providers & API keys**.

2. In the **Third-party AI providers** section, select **OpenAI-compatible** as the provider type.

3. In the configuration panel:
   - **URL**: Enter the `openai` API endpoint of the installed inference snap (for example, `http://localhost:8336/v1`).
   - **API Key**: Leave empty.

4. Click **Test Connection** to verify the connection is working correctly.

5. Click **Apply** to save the configuration.

For more information, refer to the JetBrains' [using custom models][jetbrains-ai-docs] guide.

### Assign the model to AI Assistant features

Once the endpoint is configured, you can assign specific AI Assistant features to the model:

1. Navigate to **Settings** › **Tools** › **AI Assistant** › **Providers & API keys**.

2. Scroll to the **Models Assignment** section.

3. For each AI Assistant feature group (Core features, Instant helpers, Completion model), select the model from the dropdown.

   - **Core features**: Used for in-editor code generation, commit message generation, and AI Chat
   - **Instant helpers**: Used for lightweight features like chat context collection and name suggestions
   - **Completion model**: Used for inline code completion (requires a Fill-in-the-Middle capable model)

4. Click **Apply** to save your changes.

## Use the inference snap in the IDE

Once the model is configured, you can use it in different parts of the IDE.
For example, to chat:

1. Open AI Chat by selecting **Tools** › **AI Assistant** from the menu.

2. In the chat window, select the model from the model selector dropdown. The configured model should appear under the **OpenAI-compatible** section.

3. Start typing prompts. The requests will be sent locally to the inference snap.

```{tip}
For more information about using AI Assistant, refer to the [JetBrains AI Assistant documentation][jetbrains-ai-chat].
```

<!-- links -->
[jetbrains-ai-docs]: https://www.jetbrains.com/help/ai-assistant/use-custom-models.html
[jetbrains-ai-chat]: https://www.jetbrains.com/help/ai-assistant/ai-chat.html
[jetbrains-install-ai-assistant]: https://www.jetbrains.com/help/ai-assistant/installation-guide-ai-assistant.html
