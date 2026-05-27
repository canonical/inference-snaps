(sideload-model)=
# Sideload a model
Sideloading allows testing a different LLM with an existing snap.

By default, the runtimes point to model files that are in snap components. 
This is set internally via environment variables.
The exact environment variable depends on the packaging and engines.

For example, let's sideload a 12B model into `gemma3` snap, for the `cpu` engine:

Inspect current snap environment variables to find the relevant variables for the model file:

```shell
sudo gemma3 run env
```

This command prints all environment variables available in the snap.
Look for variables `MODEL_FILE`, `MODEL_NAME`, and `MMPROJ_FILE`.

The goal is to override the values of these environment variables with paths that point to model files on the host.
For this example, download the following two model files from Hugging Face into ~/models/gemma3-12b:

```shell
mkdir ~/models/gemma3-12b
cd ~/models/gemma3-12b
wget https://huggingface.co/unsloth/gemma-3-12b-it-qat-GGUF/resolve/main/gemma-3-12b-it-qat-Q4_0.gguf
wget https://huggingface.co/unsloth/gemma-3-12b-it-qat-GGUF/resolve/main/mmproj-F16.gguf
```

Since the model files are in the user home directory, you need to connect the `home` interface to grant the snap access:

```shell
sudo snap connect gemma3:home
```

Finally, set the following environment variables:

```shell
sudo gemma3 set env.model-file=/home/$USER/models/gemma3-12b/gemma-3-12b-it-q4_0.gguf\
env.model-name=gemma-3-12b-it-q4-0\
env.mmproj-file=/home/$USER/models/gemma3-12b/mmproj-F16.gguf  
```

The `set` command prompts for confirmation and restarts the snap.

To verify, query the available models:

```
curl http://localhost:8328/v1/models | jq
```