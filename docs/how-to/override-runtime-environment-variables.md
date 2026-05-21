
(override-runtime-environment-variables)=
# Override runtime environment variables

Overriding runtime environment variables allows you to customize the behavior of the inference snap att runtime.
This tutorial relies on knowledge acquired in the {ref}`configure snap tutorial <configure-snap>`.

## Setting an env var

It is possible to set an environment variable for the inference snap using the CLI tool.
The syntax for this is as follows:

```shell
sudo <inference-snap> set env.<ENV_VAR_NAME>=<value>
```

<ENV_VAR_NAME> should be replaced with the name of the environment variable you want to set, considering the following rules:
- The environment variables are always uppercase, so <ENV_VAR_NAME> will be parsed and converted to uppercase if needed
- Only `-` are allowed as special characters, and they will be converted to `_` in the actual environment variable.

 For example, `my-env-var` will be set as `MY_ENV_VAR`.

## Examples

Here are some examples of how to use this feature.

### Sideloading model

First of all, let's inspect current environment variables in the snap:

```shell
sudo <inference-snap> run -- env
```

This will print all environment variables available in the snap.
Now, let's say we want to sideload a model that requires us to set a few environment variables:

```shell
sudo <inference-snap> set env.MODEL_FILE=/path/to/my_model
sudo <inference-snap> set env.MODEL_NAME=my_model
sudo <inference-snap> set env.MMPROJ_FILE=/path/to/my_model.mmproj # if the model also requires an mmproj file
```

Now we can restart the snap to apply the changes:

```shell
sudo snap restart <inference-snap>
```

### Getting supported env vars for llama.cpp

Let's say we want to interact with the llama.cpp engine and want to know which environment variables it supports.

```shell
sudo gemma4 run -- server.sh -h
```

This will print the help message for the server.sh script, which includes a list of supported environment variables and their description.

For example, it is possible to set the number of layers to store in vRAM using the `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.LLAMA_ARG_N_GPU_LAYERS=2
```