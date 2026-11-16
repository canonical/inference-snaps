(set-environment-variables)=

# Set environment variables

```{versionadded} 1.0.0-beta.42
```

```{versionchanged} 2.0.0
`passthrough.environment.` key prefix has been replaced with `env.`
```

Inference snaps should typically be customized using {ref}`configurations <configure-inference-snaps>`.  
Runtime environment variables can be used for debugging and to tweak a snap beyond what is possible via configurations. 

## Syntax

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
sudo <inference-snap> run env
```

## Examples

### Override `llama-server` environment variables

For `llama-server`, supported variables are documented [here](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md#usage).

For example, it is possible to set the number of layers stored in VRAM with `LLAMA_ARG_N_GPU_LAYERS`:

```shell
sudo <inference-snap> set env.llama-arg-n-gpu-layers=2
```

### Embed inputs longer than 512 tokens

When using `llama-server` for embeddings, each input must fit in the physical batch (`n_ubatch`).
The default physical batch size is 512 tokens, even if the model supports a larger context.
An input exceeding this limit can fail with HTTP 500 and a message asking you to increase the physical batch size.
Splitting a request into fewer inputs does not help: the limit applies to each input, not the total request.

Set `LLAMA_ARG_UBATCH` (physical batch size) and `LLAMA_ARG_BATCH` (logical batch size) to accommodate the model's context.
For example, to embed inputs up to the reported 2048-token context of `nomic-embed-text-v1-5`:

```shell
sudo nomic-embed-text-v1-5 set --no-restart env.llama-arg-batch=2048 env.llama-arg-ubatch=2048
sudo snap restart nomic-embed-text-v1-5
```

Verify the environment passed to the runtime:

```shell
sudo nomic-embed-text-v1-5 run printenv LLAMA_ARG_BATCH LLAMA_ARG_UBATCH
```

Both values should be `2048`.
Retry the embedding request after the server is ready.
Larger batches can increase memory usage.
These settings do not extend the model's context window; split or truncate longer texts before submitting them.

To restore the defaults, unset the overrides and restart the service:

```shell
sudo nomic-embed-text-v1-5 unset --no-restart env.llama-arg-batch
sudo nomic-embed-text-v1-5 unset --no-restart env.llama-arg-ubatch
sudo snap restart nomic-embed-text-v1-5
```