(parallel-install)=

# Serve multiple models from one inference snap in parallel

Snapd's parallel installs let you install the same snap several times on one machine, each as an independent instance with its own data, configuration, and services.
This is useful when you want to serve several different models from the same inference snap at the same time, each on its own endpoint.

This how-to uses the `gemma4` snap to run three instances side by side, serving the `e2b`, `e4b`, and `12b` models.
Each instance is named after the model size it serves, so the model choice is reflected in the instance name and its CLI.

## Prerequisites

Parallel installs are an experimental snapd feature that you must enable once:

```shell
sudo snap set system experimental.parallel-instances=true
```

Each instance is identified by appending an instance key to the snap name, separated by an underscore (for example `gemma4_e2b`).
The instance key can contain letters and digits.

## Install three instances

Install three instances, naming each one after the model size it will serve:

```shell
sudo snap install gemma4_e2b gemma4_e4b gemma4_12b
```

Each instance exposes its own CLI, named after the instance.
Use the corresponding CLI to manage each instance:

```shell
gemma4_e2b status
gemma4_e4b status
gemma4_12b status
```

## Assign distinct ports to each instance

Every instance exposes an API port and a web UI port, which default to the same values (`8336` and `8337` for `gemma4`, see {ref}`network-ports`), so the instances would conflict.
Keep `gemma4_e2b` on its default ports `8336` and `8337`, and give the other instances unique free port pairs:

```shell
sudo gemma4_e4b set http.port=8436 webui.http.port=8437
sudo gemma4_12b set http.port=8536 webui.http.port=8537
```

When prompted, type `Y` to confirm the restart of each instance.

## Change the engine

Some engines may not support every model size.
To use the `e2b`, `e4b`, and `12b` models in this example, select the `cpu` engine for all three instances.

List the engines available on your host:

```shell
gemma4_12b list-engines
```

```shell
sudo gemma4_e2b use-engine cpu
sudo gemma4_e4b use-engine cpu
sudo gemma4_12b use-engine cpu
```

When prompted, type `Y` to confirm the installation of any missing components and to confirm the restart of each instance.

For more details, see {ref}`switch-between-engines`.

## Select a model for each instance

Use `list-models` to see the model IDs supported by the snap:

```shell
gemma4_e2b list-models
```

Select the model matching each instance's size:

```shell
sudo gemma4_e2b use-model gemma4-e2b 
sudo gemma4_e4b use-model gemma4-e4b
sudo gemma4_12b use-model gemma4-12b
```

As before, when prompted, type `Y` to confirm.

## Verify the instances

Query each endpoint to confirm the served model:

```console
$ curl http://127.0.0.1:8336/v1/models | jq '.data[].id'
"gemma4-e2b"

$ curl http://127.0.0.1:8436/v1/models | jq '.data[].id'
"gemma4-e4b"

$ curl http://127.0.0.1:8536/v1/models | jq '.data[].id'
"gemma4-12b"
```

Each instance now serves a different model on its own API and web UI ports.

## Remove an instance

Remove a single instance without affecting the others by using its full instance name:

```shell
sudo snap remove gemma4_12b
```
