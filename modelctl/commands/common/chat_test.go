package common

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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
