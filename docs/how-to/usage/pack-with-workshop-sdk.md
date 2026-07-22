(pack-with-workshop-sdk)=

# Create an inference snap with LLMs

Use this guide when you want to pack an inference snap with a Workshop environment and the `inference-snaps-sdk`.

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

The template repository contains a `Makefile` that you can use to download the model weights and a `README` needed by the agents to run the packaging pipeline. You need to edit these files to provide the correct inputs for your inference snap creation.

Let's start by preparing the `Makefile` to download the required model weights. Sometimes you may want to pack inside the same inference snap multiple models with different number of paramenters. In this case, you can add multiple targets to the `Makefile` to download all the required models. The following is an example of a `Makefile` that downloads two models and their corresponding mmproj files.

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

At this point, you should be in the Workshop shell and ready to run packaging steps.

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

The SDK pipeline walks through the standard packaging flow and generates the build outputs.

## 4. Open a PR for your newly packaged inference snap

At the end of a successful packaging run, the AI agent generates a PR description with the build outputs. You can open a PR to your inference snap repository with the generated description.
