(sideload-model)=
# Sideload a model
Sideloading allows testing a different LLM without changing snap configuration files or rebuilding the snap.

By default, the runtimes point model files that are in snap components. 
This is set internally via environment variables.
The exact environment variable depends on the packaging and engines.

For example, let's sideload a 12B model into gemma4 snap, for the cpu engine:

Inspect current snap environment variables to find the relevant variables for the model file:

```shell
sudo gemma4 run -- env
```

This command prints all environment variables available in the snap.
Look for variables `MODEL_FILE`, `MODEL_NAME`, and `MMPROJ_FILE`.

Then, set environment variables to point to the new model file:

```shell

Set model-specific variables:

```shell
sudo gemma4 set env.MODEL_FILE=/path/to/my_model
sudo gemma4 set env.MODEL_NAME=my_model
sudo gemma4 set env.MMPROJ_FILE=/path/to/my_model.mmproj
```

The `set` command prompts for confirmation and restarts the snap.