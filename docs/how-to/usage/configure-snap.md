(configure-inference-snaps)=
# Configure inference snaps

Inference snaps expose a number of configurations.
You can use the CLI to view these or set a new value for a configuration.
The available configurations depend on the active engine.

## View available configurations

Use the {ref}`CLI tool's <models-cli>` `get` command to show all currently available configurations:

```
<inference-snap> get
```

Example:
```{terminal}
:input: qwen-vl get
http.host: 127.0.0.1
http.port: 8326
verbose: false
```

## Get a configuration value

To get the current value of a configuration, pass the configuration key to the `get` command:

```
<inference-snap> get <key>
```

Example:
```{terminal}
:input: qwen-vl get http
http.base-path: v1
http.host: 127.0.0.1
http.port: 9999
:input: qwen-vl get http.port
9999
```

## Set a configuration value

You can change the value of a specific configuration using the `set` command.
This requires root permission.

```shell
sudo <inference-snap> set <key>=<new-value>
```

You'll need to restart the server to apply the changes after setting a new configuration value.

```shell
sudo snap restart <inference-snap>
```

Example:
```{terminal}
:input: sudo qwen-vl set http.port=9090

:input: sudo snap restart qwen-vl
Restarted.
```
