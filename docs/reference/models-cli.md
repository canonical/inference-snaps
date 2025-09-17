# Models CLI

This document provides a complete reference for the `ai model snaps` command line interface (CLI).
It describes all available commands, options, arguments, and expected behavior. 
The CLI interface is the same across all models.

## Overview

The CLI is the main entry point to interact with the model runtime, manage engines, configure settings, and run interactive sessions. 
The syntax used is:

```shell
<my-model> [COMMAND] [OPTIONS] [ARGS]

## Global Options

To show usage information and list all available commands:
```shell
<my-model> --help
```

For example:
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

## Basic commands

These provide the core functionality for everyday use such as checking tool status or starting an interactive chat session.

### `status`

Displays information about the server and runtime environment.

```shell
<my-model> status
```

Details shown:

* Server status and active endpoints
* Currently selected engine
* Availability of better engines, e.g., after attaching new hardware or installing drivers.
* Missing drivers for available hardware
* Resource usage by engines, including cached data

### `chat`

Launch the command-line chat application.
```shell
<my-model> chat
```


## Configuration commands

Change configuration options including reading, setting, and unsetting persistent parameters.

### `set`

Set configuration values.
```shell
<my-model> set <key>=<value>
```

The `Set` command overrides configuration values. 
Unknown keys are rejected.

### `get`

The `get` command retrieves configuration values.

```shell
<my-model> get [key]
 ```

If no key is given, all configurations are listed in tabular form.

### `unset`

Unset configuration values.
```shell
<my-model> unset <key>
```

## Management Commands

Manage available engines by listing, inspecting, and selecting them for use.

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

