(parallel-install)=

# Add a parallel Gemma 4 instance

Assume `gemma4` is already installed and serving the `gemma4-e4b` model on its default API port, `8336`. To serve `gemma4-12b` at the same time, install a second, independent instance named `gemma4_12b` and give it a different pair of ports. The existing `gemma4` instance keeps its model, configuration, and endpoints.

## Enable parallel installs

Snapd's parallel installs are experimental. Enable them once on the machine if they are not already enabled:

```shell
sudo snap set system experimental.parallel-instances=true
```

## Install the second instance

Install another instance of the same snap, using an instance key to distinguish it from the existing `gemma4`:

```shell
sudo snap install gemma4_12b
```

The new instance has its own CLI and configuration. Use `gemma4_12b` for commands that should affect the new instance; continue using `gemma4` to manage the existing e4b service.

## Assign unused ports

The existing `gemma4` instance uses API port `8336` and web UI port `8337` by default (see {ref}`network-ports`). Assign a free pair to the new instance so the services do not conflict:

```shell
sudo gemma4_12b set http.port=8436 webui.http.port=8437
```

When prompted, confirm the restart. If those ports are in use, choose another free pair.

## Select a compatible engine and model

Check which engines are available for the new instance:

```shell
gemma4_12b list-engines
```

Choose an engine compatible with the 12b model and your system. For example, to use the CPU engine:

```shell
sudo gemma4_12b use-engine cpu
```

Confirm any prompts to install missing components and restart the instance. For more details, see {ref}`switch-between-engines`.

List the model IDs supported by the snap and select the 12b model:

```shell
gemma4_12b list-models
sudo gemma4_12b use-model gemma4-12b
```

Confirm the prompt to restart the instance.

## Verify both instances

Check the existing e4b endpoint and the new 12b endpoint:

```console
$ curl http://127.0.0.1:8336/v1/models | jq '.data[].id'
"gemma4-e4b"

$ curl http://127.0.0.1:8436/v1/models | jq '.data[].id'
"gemma4-12b"
```

The original `gemma4` instance continues serving e4b on its original ports, while `gemma4_12b` serves 12b on the new ports.