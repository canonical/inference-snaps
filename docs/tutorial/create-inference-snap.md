# Creating inference snaps

In this tutorial we will create an inference snap for a large language model (LLM).
The focus is on learning the overall structure by packaging a single model and runtime, rather than making a snap that works on a wide range of machines.
Once we build the basics, it is easy to extend it to include other model weights and runtimes. 

## What you'll need

We use `snapcraft` to create snaps.
If you don't have it installed, refer to [this guide](https://documentation.ubuntu.com/snapcraft/stable/how-to/set-up-snapcraft/)

You will also need the model weights and a compatible runtime. 
But first, we need to pick an LLM. Gemma 3 is a good choice for this tutorial because it is open and available in small sizes.
Gemma 3 is free and permits responsible commercial use.
Read its terms of use [here](https://ai.google.dev/gemma/terms).

Inferencing with a CPU is always a good starting point and a simple way to do that is via {spellexception}`llama.cpp` runtime. 
The {spellexception}`llama.cpp` project is backed by a large open source community working on various silicon optimizations. 
We can use the llama-server implementation which supports model weights in the GGUF format.
Google distributes GGUF formats for most of Gemma 3 model sizes. 
Go with a relatively small model to speed up development times and add larger ones later.

The [`gemma-3-1b-it-qat-q4_0-gguf`](https://huggingface.co/google/gemma-3-1b-it-qat-q4_0-gguf) is a good first instance.
Download the `.gguf` file.

You may also use `tree` and `curl` which you can installed with:
```shell
sudo apt update
sudo apt install tree curl
```

## The snapcraft.yaml

Start by creating the `snapcraft.yaml` file. Use the following command to create an initial project tree:
```shell
snapcraft init
```

Checkout the generated file:
```{terminal}
:input: tree .
.
└── snap
    └── snapcraft.yaml

2 directories, 1 file
```

Update the snapcraft file with sensible metadata. Remove the parts.

Inference snaps are usually named after the model, e.g. `gemma3` for Gemma 3 model. 
Add your username to it as a suffix to create a unique snap name that can be uploaded to the store in order to complete this tutorial.

My snapcraft file looks like this:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:caption: snap/snapcraft.yaml
:language: yaml
:start-after: start metadata
:end-before: end metadata
```

We need to complete the snapcraft file by adding:
- parts: the build logic
- apps: commands and services 
- components: optional building blocks
- environments: global environment variables
- hooks: settings for snap hooks

Start by settings the global variables, toward preparing and building profile files, to the definition of apps.

## Global environment variables

Add the following snippet.
Setting the variables at this level make them globally available to the snap's runtime environment.

```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:caption: snap/snapcraft.yaml
:language: yaml
:start-after: start environment
:end-before: end environment
```

## Project files

We need to organize project files in a clear directory structure. 

Create a directory for the runtime, and another one for the model weights. 
Use the kebab-case naming convention. This isn't required but it makes the next steps easier as most parts of snaps support kebab-case naming only.

Runtime and model weights are considered optional building blocks of this snap. 
This does not make perfect sense right now because we only have one of each. 
But imagine this snap growing to have multiple of each, selected and downloaded on runtime based on the host machine.

We organize optional building blocks in [snap components](https://snapcraft.io/docs/components). 
This enables us to build the snap with all the building blocks but install only what is needed.

Create a `components` directory with two subdirectories: `llamacpp` and `model-1b-it-q4-0-gguf`

Move the downloaded `.gguf` file to the model directory. 

Add the following configuration file for the model:
```{literalinclude} create-inference-snap/components/model-1b-it-q4-0-gguf/component.yaml
:caption: {spellexception}`components/model-1b-it-q4-0-gguf/component.yaml`
:language: yaml
```

This is used to find the `.gguf` file within the component's directory.

Now, add `llamacpp`'s configuration file:
```{literalinclude} create-inference-snap/components/llamacpp/component.yaml
:caption: {spellexception}`components/llamacpp/component.yaml`
:language: yaml
```

These variables are passed on to the application that runs {spellexception}`llama.cpp`'s server. 

Finally, add a wrapper script for llama-server:
```{literalinclude} create-inference-snap/components/llamacpp/server
:caption: {spellexception}`components/llamacpp/server`
:language: bash
```
Make the script executable: `chmod +x engines/generic-cpu/server`

The `components` directory should now look like:
```{terminal}
:input: tree components/
components/
├── llamacpp
│   ├── component.yaml
│   └── server
└── model-1b-it-q4-0-gguf
    └── component.yaml

3 directories, 3 files
```

## Components

We have created a couple of component directories, but the snap doesn't know anything about them. 

Edit the `snapcraft.yaml` file.
Define one component for the runtime, and another one for the model weights:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start components
:end-before: end components
```

Here, the component names must be kebab-case. 

If you build the snap now, it will create the `.snap` and two `.comp` files. 
The components will be empty.
We need to specify what gets packaged in them.

Add the following to `parts` to copy the local files in the `components` directory to the corresponding component within the package.

```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start component-local-files part
:end-before: end part
```
<!-- allowed: ['append', 'caption', 'class', 'dedent', 'diff', 'emphasize-lines', 'encoding', 'end-at', 'end-before', 'force', 'language', 'lineno-match', 'lineno-start', 'linenos', 'lines', 'name', 'prepend', 'pyobject', 'start-after', 'start-at', 'tab-width']) [myst.directive_option -->

`(component/<component name>)` is a special syntax to reference a defined component.

Furthermore, add a part to build {spellexception}`llama.cpp` from source and move the artifacts to the corresponding component:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start llamacpp part
:end-before: end part
```

The `parts` section of the snapcraft file should now include two parts:
```yaml
parts:
  component-local-files:
    ...
  llamacpp:
    ...
```

The packaging logic for the components are now complete.
You can build the snap and verify the packages:
```{terminal}
:input: snapcraft pack
Packed: gemma3-jane_v3_amd64.snap, gemma3-jane+model-1b-it-q4-0-gguf.comp, gemma3-jane+llamacpp.comp
```

Next, we need to extend the snap to include

## Engine

The runtime and model weights are tightly coupled components useful for certain host machines.
Each inference snaps can include a range of such components. 
In this tutorial we create a single {ref}`engine <engines>` that uses those components. 
Engines can reuse and customize components, for example to tweak default runtime configurations.

Create an `engines` directory, with a subdirectory for our only engine named `generic-cpu`. 

Next, create an {ref}`engine manifest file <engine-manifest>`:
```{literalinclude} create-inference-snap/engines/generic-cpu/engine.yaml
:caption: engines/generic-cpu/engine.yaml
:language: yaml
```

Lastly, create a wrapper script that starts the runtime for this engine. 
In simplest cases, this simply runs the runtime component's server script:
```{literalinclude} create-inference-snap/engines/generic-cpu/server
:caption: engines/generic-cpu/server
:language: bash
```
Make this script executable: `chmod +x engines/generic-cpu/server`


Add a part in snapcraft file to include the engines directory in the snap:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start engines part
:end-before: end part
```

## Generic packaging bits

Inference snaps include some common parts for standard functionality such as CLI, server execution, and installation logic.

The most important piece is the CLI. This takes care of several things, such as hardware detection, engine selection, component installation, server startup, config management, and more.

Add a part in the snapcraft file to build and include the CLI:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:start-after: start cli part
:end-before: end part
```

The snap needs a top level script that defines how to run the engine's server script:
```{literalinclude} create-inference-snap/scripts/server.sh
:caption: scripts/server.sh
:language: bash
``` 

Include this script in the snapcraft file:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:start-after: start scripts part
:end-before: end part
```

Moreover, we need a script which takes care of the installation logic.
This script is responsible for detecting the hardware capabilities of the host machine, and installing the best matching components and engine.
```{literalinclude} create-inference-snap/snap/hooks/install
:caption: snap/hooks/install
:language: bash
```

This script gets added automatically to the snap if placed in the `snap/hooks` directory.
However, we need to grant it the necessary permissions to detect hardware.
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start hooks
:end-before: end hooks
```

Finally, add all license files to the snap. E.g. to add Gemma's license:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start gemma license part
:end-before: end part
```

## Apps to expose programs

The package is now complete, but we need to define how to run the snap.

Create an `apps` section in the snapcraft file:
```yaml
apps:
  # Name this after the snap to make it the default command
  gemma3-jane:
    ...

  # Server run as a service
  server:
    ...
```

The CLI app exposes the `modelctl` command line tool to the user, under the snap's name.:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start cli app
:end-before: end app
```

The server app runs the model as a background service:
```{literalinclude} create-inference-snap/snap/snapcraft.yaml
:language: yaml
:start-after: start server app
:end-before: end app
```

## Build and test locally

```{terminal}
:input: snapcraft pack
Packed: gemma3-jane_v3_amd64.snap, gemma3-jane+model-1b-it-q4-0-gguf.comp, gemma3-jane+llamacpp.comp

:input: sudo snap install --devmode gemma3-jane_v3_amd64.snap
gemma3-jane v3 installed
``` 

Checking `snap logs gemma3-jane` should show that the server has failed to start because no engine and components have been installed yet.
The reason for this is that the locally built snap cannot access the component files and install them automatically.

We have to install the components and select the engine manually now, but the install hook will take care of this automatically when the same snap is installed from the store.

```{terminal}
:input: sudo snap install --dangerous gemma3-jane+model-1b-it-q4-0-gguf.comp
component model-1b-it-q4-0-gguf v3 for gemma3-jane v3 installed

:input: sudo snap install --dangerous gemma3-jane+llamacpp.comp 
component llamacpp v3 for gemma3-jane v3 installed

:input: sudo gemma3-jane use-engine generic-cpu
Engine successfully changed to "generic-cpu"

:input: sudo snap start gemma3-jane
Started.
```

````{tip}
If want to repeat making changes and building the snap and its components without losing your sanity, chain the comments:
```shell
sudo snap install --dangerous gemma3-jane_v3_amd64.snap && \ 
sudo snap install --dangerous gemma3-jane+model-1b-it-q4-0-gguf.comp && \
sudo snap install --dangerous gemma3-jane+llamacpp.comp 
```

You can also use wildcards if you don't have any other snap artifacts in that directory:
```shell
sudo snap install --dangerous *.snap && sudo snap install --dangerous *.comp
```
````

After starting the server, check the logs again and make sure it started successfully:
```{terminal}
:input: snap logs gemma3-jane
...
2025-11-13T15:34:52+01:00 gemma3-jane.server[241059]: main: server is listening on http://127.0.0.1:9090 - starting the main loop
2025-11-13T15:34:52+01:00 gemma3-jane.server[241059]: srv  update_slots: all slots are idle
```

Looks like the server is up and running!

Submit a request via curl:
```{terminal}
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

That worked! We have created an inference snap for Gemma 3 model, with just one engine.


## Upload

Uploading the snap to the store unlocks the automatic engine selection and component installation functionality.
Moreover, it lets others easily install and use the snap.

Uploading requires a few steps. Here is the list:
1. Register the name
2. Get yourself added to component upload allow-list
3. Upload the snap
4. Request interface connection permission for the hardware-observe interface

First you need to register the snap name in the store.
Refer to [this guide](https://snapcraft.io/docs/registering-your-app-name) to do that.

Uploading snaps with components currently requires a special permission. This is needed per account.
Request to be added to the component upload allow-list by posting in the [Snapcraft forum](https://forum.snapcraft.io/c/store-requests/19).
There is no subcategory for component uploads, so just post without picking one.

Once you have the permission, upload the snap and its components:
```shell
snapcraft upload gemma3-jane_v3_amd64.snap \
  --component model-1b-it-q4-0-gguf=gemma3-jane+model-1b-it-q4-0-gguf.comp \
  --component llamacpp=gemma3-jane+llamacpp.comp \
  --release edge
```

At this point, the snap is in the store but it cannot access hardware information during the installation.
We can work around this by installing it in developer mode:
```shell
sudo snap install --devmode gemma3-jane --edge
```

The developer mode allows the snap to access many other system resources as well, which is not ideal.
Moreover, snaps installed in developer mode cannot be updated automatically from the store.

If we install the snap in normal (confined) mode at this point, it will install but skip hardware detection and engine selection.
We can still grant permission and resume manually:
```shell
sudo snap install gemma3-jane --edge
sudo snap connect gemma3-jane:hardware-observe
sudo gemma3-jane use-engine --auto
```

But we can do better than this.
The snap store can grant the snap permission to auto connect the hardware-observe plug.
Submit a request on the [Snapcraft forum](https://forum.snapcraft.io/c/store-requests/privileged-interfaces/27) to get this permission for your snap.
An example, extended request can be found [here](https://forum.snapcraft.io/t/auto-connection-of-hardware-observe-and-home-for-gemma3/49246).

Once the permission is granted, you can install the snap and it will be ready to use right away.

## Install and use

Install the snap from the store:
```{terminal}
:input: sudo snap install gemma3-jane --edge
gemma3-jane v3 installed

:input: gemma3-jane status
engine: generic-cpu
endpoints:
    openai: http://localhost:9090/v1
```

## Next steps

We created a basic inference snap for Gemma 3 model with a single engine. 
You can extend it in various ways, for example:
- Add an additional engine for GPU support
- Add support for arm64 architecture

You may also have noticed the the CLI doesn't offer tab completion.
You can find some guidance for adding it temporarily, e.g. for Bash `gemma3-jane completion bash -h`.
The tab completion can also be added to the snap itself.

Look into {ref}`existing inference snaps <available-snaps>` to learn more about these and other advanced topics.
