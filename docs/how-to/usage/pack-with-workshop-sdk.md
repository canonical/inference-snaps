(pack-with-workshop-sdk)=

# Pack an inference snap with Workshop and the inference-snaps-sdk

Use this guide when you want to package an inference snap with the Workshop environment and the `inference-snaps-sdk`.

## Before you start

You need:

- Workshop installed on the machine where you will run the packaging pipeline. See [Install Workshop](https://ubuntu.com/workshop/docs/tutorial/part-1-get-started/#) for instructions.
- A GitHub repository for your inference snap (typically created from the [inference snap template](https://github.com/canonical/inference-snap-template))

## 1. Open your snap project

Clone your inference snap repository and navigate into it:

```shell
git clone <your-inference-snap-repo-url>
cd <your-inference-snap-repo>
```

## 2. Inputs preparation

Prepare the Makefile to download the required model weights:
```makefile
SHELL := /bin/bash

.PHONY: download-models setup-hf-cli

all: download-models

download-models: download-model-A download-mmproj-A download-model-B-splitted download-mmproj-B

setup-hf-cli:
	sudo apt-get install -y python3-venv
	python3 -m venv .venv
	. .venv/bin/activate && pip install --upgrade pip && pip install -U huggingface_hub

download-model-4b: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>

download-mmproj-4b: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>

download-model-9b: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir-1-of-N>
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir-2-of-N>
    ...

download-mmproj-9b: setup-hf-cli
	. .venv/bin/activate && hf download <hf-link> --local-dir components/<component-dir>
```

Now open the `README.md`with your favorite editor and compile the section at the top of the file between the `<!--` and `-->` comments. This section is needed by the SDK to generate the packaging pipeline.

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

## 3. Launch Workshop

From your project directory, start Workshop:

```shell
workshop launch
```

Then open a shell inside the Workshop environment:

```shell
workshop shell
```

At this point, you should be in the Workshop shell and ready to run packaging steps.

## 3. Start OpenCode in Workshop

Inside the Workshop shell, start OpenCode:

```shell
opencode --auto
```

This opens the OpenCode TUI with the skills and agents installed by `inference-snaps-sdk`.

## 4. Run the packaging pipeline

In the OpenCode TUI, run this prompt:

```text
start snap packaging pipeline
```

The SDK pipeline walks through the standard packaging flow and generates the build outputs.

## 5. Open a PR for your newly packaged inference snap

On a successful packaging run, the SDK generates a PR description with the build outputs. You can open a PR to your inference snap repository with the generated description.
