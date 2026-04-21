(configure-openclaw)=
# Configure OpenClaw for use with inference snaps

Inference snaps that support tool calling can be used with local agents like OpenClaw, PicoClaw, NemoClaw, etc.
Refer to {ref}`available snaps <available-snaps>` to see which inference snaps support tool calling.

## Look up base URL and model ID

The inference endpoints that an inference snap exposes can be seen using the `status` command.
Agents normally use the `openai` endpoint, so make sure the snap you are using has this endpoint available.

Run the status command and note the URL printed next to `openai`.
This is the base URL that will be required when configuring the agent.

```shell
$ <inference-snap> status
engine: intel-gpu
services:
    server: active
    server-webui: active
endpoints:
    kserve: http://localhost:8328/v2
    openai: http://localhost:8328/v3
    tensorflow-serving: http://localhost:8328/v1
    webui: http://localhost:8329/
```

We also require the model ID.
This can be looked up from the `/models` path served under the `openai` endpoint.
You can use `curl` and `jq` to do this in a single command, or go to `<base-url>/models` in your browser.

```shell
$ curl http://localhost:8328/v3/models 2>/dev/null | jq -r '.data[0].id'
<model-id>
```

## Add model endpoint to OpenClaw

After installation, OpenClaw launches the interactive setup wizard.
This guides you through a number of sections to set up the agent.
To use an inference snap, this needs to be set up during the model configuration section.

If OpenClaw has already been set up, the wizard can be manually started again to run through the model configuration
section:

```shell
openclaw configure --section model
```

The wizard will ask a number of questions.
These questions can change over time, so we list the key answers here to help you answer them correctly.

Where will the Gateway run?
~ Most likely "Local (this machine)", but it depends on your setup

Model/auth provider
~ Choose "Custom Provider" 

API Base URL
~ Enter the `openai` base URL obtained previously in this how-to

How do you want to provide this API key?
~ Choose "Paste API key now"

API Key
~ Leave blank

Endpoint compatibility
~ Choose "OpenAI-compatible"

Model ID
~ Enter the model ID obtained previously in this how-to

The endpoint will now be verified by the wizard.
If it fails, but the `<inference-snap> chat` command works, it is probably because the model is too slow for OpenClaw to
use.

Endpoint ID and Model alias
~ You can give the endpoint and model alias any value of your choosing. We recommend making it clear this is the model from the inference snap

The wizard should print `Configure complete`.
You can now use OpenClaw, and it should use the inference snap as model.

If you have multiple models configured, you can select the one to use from a dropdown on the OpenClaw dashboard.
