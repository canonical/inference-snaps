(set-environment-variables)=
# Set environment variables

```{versionadded} 2.0.0
```
Inference snaps should typically be customized using {ref}`configurations <configure-inference-snaps>`.  
Runtime environment variables can be used for debugging and to tweak a snap beyond what is possible via configurations. 

## Set environment variables

Set an environment variable with:

```shell
sudo <inference-snap> set env.<env-var-name>=<value>
```

Replace `<env-var-name>` with the environment variable name converted from upper case to lowercase. Underscores must be replaced with hyphens.

For example, `MY_ENV_VAR` can be set by using `env.my-env-var` key.

An environment variable set this way overrides all internal values.  
To unset and get environment variables, use `unset` and `get` commands.

To list all the internally set environment variables, run:  

```shell
sudo <inference-snap> run -- env
```

## Examples

### Interact with `llama.cpp` environment variables

For `llama.cpp`, supported variables are documented [here](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md#usage).

For example, it is possible to set the number of layers stored in VRAM with `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.llama-arg-n-gpu-layers=2
```