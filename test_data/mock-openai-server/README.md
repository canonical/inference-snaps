# Mock OpenAI API Server

A minimal Python HTTP server that mimics a subset of the [OpenAI REST API](https://platform.openai.com/docs/api-reference) for local testing. No third-party dependencies are required — it runs on the Python 3 standard library.

## Endpoints

Both `/v1` and `/v3` prefixes are served and return identical responses.
All responses include `Access-Control-Allow-Origin: *` CORS headers.

| Method | Path | Description |
|--------|------|-------------|
| `GET`  | `/{v1,v3}/models` | Returns a list of mock models |
| `POST` | `/{v1,v3}/chat/completions` | Returns a chat completion response (streaming or non-streaming) |

### Streaming vs non-streaming

`/chat/completions` inspects the `"stream"` field in the request body:

- `"stream": true` → responds with `Content-Type: text/event-stream` SSE chunks (`ChatCompletionChunk` objects), terminated by `data: [DONE]`
- `"stream": false` or omitted → responds with `Content-Type: application/json` (`ChatCompletion` object)

## Usage

```bash
python3 server.py [--host HOST] [--port PORT] [--delay SECONDS]
```

### Arguments

| Argument | Default | Description |
|----------|---------|-------------|
| `--host` | `127.0.0.1` | Network interface to bind to |
| `--port` | `8080` | TCP port to listen on |
| `--delay` | `0.5` | Seconds to wait before each response and between SSE chunks |

### Examples

```bash
# Listen on localhost port 8080 (default)
python3 server.py

# Listen on all interfaces, port 11434, no delay
python3 server.py --host 0.0.0.0 --port 11434 --delay 0

# Simulate a slow server (2 s per token)
python3 server.py --delay 2.0
```

## Sample Responses

### GET /v1/models

```json
{
  "object": "list",
  "data": [
    {
      "id": "mock-model",
      "object": "model",
      "created": 1712345678,
      "owned_by": "mock"
    }
  ]
}
```

### POST /v1/chat/completions (non-streaming)

```json
{
  "id": "chatcmpl-mock-0000000000000001",
  "object": "chat.completion",
  "created": 1712345678,
  "model": "mock-model",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Hello! I am a mock assistant. How can I help you today?"
      },
      "finish_reason": "stop",
      "logprobs": null
    }
  ],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 12,
    "total_tokens": 22
  },
  "system_fingerprint": null
}
```

### POST /v1/chat/completions (streaming)

```
data: {"id":"chatcmpl-mock-0000000000000001","object":"chat.completion.chunk","model":"mock-model","choices":[{"index":0,"delta":{"role":"assistant","content":"Hello!"},"finish_reason":null}]}

data: {"id":"chatcmpl-mock-0000000000000001","object":"chat.completion.chunk","model":"mock-model","choices":[{"index":0,"delta":{"content":" I"},"finish_reason":null}]}

...

data: {"id":"chatcmpl-mock-0000000000000001","object":"chat.completion.chunk","model":"mock-model","choices":[{"index":0,"delta":{"content":""},"finish_reason":"stop"}]}

data: [DONE]
```

## Testing with curl

```bash
# List models
curl http://127.0.0.1:8080/v1/models | python3 -m json.tool

# Non-streaming chat completion
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"mock-model","messages":[{"role":"user","content":"Hello"}]}' \
  | python3 -m json.tool

# Streaming chat completion
curl http://127.0.0.1:8080/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{"model":"mock-model","stream":true,"messages":[{"role":"user","content":"Hello"}]}'
```
