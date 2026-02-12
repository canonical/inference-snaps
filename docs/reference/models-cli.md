(models-cli)=
# Command line interface (CLI)

This document provides a complete reference for the command line interface (CLI).
It describes all available commands, options, arguments, and expected behavior. 
The CLI interface is the same across all inference snaps.

## Overview

The CLI is the main entry point to interact with the inference snap, manage engines, configure settings, and run interactive sessions. 
The syntax used is:

```shell
<inference-snap> [COMMAND] [OPTIONS] [ARGS]
```

## Global options

To show usage information about any of the commands, set the `--help` flag.

For example:
```{terminal}
:input: gemma3 --help

gemma3 runs an engine that is optimized for your host machine,
providing a local service endpoint.

Use this command to configure the active engine, or switch to an alternative engine.

Usage:
  gemma3 [command]

Basic Commands:
  status       Show the status
  chat         Start the chat CLI

Configuration Commands:
  get          Print configurations
  set          Set configurations

Management Commands:
  list-engines List available engines
  show-engine  Print information about an engine
  use-engine   Select an engine

Additional Commands:
  show-machine Print information about the host machine
  prune-cache  Remove cached data
  help         Help about any command

Flags:
  -h, --help      help for cli
  -v, --verbose   Enable verbose logging

Use "gemma3 [command] --help" for more information about a command.
```

## Basic commands

These provide the core functionality for everyday use such as checking tool status or starting an interactive chat session.

### `status`

Displays information about the server and runtime environment.

```shell
<inference-snap> status
```

Details shown:

* Service status and active endpoints
* Currently selected engine
<!-- * Availability of better engines, e.g., after attaching new hardware or installing drivers.
* Missing drivers for available hardware
* Resource usage by engines, including cached data -->

### `chat`

Launch the command-line chat application.
```shell
<inference-snap> chat
```


## Configuration commands

Change configuration options including reading, and setting persistent parameters.

### `set`

Set configuration values.
```shell
<inference-snap> set <key>=<value>
```

The `set` command overrides configuration values. 
Unknown keys are rejected.

### `get`

The `get` command retrieves configuration values.

```shell
<inference-snap> get [key]
 ```

If no key is given, all configurations are listed in tabular form.


## Management commands

Manage available engines by listing, inspecting, and selecting them for use.

### `list-engines`

List engines with details including name, description, vendor, and compatibility.
```shell
<inference-snap> list-engines
```

The currently active engine is signified with an asterisk (*).

The default output format is `table`. 
For JSON serialization, set `--format=json`. 

### `show-engine`

Show engine details.
```shell
<inference-snap> show-engine <engine> [flags]
```

The default output format is `yaml`. 
For JSON serialization, set `--format=json`. 

The command's output includes information such as:

* Supported devices
* Default configurations
* Host compatibility

### `use-engine`

To switch to a specific engine:
```shell
<inference-snap> use-engine <engine>
```

To automatically select the best match, considering the hardware and compute capabilities of the host machine:
```shell
<inference-snap> use-engine --auto
```

Switching engines will trigger the download of engine resources (as snap components).
A confirmation prompt will appear before downloading.


```{tip}
The `use-engine --auto` command is usually called during the inference snap installation.
This command can also be used to revert to the original engine after manual changes,
or to select the best engine match after driver updates, or new hardware installation.
```

## Additional commands

### `show-machine`

Show information about the host machine.

```shell
<inference-snap> show-machine [flags]
```

The default output format is `yaml`. 
For JSON serialization, set `--format=json`.

### `prune-cache`

Remove cached resources of the snap or an engine.

To remove all cached data of the snap, including data of all inactive engines:
```shell
<inference-snap> prune-cache 
```

To remove cached data of a specific engine:
```shell
<inference-snap> prune-cache --engine=<engine>
```

The cache pruning typically involves removal of unused components, configurations, and temporary files.
Keep in mind that removal of components does not immediately free up disk space.
The components are marked for deletion and will be removed by the snap system at a later time.
