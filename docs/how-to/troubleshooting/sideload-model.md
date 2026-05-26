(sideload-model)=
# Sideload a model
Sideloading allows testing a different LLM without changing snap configuration files or rebuilding the snap.

Inspect current snap environment variables:

```shell
sudo <inference-snap> run -- env
```

This command prints all environment variables available in the snap.

Set model-specific variables:

```shell
sudo <inference-snap> set env.MODEL_FILE=/path/to/my_model
sudo <inference-snap> set env.MODEL_NAME=my_model
```

If required by the model, set `MMPROJ_FILE`.

```shell
sudo <inference-snap> set env.MMPROJ_FILE=/path/to/my_model.mmproj
```

Restart the snap to apply changes:

```shell
sudo snap restart <inference-snap>
```