# Models Command Line Interface (CLI)

This document provides a complete reference for the `ai model snaps` command line interface (CLI).
It describes all available commands, options, arguments, and expected behavior. 
The CLI interface is the same across all models.

## Overview

The CLI is the main entry point to interact with the model runtime, manage engines, configure settings, and run interactive sessions. 
The syntax used is:

```shell
<my-model> [COMMAND] [OPTIONS] [ARGS]
```

## Global Options

To show usage information about any of the commands, set the `--help` flag.

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

Change configuration options including reading, and setting persistent parameters.

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

YAML is the default output format, but JSON is also available.
You can use the `--format <output-type>` flag to specify a format.

Output includes:

* Supported devices
* Default configurations


### `use-engine`

To switch to a specific engine:
```shell
<my-model> use-engine <engine>
```

To automatically select the best engine:
```shell
<my-model> use-engine --auto
```

Switching engines will trigger the download of engine resources (as snap components).
A confirmation prompt will appear before downloading.

For regular installations, the `--auto` flag is applied by default.
This command can also be used to revert to the original engine after manual changes, driver updates, or new hardware installation.

