package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

type chatClient struct {
	baseUrl   string
	modelName string
	verbose   bool

	oaiClient openai.Client
}

func ChatClient(baseUrl string, modelName string, verbose bool) *chatClient {
	c := &chatClient{
		baseUrl:   baseUrl,
		modelName: modelName,
		verbose:   verbose,
	}

	c.oaiClient = openai.NewClient(option.WithBaseURL(c.baseUrl))

	return c
}

func (c *chatClient) Start() error {

	fmt.Printf("Using server at %v\n", c.baseUrl)

	// Check if server is reachable
	if err := c.handshake(); err != nil {
		return err
	}

	if c.modelName == "" {
		var err error
		if err = c.lookupModelName(); err != nil {
			return err
		}
	}
	if c.verbose {
		fmt.Printf("Using model %v\n", c.modelName)
	}

	// Check if server is ready to accept chat completion requests
	if err := c.checkServerReady(); err != nil {
		return err
	}

	fmt.Println("Type your prompt, then ENTER to submit. CTRL-C to quit.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt:            color.RedString("» "),
		InterruptPrompt:   "^C",
		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			// Block Ctrl+Z (suspend process)
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		return fmt.Errorf("error initializing readline: %v", err)
	}
	defer rl.Close()
	//rl.CaptureExitSignal() // Should readline capture and handle the exit signal? - Can be used to interrupt the chat response stream.
	log.SetOutput(rl.Stderr())

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant."),
		},
		Model: c.modelName,
	}

	for {
		prompt, err := rl.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(prompt) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if prompt == "exit" {
			break
		}

		if len(prompt) > 0 {
			params, err = c.handlePrompt(params, prompt)
			if err != nil {
				return err
			}
		}
	}
	fmt.Println("Closing chat")

	return nil
}

func (c *chatClient) handshake() error {
	stopProgress := StartProgressSpinner("Connecting to server")
	defer stopProgress()

	parsedURL, err := url.Parse(c.baseUrl)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), 5*time.Second)
	if errors.Is(err, syscall.ECONNREFUSED) {
		return fmt.Errorf("connection refused\n\n%s\n%s",
			SuggestServerStartup(),
			SuggestServerLogs())
	} else if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func (c *chatClient) checkServerReady() error {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Are you up?"),
		},
		Model:               c.modelName,
		MaxCompletionTokens: openai.Int(1),
		MaxTokens:           openai.Int(1), // for runtimes that don't yet support MaxCompletionTokens
	}

	stopProgress := StartProgressSpinner("Waiting for server to be ready")
	defer stopProgress()

	const (
		retryInterval = 5 * time.Second
		waitTimeout   = 60 * time.Second
	)
	start := time.Now()
	for {
		_, err := c.oaiClient.Chat.Completions.New(context.Background(), params)
		if err != nil {
			var apiError *openai.Error
			if errors.As(err, &apiError) {
				// llama-server starting up
				// Error: POST "http://localhost:8328/v1/chat/completions": 503 Service Unavailable {"message":"Loading model","type":"unavailable_error","code":503}
				if apiError.StatusCode == http.StatusServiceUnavailable && apiError.Type == "unavailable_error" {
					if time.Since(start) > waitTimeout {
						// Stop waiting
						return fmt.Errorf("no models available on server\n\n%s\n%s",
							SuggestServerStartup(),
							SuggestServerLogs())
					}
					time.Sleep(retryInterval)
					continue
				}
				return fmt.Errorf("api: %s", apiError.Error())
			} else {
				return fmt.Errorf("%s\n\n%s", err,
					SuggestServerLogs())
			}
		}

		return nil
	}
}

func (c *chatClient) lookupModelName() error {
	stopProgress := StartProgressSpinner("Looking up model name")
	defer stopProgress()

	modelService := openai.NewModelService(option.WithBaseURL(c.baseUrl))

	const (
		retryInterval = 5 * time.Second
		waitTimeout   = 60 * time.Second
	)
	start := time.Now()
	for {
		modelPage, err := modelService.List(context.Background())
		if err != nil {
			var apiError *openai.Error
			if errors.As(err, &apiError) {
				// llama-server starting up
				// Error: GET "http://localhost:8330/v1/models": 503 Service Unavailable {"message":"Loading model","type":"unavailable_error","code":503}
				if apiError.StatusCode == http.StatusServiceUnavailable && apiError.Type == "unavailable_error" {
					if time.Since(start) > waitTimeout {
						// Stop waiting
						return fmt.Errorf("no models available on server\n\n%s\n%s",
							SuggestServerStartup(),
							SuggestServerLogs())
					}
					time.Sleep(retryInterval)
					continue
				}
				return fmt.Errorf("api: %s", apiError.Error())
			}

			return fmt.Errorf("%s\n\n%s", err,
				SuggestServerLogs())
		}

		if len(modelPage.Data) == 0 {
			// This can happen when OpenVINO Model Server is starting up
			if time.Since(start) > waitTimeout {
				// Stop waiting
				return fmt.Errorf("server returned no models\n\n%s\n%s",
					SuggestServerStartup(),
					SuggestServerLogs())
			}
			time.Sleep(retryInterval)
			continue
		} else if len(modelPage.Data) > 1 {
			var names []string
			for _, model := range modelPage.Data {
				names = append(names, model.ID)
			}
			return fmt.Errorf("expected one but server returned multiple models: %s", strings.Join(names, ", "))
		}

		c.modelName = modelPage.Data[0].ID
		return nil
	} // end for
}

func (c *chatClient) handlePrompt(params openai.ChatCompletionNewParams, prompt string) (openai.ChatCompletionNewParams, error) {
	params.Messages = append(params.Messages, openai.UserMessage(prompt))

	paramDebugString, _ := json.Marshal(params)

	if c.verbose {
		fmt.Printf("Sending request: %s\n", paramDebugString)
	}

	stopProgress := StartProgressSpinner("Waiting for a response")
	stream := c.oaiClient.Chat.Completions.NewStreaming(context.Background(), params)
	stopProgress()

	appendParam, err := c.processStream(stream)
	if err != nil {
		return params, err
	}

	// Store previous prompts for context
	if appendParam != nil {
		params.Messages = append(params.Messages, *appendParam)
	}
	fmt.Println()

	return params, nil
}

func (c *chatClient) processStream(stream *ssestream.Stream[openai.ChatCompletionChunk]) (*openai.ChatCompletionMessageParamUnion, error) {
	// optionally, an accumulator helper can be used
	acc := openai.ChatCompletionAccumulator{}

	// An opening <think> tag will change the output color to indicate reasoning.
	thinking := false

	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)

		if _, ok := acc.JustFinishedContent(); ok {
			//fmt.Println("\nContent stream finished")
		}

		// if using tool calls
		if tool, ok := acc.JustFinishedToolCall(); ok {
			fmt.Printf("Tool call stream finished %d: %s %s", tool.Index, tool.Name, tool.Arguments)
		}

		if refusal, ok := acc.JustFinishedRefusal(); ok {
			fmt.Printf("Refusal stream finished: %s", refusal)
		}

		// Print chunks as they are received
		if len(chunk.Choices) > 0 {
			lastChunk := chunk.Choices[0].Delta.Content

			if strings.Contains(lastChunk, "<think>") {
				thinking = true
				fmt.Printf("%s", color.BlueString(lastChunk))
			} else if strings.Contains(lastChunk, "</think>") {
				thinking = false
				fmt.Printf("%s", color.BlueString(lastChunk))

			} else if thinking {
				fmt.Printf("%s", color.BlueString(lastChunk))

			} else {
				fmt.Printf("%s", lastChunk)
			}
		}
	}

	if err := stream.Err(); err != nil {
		if errors.Is(err, syscall.ECONNREFUSED) { // connection refused before streaming
			return nil, fmt.Errorf("connection refused\n\n%s",
				SuggestServerLogs())
		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			fmt.Println() // break the line after incomplete stream
			return nil, fmt.Errorf("connection closed by server\n\n%s",
				SuggestServerLogs())
		}
		return nil, fmt.Errorf("%s\n\n%s", err,
			SuggestServerLogs())
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	appendParam := acc.Choices[0].Message.ToParam()
	if acc.Choices[0].Message.Content == "" {
		return nil, nil
	}
	return &appendParam, nil
}
