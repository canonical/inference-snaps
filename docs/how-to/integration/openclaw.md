(openclaw)=

# Use inference snaps in OpenClaw

Inference snaps that support tool calling can be used with local agents like OpenClaw, {spellexception}`PicoClaw`, {spellexception}`NemoClaw`, Hermes, etc.
The steps below are a guideline on how to configure OpenClaw to use an inference snap like Gemma 4.
Similar steps can be used to configure other AI agents.

## Look up base URL and model ID

Run the respective inference snap's status command to look up the `openai` endpoint URL and the model name.
These values will be required when configuring the agent.

```shell
$ <inference-snap> status
engine: cpu
services:
    server: active
    server-webui: active
endpoints:
    openai: http://127.0.0.1:8336/v3
    webui: http://127.0.0.1:8337/
model:
    name: gemma4-e4b-it-int4-ov
```

## Configure using wizard

After installation, OpenClaw launches the interactive setup wizard with a number of sections to set up the agent.
The connection with an inference snap needs to be set up during the model configuration section.

If OpenClaw has already been set up, the wizard can be manually started again to run through the model configuration section:

```shell
openclaw configure --section model
```

The wizard will ask a number of questions.
Answer the questions using this guide as a reference.

Model/{spellexception}`auth` provider
~ Choose "Custom Provider"

API Base URL
~ Enter the `openai` endpoint URL obtained previously from the `status` command.

How do you want to provide this API key?
~ Choose "Paste API key now"

API Key
~ Leave blank

Endpoint compatibility
~ Choose either the legacy "OpenAI-compatible" or the newer "OpenAI Responses-compatible".

Model ID
~ Enter the model name obtained previously from the `status` command.

The endpoint will now be verified by the wizard.
If it fails, but the `<inference-snap> chat` command works, it is probably because the model does not have enough resources to be used by an agent.

Endpoint ID and Model alias
~ You can give the endpoint and model alias any value of your choosing. We recommend making it clear this is the model from the inference snap.

The wizard should print `Configure complete`.
You can now use OpenClaw, and it should use the inference snap as model.

If you have multiple models configured, the OpenClaw dashboard lets you choose the one to use from the dropdown below the chat input box.

## Configure using config file

OpenClaw's configuration is persisted to a file. One can also manually edit this file.

```
nano ~/.openclaw/openclaw.json
```

After the wizard was run, the file will contain a section listing the available models.
You can manually edit this file to change the configuration or add new models.

```
{
...
  "models": {
    "mode": "merge",
    "providers": {
      "qwen3-snap": {
        "baseUrl": "http://localhost:8338/v1",
        "api": "openai-responses",
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
...
}
```

After modifying the file, restart the OpenClaw gateway to apply the changes:

```
openclaw gateway restart
```
