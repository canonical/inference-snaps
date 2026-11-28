package common

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"

	"github.com/canonical/inference-snaps/modelctl/v2/internal/fallbackserver"
)

func processTestStream(body string) (*openai.ChatCompletionMessageParamUnion, error) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, body)
	}))
	defer server.Close()

	c := chatClient{oaiClient: openai.NewClient(option.WithBaseURL(server.URL), option.WithMaxRetries(0))}
	stream := c.oaiClient.Chat.Completions.NewStreaming(context.Background(), openai.ChatCompletionNewParams{})
	return c.processStream(stream)
}

func TestProcessStreamEmpty(t *testing.T) {
	wantErr := "server returned an empty response\n\n" + SuggestServerLogs()
	for name, body := range map[string]string{
		"no events": "",
		"done only": "data: [DONE]\n\n",
	} {
		t.Run(name, func(t *testing.T) {
			msg, err := processTestStream(body)
			if err == nil || err.Error() != wantErr {
				t.Fatalf("expected error %q, got %v (message %v)", wantErr, err, msg)
			}
		})
	}
}

func TestProcessStreamEmptyContent(t *testing.T) {
	msg, err := processTestStream("data: {\"choices\":[{\"index\":0,\"delta\":{\"content\":\"\"}}]}\n\ndata: [DONE]\n\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg != nil {
		t.Fatalf("expected no message, got %v", msg)
	}
}

func TestProcessStreamContent(t *testing.T) {
	msg, err := processTestStream("data: {\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"Hello\"}}]}\n\ndata: [DONE]\n\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg == nil || msg.OfAssistant == nil || msg.OfAssistant.Content.OfString.Value != "Hello" {
		t.Fatalf("expected assistant message Hello, got %v", msg)
	}
}

type testReply struct {
	status  int
	headers map[string]string
	body    string
}

func serveReplies(t *testing.T, replies ...testReply) (string, *atomic.Int32) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		reply := replies[min(int(requests.Add(1)), len(replies))-1]
		for key, value := range reply.headers {
			w.Header().Set(key, value)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(reply.status)
		_, _ = fmt.Fprint(w, reply.body)
	}))
	t.Cleanup(server.Close)
	return server.URL, &requests
}

func TestWaitForServer(t *testing.T) {
	fallback := testReply{
		status:  http.StatusServiceUnavailable,
		headers: map[string]string{"X-Should-Retry": "false"},
		body:    `{"error":{"message":"No active engine is set.","type":"unavailable_error","param":null,"code":"` + fallbackserver.ErrorCode + `"}}`,
	}
	loading := testReply{
		status: http.StatusServiceUnavailable,
		body:   `{"error":{"code":503,"message":"Loading model","type":"unavailable_error"}}`,
	}

	for name, wait := range map[string]struct {
		modelName string
		ready     string
		run       func(c *chatClient) error
	}{
		"lookup model name": {
			ready: `{"object":"list","data":[{"id":"example","object":"model"}]}`,
			run:   (*chatClient).lookupModelName,
		},
		"check server ready": {
			modelName: "example",
			ready:     `{"id":"1","object":"chat.completion","model":"example","choices":[]}`,
			run:       (*chatClient).checkServerReady,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Run("fallback server", func(t *testing.T) {
				t.Parallel()
				url, requests := serveReplies(t, fallback)
				err := wait.run(ChatClient(url, wait.modelName, false))
				wantErr := "No active engine is set.\n\n" + SuggestServerLogs()
				if err == nil || err.Error() != wantErr {
					t.Fatalf("expected error %q, got %v", wantErr, err)
				}
				if got := requests.Load(); got != 1 {
					t.Fatalf("expected 1 request, got %d", got)
				}
			})

			t.Run("model loading", func(t *testing.T) {
				t.Parallel()
				url, _ := serveReplies(t, loading, loading, loading, testReply{status: http.StatusOK, body: wait.ready})
				c := ChatClient(url, wait.modelName, false)
				if err := wait.run(c); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if c.modelName != "example" {
					t.Fatalf("expected model example, got %q", c.modelName)
				}
			})
		})
	}
}
