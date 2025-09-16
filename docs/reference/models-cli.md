# Models CLI

The CLI interface is the same across all models.
This document provides a complete reference for the **`ai model snaps`** command line interface (CLI).
It describes all available commands, options, arguments, and expected behavior.

## Overview

CLI syntax.
```shell
<my-model> [COMMAND] [OPTIONS] [ARGS]
```

The main entry point to interact with the model runtime, manage engines, configure settings, and run interactive sessions.

## Global Options

Show usage information and list all available commands.
```shell
<my-model> --help
```

**Example**:
```{terminal}
:input: deepseek-r1 --help

Usage:
  deepseek-r1 [command]

Basic Commands:
  status       Show the status
  chat         Start the chat CLI

Configuration Commands:
  get          Print configuration option
  set          Set configuration option
  unset        Unset configuration option

Management Commands:
  list-engines List available engines
  show-engine  Print information about an engine
  use-engine   Select an engine

Additional Commands:
  debug        Debugging commands
  help         Help about any command

Flags:
  -h, --help   help for deepseek-r1

Use "deepseek-r1 [command] --help" for more information about a command.
```

## Commands

Following commands are available:

### `status`

Display information about the server and runtime environment.

```shell
<my-model> status
```

**Details shown:**

* Server status and active endpoints
* Currently selected engine
* Availability of better engines (e.g., after attaching new hardware or installing drivers)
* Missing drivers for available hardware
* Resource usage by engines, including cached data

### `chat`

Launch the command-line chat application.
```shell
<my-model> chat
```

### `set` and `get`

Set configuration values
```shell
<my-model> set <key>=<value>
```

Get configuration values
```shell
<my-model> get [key]
 ```

* **Set**: Override configuration values. Rejects unknown keys.
* **Get**: Retrieve configuration values. If no key is given, list all configurations in tabular form.

Key syntax: follows the same hierarchical format as [**snapd**][snap_options].

**Example**:
```{terminal}
:input: deepseek-r1 set http.port=34531

:input: deepseek-r1 get http.port
34531
:input: deepseek-r1 get
Key      Value
engine   cpu-avx512
engines  {...}
http     {...}
model    model-distill-qwen-7b-q4-k-m-gguf
server   llamacpp-avx512
verbose  false
```

### `list-engines`

List engines with details including name, description, vendor, and compatibility.
```shell
<my-model> list-engines
```

### `show-engine`

Show engine details.
```shell
<my-model> show-engine <engine> [flags]
```

**Includes:**

* Supported devices
* Default configurations

```
--format <f>
```

Output format.
- `yaml` (default)
- `json`

**Example**
```{terminal}
:input: deepseek-r1 show-engine cpu-avx512 --format=json

{
   "engines.cpu-avx512":{
      "compatible":true,
      "components":[
         "llamacpp-avx512",
         "model-distill-qwen-7b-q4-k-m-gguf"
      ],
      "configurations":{
         "http.base-path":"v1",
         "model":"model-distill-qwen-7b-q4-k-m-gguf",
         "server":"llamacpp-avx512"
      },
      "description":"For CPUs with AVX512 support",
      "devices":{
         "any":[
            {
               "architecture":"amd64",
               "flags":[
                  "sse4_2",
                  "f16c",
                  "fma",
                  "avx",
                  "avx2",
                  "avx512f",
                  "avx512dq",
                  "avx512bw"
               ],
               "type":"cpu"
            }
         ]
      },
      "disk-space":"8G",
      "grade":"stable",
      "memory":"4G",
      "name":"cpu-avx512",
      "score":20,
      "vendor":"Canonical Ltd"
   }
}
```


### `use-engine`

Switch to a specific engine.
```shell
<my-model> use-engine <engine>
```

Automatically select best engine.
```shell
<my-model> use-engine --auto
```

**Behavior:**

* Switching engines triggers download of stack resources (as snap components).
* A confirmation prompt appears before downloading.
* `--auto` is run by default during installation.
* Can be used to revert to the original engine after manual changes, driver updates, or new hardware installation.

