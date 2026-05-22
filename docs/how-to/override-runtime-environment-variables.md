
(override-runtime-environment-variables)=
# Override runtime environment variables

Overriding runtime environment variables allows you to customize the behavior of the inference snap at runtime.
This tutorial builds on the concepts introduced in {ref}`Configure inference snaps <configure-inference-snaps>`.

## Setting an environment variable

You can set an environment variable for the inference snap by using the CLI tool.
Use the following syntax:

```shell
sudo <inference-snap> set env.<ENV_VAR_NAME>=<value>
```

Replace `<ENV_VAR_NAME>` with the name of the environment variable you want to set, while following these rules:
- Environment variable names are always uppercase, so `<ENV_VAR_NAME>` is converted to uppercase if needed.
- The only allowed special character is `-`, which is converted to `_` in the actual environment variable name.

For example, `my-env-var` is set as `MY_ENV_VAR`.

## Examples

Here are some examples of how to use this feature.

### Sideloading a model

First, inspect the current environment variables in the snap:

```shell
sudo <inference-snap> run -- env
```

This command prints all environment variables available in the snap.
Now suppose you want to sideload a model that requires a few environment variables:

```shell
sudo <inference-snap> set env.MODEL_FILE=/path/to/my_model
sudo <inference-snap> set env.MODEL_NAME=my_model
sudo <inference-snap> set env.MMPROJ_FILE=/path/to/my_model.mmproj # if the model also requires an mmproj file
```

Restart the snap to apply the changes:

```shell
sudo snap restart <inference-snap>
```

### Getting supported environment variables for `llama.cpp`

Suppose you want to interact with the `llama.cpp` engine and check which environment variables it supports.

```shell
sudo <inference-snap> run -- server.sh -h
```

This command prints the help message for the `server.sh` script, including a list of supported environment variables and their descriptions.

For example, you can set the number of layers to store in VRAM with `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.LLAMA_ARG_N_GPU_LAYERS=2
```