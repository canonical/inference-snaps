(pack-with-workshop-sdk)=

# Pack an inference snap with Workshop

Use this guide when you want to package an inference snap with Workshop and the `inference-snaps-sdk`.

```{admonition} Experimental workflow
This packaging workflow is experimental. The steps, prompts, and generated outputs may change, and you may need to make manual adjustments before opening a pull request.
```

## What this tutorial does

In this workflow, you start from an inference snap repository created from the template. That repository already contains the files that the packaging pipeline expects:

- `workshop.yaml` defines the Workshop environment.
- `Makefile` defines how to fetch the model files that will be packaged.
- `README.md` contains the metadata that the SDK agents read before generating packaging changes.

You then launch Workshop, start OpenCode inside that environment, and ask the `inference-snaps-sdk` to run the packaging pipeline. The pipeline reads the inputs from your repository, generates or updates packaging files, builds the snap artifacts, and prepares a pull request description for review.

You do not need to understand every packaging file before you start. For a first run, the important idea is simpler:

- you tell the pipeline which model files to download
- you tell it a few basic metadata values such as the snap name and ports
- the pipeline generates the packaging changes for you and builds the snap artifacts

This tutorial focuses on that first successful run.

## Before you start

Some of the concepts and steps in this tutorial require familiarity with snap packaging and inference snaps. If you are new to these topics, we recommend that you first read the following tutorials:
- [Create an inference snap](../tutorial/create-inference-snap.md)
- [Get started with Qwen VL](../tutorial/qwen-vl-tutorial.md)

For this tutorial, you need:

- Workshop installed on the machine where you will run the packaging pipeline. If you are not familiar with Workshop, see [Install Workshop](https://ubuntu.com/workshop/docs/tutorial/part-1-get-started/#) for instructions.
- A Linux machine where you can install and run snaps.
- A GitHub repository for your inference snap, typically created from the [inference snap template](https://github.com/canonical/inference-snap-template). The template repository already contains the `workshop.yaml`, `Makefile`, and `README.md` files used in this tutorial.
- Access to the model files that you want to package. If your model provider requires authentication, make sure you can download those files before starting the pipeline.

If you are completely new to inference snaps, use a simple first target:

- one model variant
- one engine, usually `cpu`
- no optional Web UI unless you already know that you need it

That keeps the first packaging run easier to validate.

The following references may be helpful while you follow this tutorial:

- [OpenCode integration](../how-to/integration/opencode.md)
- [Network ports registry](../reference/network-ports.md)
- [Troubleshooting](../how-to/troubleshooting/index.md)

## 1. Open your snap project

Clone your inference snap repository and navigate into it:

```shell
git clone <your-inference-snap-repo-url>
cd <your-inference-snap-repo>
```

You should see the following files in your project directory:

```text
Makefile
README.md
workshop.yaml
```

These files are the minimum inputs for the packaging process, and you will update two of them in the next step.

## 2. Prepare the inputs

The template repository contains a `Makefile` that you can use to download the model files and a `README.md` file that the SDK agents read before they generate packaging changes. In this section, you edit those two files so the pipeline knows what to build.

Start by preparing the `Makefile` to download the required model files. For a first attempt, use the smallest working example that matches your model. If you later want to package multiple model variants in one snap, you can add more download targets.

The following example downloads a single GGUF model file:

```makefile
SHELL := /bin/bash

.PHONY: download-models setup-hf-cli

all: download-models

download-models: download-model

setup-hf-cli:
	sudo apt-get install -y python3-venv
	python3 -m venv .venv
	. .venv/bin/activate && pip install --upgrade pip && pip install -U huggingface_hub

download-model: setup-hf-cli
	. .venv/bin/activate && hf download hf://unsloth/Qwen3.5-4B-GGUF/Qwen3.5-4B-UD-Q4_K_XL.gguf --local-dir components/model-q4-k-xl-gguf
```

This `Makefile` sets up the Hugging Face CLI in a local virtual environment and downloads one model file into the `components/` directory.

The `components/` directory is where the packaging pipeline expects model or runtime files that will become part of the final snap.

If your model is multimodal, you may also need one or more projector files. Projector files are extra model files used by some multimodal models, for example models that accept both text and images. If you want to package multiple variants, add extra download targets for each model file and any required projector files.

Now open `README.md` and fill in the metadata block at the top of the file, between the `<!--` and `-->` comments:

```
# This is the name of the snap. The name that is registered on the snap store and also the name of the cli command.
snap-name: e.g. qwen3-5
# This name is just a friendly name for the snap, it can be used in the documentation
snap-friendly-name: e.g. Qwen 3.5
# URL to model card from the model publisher
model-card: e.g. https://qwen.ai/blog?id=qwen3.5
# The port that the inference snap will use for its API server.
http-port: e.g. 8352
# The port that the inference snap will use for its webui server.
webui-http-port: e.g. 8353
# Optimizations
engines: e.g. cpu, nvidia-gpu
```

When you edit this block:

- Replace each `e.g.` value with the real value for your snap.
- Keep the same field names.
- Make sure `snap-name` does not contain capital letters, spaces, or special characters other than hyphens.
- Make sure your selected ports do not conflict with entries in [Network ports registry](../reference/network-ports.md).

For a first run, these fields usually mean:

- `snap-name`: the command name users will run after installation
- `snap-friendly-name`: a human-readable name for docs and descriptions
- `model-card`: the upstream page that describes the model
- `http-port`: the port used by the inference API server
- `webui-http-port`: the port used by the optional Web UI, if your snap exposes one
- `engines`: the runtime families that you want to package. If you are unsure, start with `cpu`.

Before moving on, make sure both inputs are ready:

- the `Makefile` downloads the model files you want
- the `README.md` metadata block contains real values instead of placeholders

## 3. Packaging the inference snap

Now that you have prepared the inputs, you can launch Workshop and run the packaging pipeline.

### Launch Workshop

From your project directory, start Workshop:

```shell
workshop launch
```

This command creates or starts the Workshop environment described by `workshop.yaml`.

Then open a shell inside the Workshop environment:

```shell
workshop shell
```

At this point, you should be in the Workshop shell and ready to run packaging steps. You can check that you effectively entered the Workshop shell by running:

```{terminal}
workshop@dev:/project$ whoami
workshop
```

If you do not see `workshop` as the output of `whoami`, do not continue yet. Exit and check that Workshop launched successfully.

### Start OpenCode in Workshop

The `inference-snaps-sdk` ships with OpenCode, a TUI that allows you to use LLM agents. If you want more information about OpenCode, see [OpenCode documentation](https://opencode.ai/docs).

Inside the Workshop shell, start OpenCode:

```shell
opencode --auto
```

This opens the OpenCode TUI with the skills and agents installed by `inference-snaps-sdk`. The `--auto` flag approves the permission prompts automatically, so you do not have to confirm them one by one. The Workshop environment is sandboxed, so this is usually acceptable for this workflow. If you prefer to approve permissions manually, run `opencode` without the `--auto` flag.

### Run the packaging pipeline

In the OpenCode TUI, run this prompt:

```text
start snap packaging pipeline
```

The SDK pipeline first validates your inputs, then summarizes what it found, and asks you to confirm before it starts the packaging process. After confirmation, it generates or updates packaging files, builds the snap artifacts, and prepares a PR description.

If you are new to this workflow, pay attention to the summary before you confirm. In particular, check that it matches:

- the snap name you set in `README.md`
- the model files downloaded by the `Makefile`
- the engines you want to package

Expected result:

- The run finishes without blocking errors.
- One or more build artifacts such as `.snap` and `.comp` files are available in the repository.
- A generated PR description is available for you to review.

If the pipeline stops on a validation error, fix the reported problem first and then ask the agent to continue the pipeline. Typical beginner issues here are placeholder values left in `README.md`, missing model files, or a download path that does not match the expected `components/` layout.

## 4. Test your inference snap

First, exit the Workshop shell:

```{terminal}
workshop@dev:/project$ exit
exit
user@host-machine:~$ whoami
user
```

Now navigate to your project directory and install the generated artifacts:

```shell
cd <your-inference-snap-repo>
sudo snap install *.snap *.comp --dangerous
```

If this command fails because the shell does not match any `.snap` or `.comp` files, go back to the previous step and confirm that the packaging run actually produced artifacts.

Once the snap is installed, run the snap command and confirm that it starts correctly:

```shell
<snap-name> chat
```

For a beginner, a good first check is simply that the command starts and does not fail immediately with a missing component or configuration error.

If your snap also exposes a Web UI, check which ports were configured:

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

Open the Web UI in your browser at:

```
http://127.0.0.1:<webui-http-port>
```

or 

```
http://127.0.0.1:<http-port>
```

At this point, you should have confirmed that:

- the snap installs successfully
- the CLI starts without immediate errors
- the configured ports match the values you set in `README.md`

## 5. Open a PR for your newly packaged inference snap

At the end of a successful packaging run, the AI agent generates a PR description with the build outputs.

Before opening a pull request, review the repository and confirm that:

- the generated packaging files match the model and engines you intended to package. You can check this by running `<snap-name> list-engines` and `<snap-name> list-models` and comparing the output with the values you set in `README.md` and the files downloaded by the `Makefile`
- the build artifacts were produced successfully
- the snap installs and starts locally

You can then open a pull request to your inference snap repository using the generated PR description as the starting point.
