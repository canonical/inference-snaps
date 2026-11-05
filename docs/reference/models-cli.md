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

To show usage information about any command, set the `--help` flag, for example `<inference-snap> status --help`.

```{include} ../reuse/cli-commands.txt
```
