(use-openai-api)=
# Use an Inference Snap via its OpenAI API

Some Inference Snaps expose an [OpenAI](https://github.com/openai/openai-openapi) compliant network API.
This guide shows you how to access this API.

## Identify API URL

The OpenAI API URL is reported by the `status` command:

```shell
<inference-snap> status
```

The URL is in the format `http://localhost:<port>/<base-path>`, where `<base-path>` depends on the current active engine.
This URL can be used to access the model from your local computer.

## Configure API interface and port

The Inference Snap's server binds to the configured host interface and port.
To view these configurations:

```shell
<inference-snap> get http
```

Use the `set` command to change these configurations.
You must restart the server after making a change:

```shell
<inference-snap> set http.port=9090
sudo snap restart <inference-snap>
```

If the `http.host` value is set to `127.0.0.1`, you can only access the model from the local machine.
To allow access by other machines on the network, set it to `0.0.0.0`:

```shell
<inference-snap> set http.host=0.0.0.0
sudo snap restart <inference-snap>
```

To access the API from another machine in the network, in the URL, use the IP address or hostname of the machine where the Inference Snap is running instead of `localhost`.


## Look up model name

Some servers require you to specify the exact model name when querying them via the API.
Available models are listed under `<api-url>/models`.

```shell
curl <api-url>/models
```

The response will look similar to the following, with the `id` field containing the model name:
```json
{
  "object": "list",
  "data": [
    {
      "id": "Qwen2.5-VL-3B-Instruct-ov-int4",
      "object": "model",
      "created": 1760609472,
      "owned_by": "OVMS"
    }
  ]
}
```


## Test API access with `curl`

If your installed Inference Snap supports chat completions, you can use its API URL to make a test prompt.
In this example, replace `<api-url>` and `<model-name>` with the values you obtained.

```shell
curl <api-url>/chat/completions \
-H "Content-Type: application/json" \
-d '{
"model": "<model-name>",
"messages": [{"role": "user", "content": "Hi there!"}],
"temperature": 0.7
}' | jq
```

This should return a JSON response looking similar to this:
```json
{
  "choices": [
    {
      "finish_reason": "stop",
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! How can I assist you today?"
      }
    }
  ],
  "created": 1752591235,
  "model": "<model-name>",
  "system_fingerprint": "b1-6a746cf",
  "object": "chat.completion",
  "usage": {
    "completion_tokens": 10,
    "prompt_tokens": 11,
    "total_tokens": 21
  },
  "id": "chatcmpl-YvRT8bwvQB6YMYBjrdjMf7zumamNTIB3",
  "timings": {
    "prompt_n": 1,
    "prompt_ms": 116.811,
    "prompt_per_token_ms": 116.811,
    "prompt_per_second": 8.560837592350035,
    "predicted_n": 10,
    "predicted_ms": 824.103,
    "predicted_per_token_ms": 82.41029999999999,
    "predicted_per_second": 12.134405529405912
  }
}
```

## Use other clients

Other OpenAI API compatible clients can also make use of Inference Snaps.
Configure these clients with the API URL and the model name obtained above.

If the client requires an API key, leave this blank or enter a placeholder value.
