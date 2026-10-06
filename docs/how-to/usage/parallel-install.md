(parallel-install)=

# Serve multiple model sizes in parallel

Each inference snap can serve a single model. You need to install a second instance of the snap to serve a second model, such as a different model size or quantization. This method is also useful if you want to serve models on different silicons.

## Enable parallel installs

Enable [parallel installs](https://snapcraft.io/docs/explanation/how-snaps-work/parallel-installs/) to allow installation of multiple instances of the same snap:

```shell
sudo snap set system experimental.parallel-instances=true
```

## Install the second instance

Install the second instance of the same snap, using an instance key to distinguish it from the existing installation. The instance key is an arbitrary string added as a suffix to the snap name. For example, if [Gemma4](https://snapcraft.io/gemma4) is already installed and using the E4B model, and you want a second instance to serve the 12B model, use a `_12b`:

```shell
sudo snap install gemma4_12b
```

The new instance comes with its own CLI and configuration. Use `gemma4_12b` for commands that should affect the new instance; continue using `gemma4` to manage the existing instance with the E4B model.

## Assign unused ports

Configure the new instance to use available TCP ports to serve the API and Web UI.

The Gemma4 snap uses API port `8336` and Web UI port `8337` by default (see {ref}`network-ports`). Assign other ports to the new instance so the services do not conflict:

```shell
sudo gemma4_12b set http.port=8436 webui.http.port=8437
```

When prompted, confirm the restart. If those ports are in use, choose another free pair.

## Select a compatible engine and model

Check which engines are available for the new instance:

```shell
gemma4_12b engines
```

Choose an engine compatible on your system. This could be different from the one used on the first instance. For example, to use the CPU engine:  

```shell
sudo gemma4_12b use-engine cpu
```

Confirm any prompts to install missing components and restart the instance. For more details, see {ref}`switch-between-engines`.

List the models supported by the current engine:  

```shell
gemma4_12b models
```

Check that the 12B model is available. If it is, select it:

```shell
sudo gemma4_12b use-model gemma4-12b
```

Confirm the prompt to restart the instance.

## Verify both instances

Verify served models by querying the `/v1/models` endpoints:  

```console
$ curl http://127.0.0.1:8336/v1/models | jq '.data[].id'
"gemma4-e4b"

$ curl http://127.0.0.1:8436/v1/models | jq '.data[].id'
"gemma4-12b"
```

The original `gemma4` instance continues serving E4B on its original ports, while `gemma4_12b` serves 12B on the new ports.  
