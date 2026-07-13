(configure-openclaw)=

# Use inference snaps in OpenClaw

OpenClaw is a local AI agent that supports tool calling.
This guide will help you configure OpenClaw to use an inference snap.

Similar steps can be used to configure other AI agents, such as {spellexception}`PicoClaw`, {spellexception}`NemoClaw`, and Hermes.

You must have OpenClaw and an inference snap installed to continue.

## Get the inference snap endpoint and model name

Use the inference snap's `status` command to retrieve the OpenAI endpoint URL and model name.

For example, using the `qwen3` snap:

```shell
qwen3 status
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

## Configure OpenClaw

After installation, OpenClaw launches the interactive setup wizard with a number of sections to set up the agent.
The connection with an inference snap needs to be set up during the model configuration section.

If OpenClaw has already been set up, the wizard can be manually started again to run through the model configuration section:

```shell
openclaw configure --section model
```

The wizard will ask a number of questions.
Answer the questions using this guide as a reference:

- **Model/{spellexception}`auth` provider**: Under the "More…" section, choose "Custom Provider".
- **API Base URL**: Enter the `openai` endpoint URL obtained previously from the `status` command.
- **How do you want to provide this API key?**: Choose "Paste API key now"
- **API Key**: Leave blank
- **Endpoint compatibility**: Choose "OpenAI-compatible".
- **Model ID**: Enter the `model.name` value obtained previously from the `status` command.

The wizard will now run a live completion probe to verify the endpoint.
If it fails, it could be due to insufficient RAM or VRAM, or missing tool-calling support in the model.
Refer to the [OpenClaw troubleshooting guide](https://docs.openclaw.ai/help/troubleshooting) for further help.

- **Endpoint ID**: Any desired value, but we recommend making it clear this is the connection with the inference snap, e.g. `qwen3-snap-endpoint`.
- **Model alias (optional)**: This is the name of the model that will be shown in the OpenClaw dashboard. If not set here, it will fall back to the model ID. We recommend making it clear this is the model from the inference snap, e.g. `qwen3-snap-model`.

The wizard should print `Configure complete`.
You can now use OpenClaw, and it should use the inference snap as its model.

If you have multiple models configured, the OpenClaw dashboard lets you choose the one to use from the dropdown below the chat input box.
