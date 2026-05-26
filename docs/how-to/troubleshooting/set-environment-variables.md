(set-environment-variables)=
# Environment variables for the inference snap

Runtime environment variables can be used to customize inference snap behavior.
This guide builds on the concepts introduced in {ref}`Configure inference snaps <configure-inference-snaps>`.

## Set environment variables

Set an environment variable with:

```shell
sudo <inference-snap> set env.<env-var-name>=<value>
```

Replace `<env-var-name>` with the variable name. Naming rules:
- Environment variable names are always uppercase, so `<env-var-name>` is converted to uppercase if needed.
- The only allowed special character is `-`, which is converted to `_` in the actual environment variable name.

For example, `my-env-var` is converted to `MY_ENV_VAR`.

New variables can be added, and existing variables can be overridden.
To list currently set variables, run:

```shell
sudo <inference-snap> run -- env
```

## Examples

### Interact with `llama.cpp` environment variables

For `llama.cpp`, supported variables are documented [here](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md#usage).

For example, it is possible to set the number of layers stored in VRAM with `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.LLAMA_ARG_N_GPU_LAYERS=2
```