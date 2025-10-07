(use-openai-api)=
# Use an AI model snap via its OpenAI API

Some AI model snaps expose an [OpenAI](https://github.com/openai/openai-openapi) compliant network API.
This guide shows how to access this API.

## Identify API URL

The OpenAI API URL is reported by the `status` command:

```shell
<ai-snap> status
```

The URL is in the format `http://localhost:<port>/<base-path>`, where the `<base-path>` depends on the currently active engine.
This URL can be used to access the model from your local computer.

## API interface and port

The model server binds to the configured host interface and port.
These configurations can be viewed with:

```shell
<ai-snap> get http
```

It can also be changed by using the `set` command.
A server restart is required after changing these values.

```shell
<ai-snap> set http.port=9090
sudo snap restart <ai-snap>
```

If the `http.host` value is set to `127.0.0.1`, the model can only be accessed from the local machine.
Set it to `0.0.0.0` if you want to access it from other machines on the network:

```shell
<ai-snap> set http.host=0.0.0.0
sudo snap restart <ai-snap>
```

To access the API from elsewhere, use the IP address or hostname of the machine where the model is running, instead of `localhost`.


## Look up model name

Some servers require you to specify the exact model name when querying it via the API.
An OpenAI compliant API lists available models under `<api-url>/models`.

```shell
curl <api-url>/models
```


## Test API access with `curl`

If your installed model supports chat completions, you can use its API URL to make a test prompt.
In this example, replace the `<api-url>` and `<model-name>` fields with the values you obtained previously.

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

Other OpenAI API compatible clients can also make use of these snaps.
Configure these clients with the API URL and the model name obtained above.

If the client requires you to provide an API key, you can leave this blank or provide a placeholder value.
