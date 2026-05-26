(sideload-model)=
# Sideload a model
Sideloading allows testing a different LLM without changing snap configuration files or rebuilding the snap.

By default, the runtimes point model files that are in snap components. 
This is set internally via environment variables.
The exact environment variable depends on the packaging and engines.

For example, let's sideload a 12B model into gemma3 snap, for the cpu engine:

Inspect current snap environment variables to find the relevant variables for the model file:

```shell
sudo gemma3 run -- env
```

This command prints all environment variables available in the snap.
Look for variables `MODEL_FILE`, `MODEL_NAME`, and `MMPROJ_FILE`.

Then, set environment variables to point to the new model file:

```shell
sudo gemma3 set env.model-file=/path/to/model/gemma-3-12b-it-q4_0.gguf
sudo gemma3 set env.model-name=gemma-3-12b-it-q4-0
sudo gemma3 set env.mmproj-file=/path/to/mmproj/mmproj-model-f16-12B.gguf
```

The `set` command prompts for confirmation and restarts the snap.