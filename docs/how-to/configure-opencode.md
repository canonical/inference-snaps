(configure-opencode)=
# Configure OpenCode for use with inference snaps

[OpenCode][opencode-docs] is an open source AI coding agent. It’s available as a terminal-based interface, desktop app, or IDE extension.
This guide will help you configure OpenCode to use an inference snap as a provider for AI-powered features.

You must have OpenCode and an inference snap installed to continue.

````{tip}

OpenCode is available as a snap.
To install [the snap](https://snapcraft.io/opencode), run:

```shell
sudo snap install opencode --classic
```

````

## Get the inference snap endpoint and model name

Use the inference snap's `status` command to retrieve the OpenAI endpoint URL and model name.

For example, with the `gemma4` snap:

```shell
gemma4 status
```

The output will show the endpoint and model information:

```
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

Note the `openai` endpoint URL and the `model.name` value, as you'll need them for configuration.

## Configure OpenCode

Open or create the OpenCode configuration file at `~/.config/opencode/opencode.jsonc`.

Add a new provider configuration with the endpoint and model information:

```{code-block} json
:emphasize-lines: 3-12  
{
  "$schema": "https://opencode.ai/config.json",
  "provider": {
    "gemma4": {
      "api": "openai",
      "name": "Gemma 4 Inference Snap",
      "options": {
        "baseURL": "http://localhost:8336/v1"
      },
      "models": {
        "gemma4-e4b-q4-k-m": {}
      }
    }
  }
}
```

Update the configuration values to match your specific snap:

- **Provider key** (`gemma4`): A unique identifier for this provider
- **`name`**: A descriptive name for the provider (e.g., "Gemma 4 Inference Snap")
- **`baseURL`**: The OpenAI endpoint URL from the `status` command
- **`models`**: An object containing the model name from the `status` command

## Use the inference snap with OpenCode

Once configured, you can select the inference snap served model within OpenCode.

Use `/models` in OpenCode and select the newly created model, e.g. `gemma4-e4b-q4-k-m Gemma 4 Inference Snap`.

Refer to the [OpenCode documentation][opencode-docs] for details on using AI features.

<!-- links -->
[opencode-docs]: https://opencode.ai/docs
