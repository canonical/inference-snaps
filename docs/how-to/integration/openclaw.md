(configure-openclaw)=

# Use inference snaps in OpenClaw

OpenClaw is a local AI agent that supports tool calling.
This guide will help you configure OpenClaw to use an inference snap.

Similar steps can be used to configure other AI agents, such as {spellexception}`PicoClaw`, {spellexception}`NemoClaw`, and Hermes.

You must have OpenClaw and an inference snap installed to continue.

## Get the inference snap endpoint and model name

Use the inference snap's `status` command to retrieve the OpenAI endpoint URL and model name.

For example, run:

```shell
<inference-snap> status
```

The output will show the endpoint and model information:

```yaml
engine: cpu
services:
    server: active
    server-webui: active
endpoints:
    openai: http://127.0.0.1:8338/v1
    webui: http://127.0.0.1:8339/
model:
    name: qwen3-8b-q4-k-m
```

Note down the `openai` endpoint URL and the `model.name` value, as you'll need them for configuration.

## Configure using the wizard

After installation, OpenClaw launches the interactive setup wizard with a number of sections to set up the agent.
The connection with an inference snap needs to be set up during the model configuration section.

If OpenClaw has already been set up, the wizard can be manually started again to run through the model configuration section:

```shell
openclaw configure --section model
```

The wizard will ask a number of questions.
Answer the questions using this guide as a reference:

- **Model/{spellexception}`auth` provider**: Choose "Custom Provider"
- **API Base URL**: Enter the `openai` endpoint URL obtained previously from the `status` command.
- **How do you want to provide this API key?**: Choose "Paste API key now"
- **API Key**: Leave blank
- **Endpoint compatibility**: Choose "OpenAI-compatible".
- **Model ID**: Enter the `model.name` value obtained previously from the `status` command.

The endpoint will now be verified by the wizard.
If it fails, but the `<inference-snap> chat` command works, it is probably because the model does not have enough resources to be used by an agent.

- **Endpoint ID and Model alias**: You can give the endpoint and model alias any value of your choosing. We recommend making it clear this is the model from the inference snap.

The wizard should print `Configure complete`.
You can now use OpenClaw, and it should use the inference snap as its model.

If you have multiple models configured, the OpenClaw dashboard lets you choose the one to use from the dropdown below the chat input box.

## Configure using the configuration file

OpenClaw's configuration is persisted to a file. One can also manually edit this file.

```shell
nano ~/.openclaw/openclaw.json
```

After the wizard was run, the file will contain a section listing the available models.
For example, a configuration file with both the qwen3 and gemma4 inference snaps configured would look like this:

```json
{
  // other configuration sections...
  "models": {
    "mode": "merge",
    "providers": {
      "qwen3-snap": {
        "baseUrl": "http://127.0.0.1:8338/v1",
        "api": "openai-completions",
        "models": [
          {
            "id": "qwen3-8b-q4-k-m",
            "name": "qwen3-8b-q4-k-m (Custom Provider)",
            "input": [
              "text"
            ]
          }
        ]
      },
      "gemma4-snap": {
        "baseUrl": "http://127.0.0.1:8336/v3",
        "api": "openai-completions",
        "models": [
          {
            "id": "gemma4-e4b-it-int4-ov",
            "name": "gemma4-e4b-it-int4-ov (Custom Provider)",
            "input": [
              "text",
              "image"
            ]
          }
        ]
      }
    }
  },
  // other configuration sections...
}
```

After modifying the file, restart the OpenClaw gateway to apply the changes:

```shell
openclaw gateway restart
```
