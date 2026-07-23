(pack-with-workshop-sdk)=

# Create an inference snap with LLMs

Use this guide when you want to pack an inference snap with a Workshop environment and the `inference-snaps-sdk`.

```{admonition} Experimental workflow
This packaging workflow is experimental. The steps, prompts, and generated outputs may change, and you may need to make manual adjustments before opening a pull request.
```

## Before you start

You need:

- Workshop installed on the machine where you will run the packaging pipeline. See [Install Workshop](https://ubuntu.com/workshop/docs/tutorial/part-1-get-started/#) for instructions.
- A GitHub repository for your inference snap (typically created from the [inference snap template](https://github.com/canonical/inference-snap-template))

Helpful references while you follow this tutorial:

- [OpenCode integration](../how-to/integration/opencode.md)
- [Network ports registry](../reference/network-ports.md)
- [Troubleshooting](../how-to/troubleshooting/index.md)

## 1. Open your snap project

Clone your inference snap repository and navigate into it:

```shell
git clone <your-inference-snap-repo-url>
cd <your-inference-snap-repo>
```

## 2. Inputs preparation

The template repository contains a `Makefile` that you can use to download the model weights and a `README` needed by the agents to run the packaging pipeline. You need to edit these files to provide the correct inputs for your inference snap creation.

Let's start by preparing the `Makefile` to download the required model weights. If you want to package multiple model variants in one snap, add multiple targets to download each model and related projector files.

The following example downloads two models and their corresponding projectors:

```makefile
SHELL := /bin/bash

.PHONY: download-models setup-hf-cli

all: download-models

download-models: download-model-A download-mmproj-A download-model-B-splitted download-mmproj-B

setup-hf-cli:
	sudo apt-get install -y python3-venv
	python3 -m venv .venv
	. .venv/bin/activate && pip install --upgrade pip && pip install -U huggingface_hub

download-model-A: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>

download-mmproj-A: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>

download-model-B: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir-1-of-N>
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir-2-of-N>
    ...

download-mmproj-B: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>
```

Now open the `README.md` with your favorite editor and compile the section at the top of the file between the `<!--` and `-->` comments

```
# This is the name of the snap. The name that is registered on the snap store and also the name of the cli command.
snap-name: e.g. gemma4
# This name is just a friendly name for the snap, it can be used in the documentation
snap-friendly-name: e.g. Gemma 4B
# URL to model card from the model publisher
model-card: e.g. https://ai.google.dev/gemma/docs/core/model_card_4
# The port that the inference snap will use for its API server.
http-port: e.g. 8080
# The port that the inference snap will use for its webui server.
webui-http-port: e.g. 8081
# Optimizations
engines: e.g. cpu, nvidia-gpu
```
Make sure your selected ports do not conflict with entries in [Network ports registry](../reference/network-ports.md).

## 3. Packaging the inference snap

Now that you have prepared the inputs, you can launch Workshop and run the agentic packaging pipeline.

### Launch Workshop

From your project directory, start Workshop:

```shell
workshop launch
```

Then open a shell inside the Workshop environment:

```shell
workshop shell
```

At this point, you should be in the Workshop shell and ready to run packaging steps. You can check that you effectively entered the Workshop shell by running:

```{terminal}
workshop@dev:/project$ whoami
workshop
```


### Start OpenCode in Workshop

Inside the Workshop shell, start OpenCode:

```shell
opencode --auto
```

This opens the OpenCode TUI with the skills and agents installed by `inference-snaps-sdk`. The `--auto` flag raises the permission prompt automatically, so you don't have to manually approve the permissions. Workshop environment is sandboxed so you can safely approve the permissions without worrying about security issues. Nonetheless, if you want to approve the permissions manually, you can run `opencode` without the `--auto` flag.

### Run the packaging pipeline

In the OpenCode TUI, run this prompt:

```text
start snap packaging pipeline
```

The SDK pipeline will validate inputs, summarize them and ask you to confirm before starting the packaging process. The pipeline will then run the packaging steps and generate outputs. 
Expected result:
- The run finishes without blocking errors.
- Build outputs and a generated PR description are available.

## 4. Test your inference snap
First of all let's exit the Workshop shell:

```{terminal}
workshop@dev:/project$ exit
exit
user@host-machine:~$ whoami
user
```

Now navigate to your project directory and run the following command to install the snap:

```shell
cd <your-inference-snap-repo>
sudo snap install *.snap *.comp --dangerous
```

Once the snap is installed, you can run the chat command of the inference snap and test it:

```shell
<snap-name> chat
```

Or open the webui of the inference snap in your browser. At the beginning of this guide you choose two different ports for the API server and the webui server. If you want to check them run:

```shell
<snap-name> get
```

You will see an output similar to this:

```text
http.host: 127.0.0.1
http.port: 8330
...
webui.http.host: 127.0.0.1
webui.http.port: 8331
```

Open with your favorite browser the following URL:
```
http://127.0.0.1:<webui-http-port>
```
or 

```
http://127.0.0.1:<http-port>
```

## 5. Open a PR for your newly packaged inference snap

At the end of a successful packaging run, the AI agent generates a PR description with the build outputs. You can open a PR to your inference snap repository with the generated description.
