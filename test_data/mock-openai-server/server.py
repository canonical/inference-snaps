#!/usr/bin/env python3
"""
Mock OpenAI-compatible API server for testing.

Serves a subset of the OpenAI API spec under both /v1 and /v3 prefixes.

Endpoints:
  GET  /{v1,v3}/models
  POST /{v1,v3}/chat/completions

Usage:
  python3 server.py --host 127.0.0.1 --port 8080
"""

import argparse
import json
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

# Global delay in seconds applied before every response and between SSE chunks.
RESPONSE_DELAY = 0.5

# Delay in seconds before the first SSE chunk is sent (time to first token).
TIME_TO_FIRST_TOKEN = 0.0

# When True, responses include reasoning_content in addition to content.
INCLUDE_REASONING = False


# ---------------------------------------------------------------------------
# Static response payloads
# ---------------------------------------------------------------------------

def _models_payload():
    return {
        "object": "list",
        "data": [
            {
                "id": "mock-model",
                "object": "model",
                "created": int(time.time()),
                "owned_by": "mock",
                "permission": [],
                "root": "mock-model",
                "parent": None,
            },
        ],
    }


def _extract_last_user_message(request_json):
    """Return the content of the last user message in the request, or a fallback."""
    messages = request_json.get("messages", [])
    for msg in reversed(messages):
        if msg.get("role") == "user":
            return msg.get("content", "")
    return ""


def _chat_completion_payload(reply):
    """Plain (non-streaming) chat.completion response."""
    message = {
        "role": "assistant",
        "content": reply,
    }
    if INCLUDE_REASONING:
        message["reasoning_content"] = reply
    return {
        "id": "chatcmpl-mock-0000000000000001",
        "object": "chat.completion",
        "created": int(time.time()),
        "model": "mock-model",
        "choices": [
            {
                "index": 0,
                "message": message,
                "finish_reason": "stop",
                "logprobs": None,
            }
        ],
        "usage": {
            "prompt_tokens": 10,
            "completion_tokens": len(reply.split()),
            "total_tokens": 10 + len(reply.split()),
        },
        "system_fingerprint": None,
    }


def _chat_completion_chunk(content=None, reasoning_content=None, finish_reason=None, is_first=False):
    """Build a single ChatCompletionChunk payload."""
    delta = {}
    if is_first:
        delta["role"] = "assistant"
    if content is not None:
        delta["content"] = content
    if reasoning_content is not None:
        delta["reasoning_content"] = reasoning_content

    return {
        "id": "chatcmpl-mock-0000000000000001",
        "object": "chat.completion.chunk",
        "created": int(time.time()),
        "model": "mock-model",
        "choices": [
            {
                "index": 0,
                "delta": delta,
                "finish_reason": finish_reason,
                "logprobs": None,
            }
        ],
        "system_fingerprint": None,
    }


def _chat_completion_sse_chunks(reply):
    """
    Yield SSE-formatted lines for a complete streaming chat response.
    Each item is a bytes object ready to be written to the socket.
    First emits all reasoning_content chunks, then all content chunks.
    """
    words = reply.split(" ")

    # Time-to-first-token delay before any chunk is sent
    if TIME_TO_FIRST_TOKEN > 0:
        time.sleep(TIME_TO_FIRST_TOKEN)

    # Phase 1: reasoning_content chunks (only when reasoning is enabled)
    if INCLUDE_REASONING:
        for i, word in enumerate(words):
            time.sleep(RESPONSE_DELAY)
            token = word if i == 0 else " " + word
            chunk = _chat_completion_chunk(reasoning_content=token, is_first=(i == 0))
            yield f"data: {json.dumps(chunk)}\n\n".encode("utf-8")

    # Phase 2: content chunks
    for i, word in enumerate(words):
        time.sleep(RESPONSE_DELAY)
        token = word if i == 0 else " " + word
        chunk = _chat_completion_chunk(content=token, is_first=(i == 0 and not INCLUDE_REASONING))
        yield f"data: {json.dumps(chunk)}\n\n".encode("utf-8")

    # Final chunk: empty delta, finish_reason="stop"
    time.sleep(RESPONSE_DELAY)
    stop_chunk = _chat_completion_chunk(content="", finish_reason="stop")
    yield f"data: {json.dumps(stop_chunk)}\n\n".encode("utf-8")

    # SSE stream terminator
    yield b"data: [DONE]\n\n"


# ---------------------------------------------------------------------------
# Request handler
# ---------------------------------------------------------------------------

class MockOpenAIHandler(BaseHTTPRequestHandler):

    # Suppress default request logging; override to customise if desired.
    def log_message(self, fmt, *args):  # noqa: N802
        print(f"[mock-openai] {self.address_string()} - {fmt % args}")

    # ------------------------------------------------------------------
    # Helpers
    # ------------------------------------------------------------------

    def _send_cors_headers(self):
        self.send_header("Access-Control-Allow-Origin", "*")
        self.send_header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        self.send_header("Access-Control-Allow-Headers", "Content-Type, Authorization")

    def _send_json(self, payload, status=200):
        time.sleep(RESPONSE_DELAY)
        body = json.dumps(payload, indent=2).encode("utf-8")
        self.send_response(status)
        self._send_cors_headers()
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _send_error_json(self, status, message, error_type="invalid_request_error"):
        payload = {
            "error": {
                "message": message,
                "type": error_type,
                "param": None,
                "code": None,
            }
        }
        self._send_json(payload, status)

    def _strip_api_prefix(self, path):
        """
        Strip /v1 or /v3 prefix and return (version, remainder).
        Returns (None, path) if no recognised prefix is found.
        """
        for version in ("/v1", "/v3"):
            if path == version or path.startswith(version + "/"):
                remainder = path[len(version):]
                return version, remainder or "/"
        return None, path

    def _route(self, method):
        """Dispatch a request to the appropriate handler."""
        path = self.path.split("?")[0]  # ignore query string
        version, route = self._strip_api_prefix(path)

        if version is None:
            self._send_error_json(404, f"No route matched: {path}")
            return

        if method == "GET" and route == "/models":
            self._send_json(_models_payload())
        elif method == "POST" and route == "/chat/completions":
            self._handle_chat_completions()
        else:
            self._send_error_json(
                404,
                f"Unknown endpoint: {method} {path}",
            )

    # ------------------------------------------------------------------
    # HTTP verb handlers
    # ------------------------------------------------------------------

    def do_GET(self):  # noqa: N802
        self._route("GET")

    def do_POST(self):  # noqa: N802
        self._route("POST")

    def do_OPTIONS(self):  # noqa: N802
        """Return CORS pre-flight headers so browser-based clients work."""
        self.send_response(204)
        self._send_cors_headers()
        self.end_headers()

    # ------------------------------------------------------------------
    # Endpoint implementations
    # ------------------------------------------------------------------

    def _handle_chat_completions(self):
        # Read the request body so we can inspect the "stream" field.
        content_length = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(content_length) if content_length > 0 else b"{}"

        try:
            request_json = json.loads(body)
        except json.JSONDecodeError:
            request_json = {}

        reply = _extract_last_user_message(request_json)

        if request_json.get("stream", False):
            # Streaming response: Server-Sent Events (SSE) of ChatCompletionChunk
            self.send_response(200)
            self._send_cors_headers()
            self.send_header("Content-Type", "text/event-stream")
            self.send_header("Cache-Control", "no-cache")
            self.end_headers()
            try:
                for sse_line in _chat_completion_sse_chunks(reply):
                    self.wfile.write(sse_line)
                self.wfile.flush()
            except BrokenPipeError:
                pass  # client disconnected mid-stream
        else:
            # Non-streaming response: plain application/json ChatCompletion
            self._send_json(_chat_completion_payload(reply))


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------

def main():
    global RESPONSE_DELAY, TIME_TO_FIRST_TOKEN, INCLUDE_REASONING

    parser = argparse.ArgumentParser(
        description="Mock OpenAI-compatible HTTP API server for testing."
    )
    parser.add_argument(
        "--host",
        default="127.0.0.1",
        help="Interface to listen on (default: 127.0.0.1)",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=8080,
        help="Port to listen on (default: 8080)",
    )
    parser.add_argument(
        "--delay",
        type=float,
        default=RESPONSE_DELAY,
        metavar="SECONDS",
        help="Delay in seconds before each response and between SSE chunks (default: 0)",
    )
    parser.add_argument(
        "--ttft",
        type=float,
        default=TIME_TO_FIRST_TOKEN,
        metavar="SECONDS",
        help="Delay in seconds before the first SSE chunk is sent, i.e. time to first token (default: 0)",
    )
    parser.add_argument(
        "--reasoning",
        action="store_true",
        default=False,
        help="Include reasoning_content in responses (default: off)",
    )
    args = parser.parse_args()

    RESPONSE_DELAY = args.delay
    TIME_TO_FIRST_TOKEN = args.ttft
    INCLUDE_REASONING = args.reasoning

    server_address = (args.host, args.port)
    httpd = HTTPServer(server_address, MockOpenAIHandler)
    print(f"[mock-openai] Listening on http://{args.host}:{args.port}")
    print(f"[mock-openai] Response delay: {RESPONSE_DELAY}s")
    print(f"[mock-openai] Time to first token: {TIME_TO_FIRST_TOKEN}s")
    print(f"[mock-openai] Reasoning: {'on' if INCLUDE_REASONING else 'off'}")
    print("[mock-openai] Serving /v1 and /v3 prefixes")
    print("[mock-openai] Endpoints:")
    print(f"  GET  http://{args.host}:{args.port}/v1/models")
    print(f"  POST http://{args.host}:{args.port}/v1/chat/completions")
    print(f"  GET  http://{args.host}:{args.port}/v3/models")
    print(f"  POST http://{args.host}:{args.port}/v3/chat/completions")
    print("[mock-openai] Press Ctrl+C to stop.")
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\n[mock-openai] Shutting down.")
        httpd.server_close()


if __name__ == "__main__":
    main()

