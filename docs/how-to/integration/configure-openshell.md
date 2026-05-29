(configure-openshell)=
# Use inference snaps in OpenShell

[OpenShell][openshell-docs] is NVIDIA's safe, private runtime for autonomous AI agents. It provides sandboxed execution environments with controlled inference routing, keeping AI traffic private by proxying model requests through a local gateway.
This guide will help you configure OpenShell to use an inference snap as the inference backend for AI agents running inside sandboxes.

You must have OpenShell and an inference snap installed to continue.

````{tip}

OpenShell is available as a snap.
To install [the snap](https://snapcraft.io/openshell), run:

```shell
sudo snap install openshell
```

````

## Network configuration

Since OpenShell's sandboxes live inside a confined environment, inference snaps have to be configured to accept network connections from different IP addresses. To allow that, the `http.host` property needs to be modified.

In this guide we will set it to `0.0.0.0` to allow connections from any IP:

```shell
sudo gemma4 set http.host=0.0.0.0
```

Substitute **`gemma4`** with the name of the inference snap you are using. Ensure the inference snap is restarted to apply the change.

## Get the inference snap endpoint and model name

Use the inference snap's `status` command to retrieve the OpenAI endpoint URL and model name.

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

Note the `openai` endpoint URL and the `model.name` value, as you'll need them for configuration.

## Create an OpenShell provider

Create a provider record that points OpenShell to the inference snap's OpenAI-compatible endpoint.

```shell
openshell provider create \
    --name gemma4-snap \
    --type openai \
    --credential OPENAI_API_KEY=not-needed \
    --config OPENAI_BASE_URL=http://$(hostname).local:8336/v1
```

Update the configuration values to match your specific snap:

- **`--name`**: A unique identifier for this provider (e.g., `gemma4-snap`)
- **`--credential OPENAI_API_KEY`**: A placeholder value; the inference snap does not require authentication
- **`--config OPENAI_BASE_URL`**: The inference snap's OpenAI endpoint, where the host is replaced with your machine's hostname, obtained by concatenating the output of `hostname` with the string `.local`

```{note}
OpenShell sandboxes cannot reach `localhost` or `127.0.0.1` on the host directly.
Use your machine's **hostname** when building the provider URL, or use the host machine's LAN IP address.
```

## Set inference routing

Point OpenShell's `inference.local` endpoint at the provider you created:

```shell
openshell inference set \
    --provider gemma4-snap \
    --model gemma4-e4b-q4-k-m
```

Make sure to substitute `gemma4-snap` with the name you've chosen in the previous step, and that the model name matches the one obtained with the `status` command.

## Use the inference snap from a sandbox

Once configured, any agent running inside an OpenShell sandbox can reach the inference snap through the `https://inference.local` endpoint.
OpenShell injects the provider credentials and rewrites the model before forwarding each request, so no API key is needed inside the sandbox.

Launch a new sandbox:

```shell
openshell sandbox create
```

Within the sandbox, make a completion request to the inference snap:

```shell
curl https://inference.local/v1/responses \
    -H "Content-Type: application/json" \
    -d '{
      "instructions": "You are a helpful assistant.",
      "input": "Hello!"
    }'
```

Refer to the [OpenShell documentation][openshell-docs] for details on creating sandboxes and running agents.

<!-- links -->
[openshell-docs]: https://docs.nvidia.com/openshell/latest/home
