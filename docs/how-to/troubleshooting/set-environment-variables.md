(set-environment-variables)=
# Set environment variables for the inference snap

Setting runtime environment variables allows you to customize the behavior of the inference snap at runtime.
This guide builds on the concepts introduced in {ref}`Configure inference snaps <configure-inference-snaps>`.

## Setting an environment variable

You can set an environment variable for the inference snap by using the CLI tool.
Use the following syntax:

```shell
sudo <inference-snap> set env.<env-var-name>=<value>
```

Replace `<env-var-name>` with the name of the environment variable you want to set, while following these rules:
- Environment variable names are always uppercase, so `<env-var-name>` is converted to uppercase if needed.
- The only allowed special character is `-`, which is converted to `_` in the actual environment variable name.

For example, `my-env-var` is converted to `MY_ENV_VAR`.

It is possible to add a new environment variable that is not currently set inside the snap or override an existing one.
To check which environment variables are currently set in the snap, run:

```shell
sudo <inference-snap> run -- env
```

## Examples

Here are some examples of how to use this feature.

### Getting supported environment variables for `llama.cpp`

Suppose you want to interact with the `llama.cpp` engine and check which environment variables it supports.

```shell
sudo snap run <inference-snap>.server -h
```

This command prints the help message of `llama.cpp`, including a list of supported environment variables and their descriptions.

For example, you can set the number of layers to store in VRAM with `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.LLAMA_ARG_N_GPU_LAYERS=2
```