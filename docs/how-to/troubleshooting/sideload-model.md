(sideload-model)=

# Sideload a model

Sideloading allows testing a different LLM with an existing snap.

By default, the runtimes point to model files that are in snap components. 
This is set internally via environment variables.
The exact environment variable depends on the packaging and engines.

For example, let's sideload a 12B model into `gemma4` snap, for the `cpu` engine.

If not installed, install the snap:

```shell
sudo snap install gemma4
```

If the snap is already installed, refresh it to ensure that it is the latest version:

```shell
sudo snap refresh gemma4 --channel latest/stable```

Inspect current snap environment variables to find the relevant variables for the model file:

```shell
sudo gemma4 run env
```

This command prints all environment variables available in the snap.
Look for variables `MODEL_FILE`, `MODEL_NAME`, and `MMPROJ_FILE`.

The goal is to override the values of these environment variables with paths that point to model files on the host.
For this example, download the following two model files from Hugging Face into ~/models/gemma4-12b:

```shell
mkdir ~/models/gemma4-e4b
cd ~/models/gemma4-e4b
wget https://huggingface.co/unsloth/gemma-4-E4B-it-GGUF/resolve/main/gemma-4-E4B-it-Q4_K_S.gguf
wget https://huggingface.co/unsloth/gemma-4-E4B-it-GGUF/resolve/main/mmproj-BF16.gguf
```

Since the model files are in the user home directory, you need to connect the `home` interface to grant the snap access:

```shell
sudo snap connect gemma4:home
```

Finally, set the following environment variables:

```shell
sudo gemma4 set \
  env.model-file=/home/$USER/models/gemma4-e4b/gemma-4-E4B-it-Q4_K_S.gguf \
  env.model-name=gemma-4-E4B-it-Q4_K_S \
  env.mmproj-file=/home/$USER/models/gemma4-e4b/mmproj-BF16.gguf  
```

The `set` command prompts for confirmation and restarts the snap.

To verify, query the available models:

```
curl http://localhost:8336/v1/models | jq
```
