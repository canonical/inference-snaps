# Your first inference snap

Inference snaps provide an efficient method to distribute large language models (LLMs) for inference on a variety of hardware.
Each inference snap can include multiple runtimes and model weights optimized for different hardware configurations.
During installation, the snap automatically detects the hardware and picks the best matching runtime and model weights.
Only the selected components are downloaded and installed. 

In this tutorial, you will create an inference snap for the Gemma 3 LLM. 
The tutorial uses Gemma 3 as an open source example model. 
Gemma models are free and permit responsible commercial use.
Read their full [terms of use](https://ai.google.dev/gemma/terms).

The focus of the tutorial is on learning the overall structure of inference snaps by packaging a single model and runtime, rather than making a snap that works on a wide range of machines.
Once you learn the basics, you'll be able to apply the newly gained skill for other model weights and runtimes. 

To succeed, you should be familiar with LLMs and Linux command line. 
You should also have access to a Linux machine with the amd64 architecture, and have at least 20GB of free disk space.

## Get started

Start by preparing the necessary tools.

The {command}`snapcraft` utility is used to create snaps.
If you don't have it installed, refer to this [setup guide](https://documentation.ubuntu.com/snapcraft/stable/how-to/set-up-snapcraft/).

You may also use {command}`tree` and {command}`curl` which you can install with:
```shell
sudo apt update
sudo apt install tree curl
```

## snapcraft.yaml

The {file}`snapcraft.yaml` file is where the snap's packaging logic is defined.

Create and enter a new directory for the project:
```shell
mkdir gemma3-snap
cd gemma3-snap
```

Use the following command to create a template in an initial project tree:
```shell
snapcraft init
```

This will create a `snap` directory with a {file}`snapcraft.yaml` file inside:
```{terminal}
:dir: gemma3-snap
:input: tree .
.
└── snap
    └── snapcraft.yaml

2 directories, 1 file
```

Update the snapcraft file with sensible metadata.

Remove the generated `parts` section.

Inference snaps are usually named after the model (e.g., `gemma3` for Gemma 3 model). 
Add your username to it as a suffix to create a unique snap name that can be uploaded to the store in order to complete this tutorial.
In this tutorial, `jane` is used as a placeholder for your username.

The snapcraft file may look like this:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:caption: snap/snapcraft.yaml
:language: yaml
:start-after: start metadata
:end-before: end metadata
```


## Craft the snap

You need to extend the snapcraft file by adding:
- `parts`: the build logic
- `apps`: commands and services 
- `components`: optional building blocks
- `environments`: global environment variables
- `hooks`: settings for scripts that trigger on certain events

You also need to add a few supporting scripts.

### Global environment variables

It is useful to set some variables globally. Do this only for variables that are read from different parts of the snap's runtime environment.

Add the following to the snapcraft file for the common environment variables:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:caption: snap/snapcraft.yaml
:language: yaml
:start-after: start environment
:end-before: end environment
```

### The CLI app

Inference snaps include a simple and powerful CLI. 
The CLI provides hardware detection, engine selection, component installation, server startup, config management, and other functionalities.

An open source implementation of the CLI is available [on GitHub](https://github.com/canonical/inference-snaps-cli). Add it to your snap.

Create a `parts` section in the snapcraft file:
```yaml
parts:
  ...
```

Add a *part* in the snapcraft file to build and include the CLI:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start cli part
:end-before: end part
```
This will add `modelctl` command line tool to the snap package.
You will use it inside the snap for various tasks, such as engine selection and configuration.

Then, create an `apps` entry in the snapcraft file:
```yaml
apps:
  ...
```

Add the following *app*:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start cli app
:end-before: end app
```
This app exposes the `modelctl` command line tool to the user, under the snap's name.

### Components

In inference snaps, the model weights and runtimes are packages into [snap components](https://snapcraft.io/docs/components). 
Using snap components enables building the snap with all the building blocks but installing only what is relevant.

Start with the bare minimum: one runtime and one model weights component.

Inference with a CPU is always a good starting point. One way to do it is with a {spellexception}`llama.cpp`-based runtime. 
The [{spellexception}`llama.cpp`](https://github.com/ggml-org/llama.cpp) project is backed by a large open-source community working on various silicon optimizations. 
Use the [{spellexception}`llama.cpp` HTTP server](https://github.com/ggml-org/llama.cpp/tree/master/tools/server) implementation which supports model weights in the GGUF format.

Google distributes GGUF formats for most of Gemma 3 model sizes. 
Go with a relatively small model to speed up development times and add larger ones later.
A good first model to use is the [`gemma-3-1b-it-qat-q4_0-gguf`](https://huggingface.co/google/gemma-3-1b-it-qat-q4_0-gguf) variant on Hugging Face.  
On that page, download the {file}`gemma-3-1b-it-q4_0.gguf` file (the file name differs slightly from the model name).

Create a `components` directory with two subdirectories:
- `model-1b-it-q4-0-gguf` for the model weights
- `llama-cpp` for the runtime

```{tip}
Use the [kebab-case](https://developer.mozilla.org/en-US/docs/Glossary/Kebab_case) naming convention for the directories.
This isn't required for the directory names but it makes the next steps easier as most parts of snaps support kebab-case naming only.
```

Move the downloaded GGUF file to the model directory.

The {spellexception}`llama.cpp`'s HTTP server binary will be built from source during packaging and placed in the expected location. More on that later in the same section.

Add the following configuration file the model:
```{literalinclude} create-inference-snap/components/model-1b-it-q4-0-gguf/component.yaml
:caption: {spellexception}`components/model-1b-it-q4-0-gguf/component.yaml`
:language: yaml
```

And the following for the runtime:
```{literalinclude} create-inference-snap/components/llama-cpp/component.yaml
:caption: {spellexception}`components/llama-cpp/component.yaml`
:language: yaml
```

The variables set above per component are passed on to the execution environment.

Finally, add a wrapper script for starting the server provided by this component:
```{literalinclude} create-inference-snap/components/llama-cpp/server
:caption: {spellexception}`components/llama-cpp/server`
:language: bash
```

Make the script executable:
```shell
chmod +x components/llama-cpp/server
```

The above script loads the configurations and starts the HTTP server.

At this point, the {file}`components` directory should look like this:
```{terminal}
:dir: gemma3-snap
:input: tree components/
components/
├── llama-cpp
│   ├── component.yaml
│   └── server
└── model-1b-it-q4-0-gguf
    ├── component.yaml
    └── gemma-3-1b-it-q4_0.gguf

3 directories, 4 files
```

You've created a couple of component directories, but the snap doesn't know anything about them yet.

Edit the snapcraft file.
Define one *component* for the runtime, and another one for the model weights:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start components
:end-before: end components
```

Here, the component names must be kebab-case. 

If you build the snap now, {command}`snapcraft` will create one `.snap` and two `.comp` archives. 
The components will be empty.
You need to specify what gets packaged in them.

Add a new *part* to copy the local files from the `components` directory to the corresponding component within the package.

```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start component-local-files part
:end-before: end part
```

`(component/<component name>)` is a special syntax to reference a defined component.

Furthermore, add a *part* to build {spellexception}`llama.cpp` from source and move the artifacts to the corresponding component:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start llama-cpp part
:end-before: end part
```

You don't need all the {spellexception}`llama.cpp` artifacts, but that's okay for now. 
You can refine this *part* later to only include what is necessary.

The packaging logic for the components is now complete.
You can build the snap and verify the packages:
```{terminal}
:dir: gemma3-snap
:input: snapcraft pack
Packed: gemma3-jane_v3_amd64.snap, gemma3-jane+model-1b-it-q4-0-gguf.comp, gemma3-jane+llama-cpp.comp
```

### Engine

The runtime and model weights are tightly coupled components useful for certain host machines.
Each inference snap can include a range of such components. 
In this tutorial you will create a single {ref}`engine <engines>` that uses those components. 
Engines can reuse and customize components, for example to bundle a runtime with different model weights or to tweak the command line arguments.

Create an `engines` directory, with a subdirectory for our only engine named `generic-cpu`. 

Inside, create an `engine.yaml` file that defines the engine:
```{literalinclude} create-inference-snap/engines/generic-cpu/engine.yaml
:caption: engines/generic-cpu/engine.yaml
:language: yaml
```

That is the {ref}`engine manifest file <engine-manifest>`. 
It describes the engine as well as its hardware and software requirements.
The manifest file can be extended to restrict the hardware requirements and also set default configurations for the runtime.

Lastly, create a wrapper script that starts the runtime for this engine. 
In simplest cases, this simply runs a program provided by one of the components:
```{literalinclude} create-inference-snap/engines/generic-cpu/server
:caption: engines/generic-cpu/server
:language: bash
```

Make this script executable:
```shell
chmod +x engines/generic-cpu/server
```

Add another *part* in snapcraft file to include the `engines` directory in the snap:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start engines part
:end-before: end part
```

### The server app

Now, you need to expose the engine's server to the snap. 
The snap should run the server (provided by the selected engine) as a background service.

In a top level `scripts` directory, add a wrapper script that abstracts away the engine selection and runs the server script of the selected engine:
```{literalinclude} create-inference-snap/scripts/server.sh
:caption: scripts/server.sh
:language: bash
``` 

Add executable permission to the script:
```shell
chmod +x scripts/server.sh
```

Include the directory in the package, via a new *part*:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start scripts part
:end-before: end part
```

Then, add a server *app* that runs `server.sh` as a background service:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start server app
:end-before: end app
```
The above goes to the `apps` section of the snapcraft file.

### The install hook

Moreover, you need a script which takes care of the installation logic.
In snaps, the `install` hook is a script that runs when the snap is installed.

The install hook should trigger hardware detection and install the best matching components and engine.
Do this using the `modelctl` CLI tool included earlier:
```{literalinclude} create-inference-snap/snap/hooks/install
:caption: snap/hooks/install
:language: bash
```

This script gets added automatically to the snap if placed in the `snap/hooks` directory.
However, you need to grant it the necessary permissions to detect hardware.

Add the following to the snapcraft file to declare the required interface for the `install` hook:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start hooks
:end-before: end hooks
```

### Licensing

Finally, add all license files to the snap. E.g. to add Gemma's license, add the following *part*:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start gemma license part
:end-before: end part
```

## Pack and deploy

With the crafting completed, it's time to pack, test, and upload the snap.

### Test locally

Pack the snap:
```shell
snapcraft pack
```

This will create the following files:
- `gemma3-jane_v3_amd64.snap` - the snap package
- `gemma3-jane+model-1b-it-q4-0-gguf.comp` - component package for the model weights
- `gemma3-jane+llama-cpp.comp` - component package for the runtime

Install the snap and the components:
```shell
sudo snap install --dangerous *.snap && \
sudo snap install --dangerous *.comp
```

The server will fail to start, so you can check the logs with: `snap logs gemma3-jane`.
The reason for this is that the snap does not have access to detect the hardware and select the engine.

Grant the snap access to the `hardware-observe` interface:
```shell
sudo snap connect gemma3-jane:hardware-observe
```

Now, request auto selection of the engine:
```{terminal}
:dir: gemma3-snap
:input: sudo gemma3-jane use-engine --auto
Evaluating engines for optimal hardware compatibility:
✔ generic-cpu: compatible, score=12
Selected engine: generic-cpu
Engine successfully changed to "generic-cpu"
```
As the output shows, the snap has detected that `generic-cpu` engine is compatible with the host machine. 
If there were other engines, it would have picked the best matching one.

Start the server:
```shell
sudo snap start gemma3-jane
```

Next, check the logs again and make sure the server started successfully:
```{terminal}
:dir: gemma3-snap
:input: sudo snap logs gemma3-jane
...
2025-11-13T15:34:52+01:00 gemma3-jane.server[241059]: main: server is listening on http://127.0.0.1:9090 - starting the main loop
2025-11-13T15:34:52+01:00 gemma3-jane.server[241059]: srv  update_slots: all slots are idle
```

Looks like the server is up and running!

Submit a request using {command}`curl`:
```{terminal}
:dir: gemma3-snap
:input: curl -s http://localhost:9090/v1/chat/completions   -H "Content-Type: application/json"   -d '{
    "max_tokens": 30,
    "temperature": 0,
    "stream": false,
    "messages": [
      { "role": "system", "content": "You are a helpful assistant." },
      { "role": "user", "content": "What are the 3 main tourist attractions in Paris?" }
    ]                                                                                   
  }'
{"choices":[{"finish_reason":"length","index":0,"message":{"role":"assistant","content":"Okay, here are the 3 main tourist attractions in Paris:\n\n1.  **The Eiffel Tower:** Undoubtedly the most iconic landmark in the world"}}],"created":1763044754,"model":"gpt-3.5-turbo","system_fingerprint":"b1-6a746cf","object":"chat.completion","usage":{"completion_tokens":30,"prompt_tokens":28,"total_tokens":58},"id":"chatcmpl-PLT9wXbZiJ6XXxHCQUBsTz9wj59y2KHU","timings":{"prompt_n":28,"prompt_ms":90.995,"prompt_per_token_ms":3.249821428571429,"prompt_per_second":307.7092147920215,"predicted_n":30,"predicted_ms":833.198,"predicted_per_token_ms":27.773266666666665,"predicted_per_second":36.005847349609574}}
```

That worked! You've created an inference snap for Gemma 3 model, with just one engine.


### Publish to the store

Uploading the snap to the store unlocks the automatic engine selection and component installation functionality.
Moreover, it lets others install and use the snap.

Uploading requires a few steps. Here is the list:
1. Register the name
2. Get yourself added to component upload allow-list
3. Upload the snap
4. Request auto connection permission for the hardware-observe interface

First you need to register the snap name in the store.
Refer to [this guide](https://snapcraft.io/docs/registering-your-app-name) to do that.

Uploading snaps with components currently requires a special permission. This is needed per account.
Request to be added to the component upload allow-list by posting in the [Snapcraft forum](https://forum.snapcraft.io/c/store-requests/19).
There is no subcategory for component uploads, so just post without picking one.

Once you have the permission, upload the snap and its components:
```shell
snapcraft upload gemma3-jane_v3_amd64.snap \
  --component model-1b-it-q4-0-gguf=gemma3-jane+model-1b-it-q4-0-gguf.comp \
  --component llama-cpp=gemma3-jane+llama-cpp.comp \
  --release edge
```

The snap is now in the store!

Remove the locally installed snap and install from the store and test the snap in developer mode:
```{terminal}
:dir: gemma3-snap
:input: sudo snap remove gemma3-jane
gemma3-jane removed

:input: sudo snap install --devmode gemma3-jane --edge
gemma3-jane (edge) v3 from Jane Doe installed

:input: gemma3-jane status
engine: generic-cpu
endpoints:
    openai: http://localhost:9090/v1
```

As shown, the snap has automatically detected generic-cpu this time. 

The [developer mode](https://snapcraft.io/docs/install-modes#developer-mode) allows the snap to access the system without the usual confinement constraints.
This isn't ideal for production use.

If you install the snap in normal (confined) mode at this point, it will install but skip hardware detection and engine selection.
You can connect the [hardware-observe interface](https://snapcraft.io/docs/hardware-observe-interface) and then run the auto engine selection manually:
```shell
sudo snap install gemma3-jane --edge
sudo snap connect gemma3-jane:hardware-observe
sudo gemma3-jane use-engine --auto
```

But you can do better than that.

The snap store can grant the snap permission to auto connect the hardware-observe interface. This requires a manual review by the store security team.

Submit a request on the [Snapcraft forum](https://forum.snapcraft.io/c/store-requests/privileged-interfaces/27) to get this permission for your snap.
An example, extended request can be found [here](https://forum.snapcraft.io/t/auto-connection-of-hardware-observe-and-home-for-gemma3/49246).

Once the permission is granted, you can install the snap in confined mode, without setting `--devmode` and it will be ready to use right away.

## Next steps

You've created a basic inference snap for Gemma 3 model with a single engine. 
You can extend it in various ways, for example to add an additional engine for GPU support.
This could be another build of {spellexception}`llama.cpp` targeting CUDA, or a different runtime altogether, such as [OpenVINO Model Server](https://github.com/openvinotoolkit/model_server) which supports Intel GPUs and NPUs.

You may also have noticed the CLI doesn't offer tab completion.
You can find some guidance for adding it temporarily (e.g., for Bash by running {command}`gemma3-jane completion bash -h`).
Ideally, the tab completion should be added to the snap itself so it becomes available during the installation.

Look into {ref}`existing inference snaps <available-snaps>` to learn more about these and other advanced topics.
