(configure-snap)=
# Configure Inference Snaps

Inference snaps expose a number of configurations.
These configurations depend on the active engine.

## Explore available configurations

The CLI implements a `get` command which will show all currently available configurations:

```
<inference-snap> get
```

Example:
```{terminal}
:input: qwen-vl get
http.host: 127.0.0.1
http.port: 9999
verbose: false
```

## Get a configuration value

A configuration key can be passed to the `get` command to get the current value.

```
<inference-snap> get <key>
```

Examples:
```{terminal}
:input: qwen-vl get http
http.base-path: v1
http.host: 127.0.0.1
http.port: 9999
:input: qwen-vl get http.port
9999
```

## Set a configuration value

The value of a specific configuration can be changed with the set command.
This requires root permission.

```shell
sudo <inference-snap> set <key>=<new-value>
```

After changing a configuration, the server needs to be restarted to apply the configuration changes.

```shell
sudo snap restart <inference-snap>
```

Example:
```{terminal}
:input: sudo qwen-vl set http.port=9090

:input: sudo snap restart qwen-vl
Restarted.
```
