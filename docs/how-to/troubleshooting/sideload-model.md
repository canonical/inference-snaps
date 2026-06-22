(sideload-model)=

# Sideload a model

Sideloading allows testing a different LLM with an existing snap.

By default, the runtimes point to model files that are in snap components. 
This is set internally via environment variables.
The exact environment variable depends on the packaging and engines.

For example, let's sideload a 12B gemma4 model fine tuned for coding and agentic work  into `gemma4` snap, for the `cpu` engine.

If not installed, install the snap:

```shell
sudo snap install gemma4
```

If the snap is already installed, refresh it to ensure that it is the latest version:

```shell
sudo snap refresh gemma4 --channel latest/stable
```

Inspect current snap environment variables to find the relevant variables for the model file:

```shell
sudo gemma4 run env
```

This command prints all environment variables available in the snap.
Look for variables `MODEL_FILE`, `MODEL_NAME`, and `MMPROJ_FILE`.

The goal is to override the values of these environment variables with paths that point to model files on the host.
For this example, download the following two model files from Hugging Face into ~/models/gemma4-12b:

```shell
mkdir ~/models/gemma4-12b
cd ~/models/gemma4-12b
wget https://huggingface.co/yuxinlu1/gemma-4-12B-agentic-fable5-composer2.5-v2-3.5x-tau2-GGUF/resolve/main/gemma4-v2-Q4_K_M.gguf
wget https://huggingface.co/unsloth/gemma-4-12b-it-GGUF/resolve/main/mmproj-F16.gguf
```

Since the model files are in the user home directory, you need to connect the `home` interface to grant the snap access:

```shell
sudo snap connect gemma4:home
```

Finally, set the following environment variables:

```shell
sudo gemma4 set \
  env.model-file=/home/$USER/models/gemma4-12b/gemma4-v2-Q4_K_M.gguf \
  env.model-name=gemma4-v2-Q4_K_M \
  env.mmproj-file=/home/$USER/models/gemma4-12b/mmproj-F16.gguf  
```

The `set` command prompts for confirmation and restarts the snap.

To verify, query the available models:

```shell
curl http://localhost:8336/v1/models | jq
```
