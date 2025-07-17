# Use an AI model snap via its OpenAI API

Some AI model snaps expose an [OpenAI](https://github.com/openai/openai-openapi) compliant network API.
This guide shows how to find the configurations needed to access this API.

## API URL

The model server binds to the host interface and port that are configured as snap options.
These snap options can be viewed as follows:

```{terminal}
:input: sudo snap get <ai-snap> http
Key             Value
http.base-path  v1
http.host       127.0.0.1
http.port       8080
```

The base path of the OpenAI compliant API might be different between stacks.
This is also reported here as a snap configuration option.

The URL for the API should be built up using these options, in the format: `http://localhost:<http.port>/<http.base-path>/`.
In this example that would be http://localhost:8080/v1/.
The API can be accessed from your local computer using this URL.

If the `http.host` value is set to `127.0.0.1`, the model can only be accessed from the local device.
Set it to `0.0.0.0` if you want to access it from other devices on the network.
From other devices, the full URL will be similar to above, but with `localhost` replaced by the IP address or hostname of the device where the model is running.

## Model name

Some servers require you to specify the exact model name when querying it via the API.
An OpenAI compliant API lists available models under `<api-url>/models`.
If the server does not implement this endpoint, the model name can also be looked up in the snap options:

```{terminal}
:input: sudo snap get <ai-model> 
Key            Value
http           {...}
model-name     <model-name>
...
```

## Test with `curl`

If your installed model supports chat completions, you can use its API URL to make a test prompt.
In this example, replace the URL and `<model-name>` fields with the values you obtained previously.

```sh
curl http://localhost:8080/v1/chat/completions \
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
  "model": "<ai-model>",
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

## Using other clients

Other OpenAI API compatible clients can also make use of these snaps.
Configure these clients with the API URL and the model name obtained above.
