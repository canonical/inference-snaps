(gemma-4-tutorial)=

# Get started with Inference Snaps using Gemma 4

This tutorial walks you through the process of installing and using an inference snap.
Specifically, you will install and run *Gemma 4*, an open-weight and natively multimodal large language model (LLM). This model is also referred to as a Vision Language Model (VLM) because it can ingest and interpret images, along with text prompts.

## Set up your computer

To complete this tutorial you will need:

- A computer running an operating system that supports snaps, e.g., Ubuntu 24.04 LTS. Older Ubuntu versions aren't supported.
- Hardware drivers. These may be pre-installed in your OS, but some devices require {ref}`additional drivers <install-drivers>`.

## Install the snap

Install the Gemma 4 snap (which downloads the model with an effective 4 billion parameters by default). Please note that the snap version you see might be newer than the one shown here:

```{terminal}
sudo snap install gemma4

gemma4 2.0.0-beta.14+0237ac7 from Ubuntu✓ installed
```

During installation, the snap detects information about your system such as installed hardware and available memory.
It uses this information to select an appropriate engine for your setup. 

In this context, {ref}`engine <engines>` refers to a combination of model weights and a runtime, along with some utilities to make it work, like a server exposing a standard API.

## Check the status

Once the installation is complete, the server used for inference with the model starts automatically.
This server comes from the engine.

Run {command}`gemma4 status` to see the selected engine and the status of the server.
The output should include the server's status and the URL where the API can be accessed:

```{terminal}
gemma4 status 

engine: nvidia-gpu
services:
    server: active
    server-webui: active
entrypoints:
    openai:
        url: http://127.0.0.1:8336/v1
    webui:
        url: http://127.0.0.1:8337
model:
    name: gemma4-e4b

```

In this case, the `nvidia-gpu` engine was automatically selected during installation because it matches the best available hardware. You may end up with a different engine. You can check the list of the available engines and their compatibility with your hardware (last column):

```{terminal}
gemma4 engines

ENGINE        VENDOR         SUMMARY                                     COMPAT
nvidia-gpu*   Canonical Ltd  CUDA optimized for workstation GPUs         yes   
intel-gpu     Canonical Ltd  Optimized for Intel GPUs                    yes   
intel-cpu     Canonical Ltd  Optimized for Intel CPUs                    yes   
cpu           Canonical Ltd  CPU-optimized for workstations              yes   
intel-onemkl  Canonical Ltd  Intel CPU-optimized for workstations usin…  exptl 
amd-gpu       Canonical Ltd  ROCm-optimized for workstation GPUs         no   
```

The list shows the currently selected engine with an asterisk. You can also switch to another engine if you want to, which may need to download additional components:

```{terminal}
sudo gemma4 use-engine cpu  

Engine changed to "cpu".
```

Switch back to the engine selected during installation (the one shown by the first {command}`gemma4 status`), so that the rest of the tutorial matches. For example, if it was `nvidia-gpu`, run {command}`sudo gemma4 use-engine nvidia-gpu`.

The snap and command-line interface (CLI) versions can be checked to see which versions you are running:

```{terminal}
gemma4 version

snap: 2.0.0-beta.14+0237ac7
cli: v2.0.0-beta.14
```

## Chat with the LLM

The `gemma4` snap includes a basic chat client that can be used to chat with the LLM.
Use the `chat` command to start a chat session:

```{terminal}
gemma4 chat

Using server at http://127.0.0.1:8336/v1
Type your prompt, then ENTER to submit. CTRL-C to quit.
» Hi there
Hello! How can I help you today? 😊
»  
```

This chat client only supports text input.
For queries that include images, you can use the Web UI.

## Chat with the LLM using the Web UI

The provided Web UI is a more advanced client interface for working with LLMs.
You can access the web interface at the URL shown for the webui in {command}`gemma4 status` command (`http://127.0.0.1:8337/` in the example above):

```{image} web-ui-gemma4.png
:alt: Web interface for chatting with Gemma 4 model
:align: center
```

Moreover, the web interface shows options for controlling reasoning, checking status, and attaching files to prompts for additional context.

To try the model with an image, download this photo of the [{spellexception}`Dom Luís I Bridge`](https://commons.wikimedia.org/wiki/File:Portugal,_Porto,_Dominic_Lu%C3%ADs_l_Bridge_(52594422458).jpg) and attach it to a new chat. Then ask a question that needs some understanding of the image, such as "What is the significance of the triangles?"

On a laptop with an NVIDIA RTX PRO 1000 Blackwell, the response was almost immediate. The model explained that "the triangles are the scaffolding of stability", giving the bridge rigidity, distributing forces, and allowing long spans. This is correct (the screenshot below shows only part of the reply):

```{image} web-ui-gemma4-image-response.png
:alt: Response of the model
:align: center
```

## More advanced connection to the LLM via the API

The installed snap also exposes an OpenAI-compatible API. You can use it to integrate the model into your own applications or agentic workflows. For example, to submit a chat completion request using {command}`curl`:

```shell
curl --silent --show-error http://127.0.0.1:8336/v1/chat/completions   --data '{
    "max_tokens": 300,
    "messages": [
      { "role": "system", "content": "You are a helpful assistant." },
      { "role": "user", "content": "What are the 3 main tourist attractions in Paris?" }
    ]                                                              
  }' | jq

```

which returns the following output (trimmed to fit) plus some request metadata and performance metrics:

```json
{
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "While Paris has countless beautiful sights, if you are looking for the absolute most iconic and essential ...",
        "reasoning_content": "Here's a thinking process to answer the request: ..."
      }
    }
  ],
  // more metadata omitted
  "model": "gemma4-e4b",
  "usage": {
    "completion_tokens": 1018,
    // more usage fields omitted
  },
  "timings": {
    "cache_n": 0,
    "prompt_n": 33,
    "prompt_ms": 481.579,
    "prompt_per_token_ms": 14.59330303030303,
    // more timing fields omitted
  }
}

```

Note that the above command needs {command}`jq` to be installed (which is required only to render the JSON output properly).

## Next steps

To become more familiar with inference snaps:
- Learn more about how to {ref}`manage and configure inference snaps <usage>` after installation
- Configure an inference snap to {ref}`work with your (agentic) clients <integration>`
