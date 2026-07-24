(pack-with-workshop-sdk)=

# Pack an inference snap with Workshop

In this tutorial, you'll package an inference snap with Workshop and the `inference-snaps-sdk`.
Starting from a repository created from the inference snap template, you'll prepare two input files, run the packaging pipeline through an LLM agent, then install and start the snap on your machine.
You'll do this with the *Qwen 3.5* model, and by the end you'll have a working `qwen3-5` snap that starts and responds to a prompt.

```{admonition} Experimental workflow
This packaging workflow is experimental. The steps, prompts, and generated outputs may change, and you may need to make manual adjustments before opening a pull request.
```

## Before you start

This tutorial builds on snap packaging and inference snap concepts. Work through these tutorials first:

- [Create an inference snap](../tutorial/create-inference-snap.md)
- [Get started with Qwen VL](../tutorial/qwen-vl-tutorial.md)

To follow this tutorial, you need:

- Workshop installed on your machine. If you are new to Workshop, see [Install Workshop](https://ubuntu.com/workshop/docs/tutorial/part-1-get-started/#).
- A Linux machine where you can install and run snaps.
- A GitHub repository for your inference snap, created from the [inference snap template](https://github.com/canonical/inference-snap-template). The template already contains the `workshop.yaml`, `Makefile`, and `README.md` files used here.
- Network access to download the *Qwen 3.5* model files from Hugging Face.

For more detail on the tools used along the way, see [OpenCode integration](../how-to/integration/opencode.md), the [Network ports registry](../reference/network-ports.md), and [Troubleshooting](../how-to/troubleshooting/index.md).

## 1. Open your snap project

Clone your inference snap repository and move into it:

```shell
git clone <your-inference-snap-repo-url>
cd <your-inference-snap-repo>
```

List the directory. You should see three files:

```{terminal}
:input: ls
Makefile  README.md  workshop.yaml
```

These files are the inputs for the packaging pipeline.

## 2. Prepare the inputs

The pipeline needs two things from you: the model files to download, and the snap metadata.
You provide these by editing the `Makefile` and the `README.md`.

### Makefile

First, open the `Makefile` and replace its contents with the following, which downloads a single *Qwen 3.5* model file:

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

This sets up the Hugging Face CLI in a local virtual environment and downloads the model file into the `components/` directory, where the pipeline expects the files that become part of the snap.

### README.md

Next, open `README.md` and fill in the metadata block at the top of the file, between the `<!--` and `-->` comments, with these values:

```text
snap-name: qwen3-5
snap-friendly-name: Qwen 3.5
model-card: https://qwen.ai/blog?id=qwen3.5
http-port: 8352
webui-http-port: 8353
engines: cpu
```

Notice a few things about this block:

- `snap-name` uses only lowercase letters, digits, and hyphens. It becomes the CLI command users run after installation.
- The ports `8352` and `8353` must not clash with entries in the [Network ports registry](../reference/network-ports.md).
- `engines: cpu` keeps this first build simple by targeting a single runtime.

Your two inputs are now ready: the `Makefile` downloads the model, and the `README.md` metadata describes the snap.

## 3. Package the inference snap

With the inputs ready, launch Workshop and run the packaging pipeline.

### Launch Workshop

From your project directory, start the Workshop environment described by `workshop.yaml`:

```shell
workshop launch
```

Then open a shell inside it:

```shell
workshop shell
```

Confirm that you are inside the Workshop shell:

```{terminal}
workshop@dev:/project$ whoami
workshop
```

If the output is not `workshop`, stop here. Exit and check that Workshop launched correctly before you continue.

### Start OpenCode

The `inference-snaps-sdk` ships with OpenCode, a terminal UI for running LLM agents. For more information, see the [OpenCode documentation](https://opencode.ai/docs).

Inside the Workshop shell, start OpenCode:

```shell
opencode --auto
```

The OpenCode TUI opens with the skills and agents installed by `inference-snaps-sdk`. The `--auto` flag approves permission prompts automatically, which is safe here because the Workshop environment is sandboxed.

OpenCode can be configured to use a local LLM or a remote API. If you want to use a local model, see [OpenCode configuration](../how-to/integration/opencode.md#configure-opencode). You can also use a remote model by providing an API key for a service like Claude or OpenAI. 
In order to do that tou can send this command in the OpenCode TUI:

```text
/connect
```

You will be prompted to select a provider and enter your API key or redirected to a web page to log in. After you connect, you can use the remote model for the packaging pipeline.
By default, Opencode uses the `Big Pickle` model from Opencode Zen. You are free to use any model you prefer.

### Run the packaging pipeline

In the OpenCode TUI, enter this prompt:

```text
start snap packaging pipeline
```

The pipeline validates your inputs, summarizes what it found, and asks you to confirm before it starts. Read the summary and check that it matches:

- the snap name `qwen3-5`
- the *Qwen 3.5* model file from the `Makefile`
- the `cpu` engine

Confirm to continue. The pipeline generates the packaging files, builds the snap, and prepares a pull request description.

When the run finishes, you'll have:

- one or more build artifacts, such as `.snap` and `.comp` files, in the repository
- a generated PR description ready for review

If the pipeline stops on a validation error, fix the reported problem and ask the agent to continue.

## 4. Test your snap

Exit the Workshop shell to return to your machine:

```{terminal}
workshop@dev:/project$ exit
exit
user@host-machine:~$ whoami
user
```

Move into your project directory and install the generated artifacts:

```shell
cd <your-inference-snap-repo>
sudo snap install *.snap *.comp --dangerous
```

Now start the snap and chat with the model:

```shell
qwen3-5 chat
```

The command connects to the model server and waits for your prompt. Type a message and press ENTER, and the model responds.

Next, check the configured ports:

```{terminal}
:input: qwen3-5 get
http.host: 127.0.0.1
http.port: 8352
...
webui.http.host: 127.0.0.1
webui.http.port: 8353
```

Notice that the ports match the `8352` and `8353` values you set in `README.md`. Open the Web UI in your browser at `http://127.0.0.1:8353`.

You have now confirmed that the `qwen3-5` snap installs, starts, and exposes the ports you configured.

## 5. Open a PR for your newly packaged inference snap

The packaging run generated a PR description together with the build outputs.

Before you open a pull request, confirm that the packaging matches what you intended:

```{terminal}
:input: qwen3-5 list-engines
cpu
```

```{terminal}
:input: qwen3-5 list-models
model-q4-k-xl-gguf
```

The `cpu` engine and the *Qwen 3.5* model match the values you set in `README.md` and the file downloaded by the `Makefile`.

With the build artifacts produced and the snap running locally, open a pull request to your inference snap repository using the generated PR description as the starting point.

You've now packaged an inference snap end to end with Workshop: you prepared the inputs, ran the packaging pipeline, installed the snap, and confirmed it runs.
