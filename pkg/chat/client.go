package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/canonical/inference-snaps-cli/cmd/cli/common"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/packages/ssestream"
)

// TODO
// Create a client object with all the reused data (baseUrl, modelName, reasoningModel, client, verbose, ...)
// This will avoid passing too many parameters around.

func Client(baseUrl string, modelName string, reasoningModel bool, verbose bool) error {
	fmt.Printf("Using server at %v\n", baseUrl)

	if modelName == "" {
		var err error
		modelName, err = findModelName(baseUrl, verbose)
		if err != nil {
			return err
		}
	}
	if verbose {
		fmt.Printf("Using model %v\n", modelName)
		if reasoningModel {
			fmt.Println("Reasoning model")
		}
	}

	// OpenAI API Client
	client := openai.NewClient(option.WithBaseURL(baseUrl))

	if err := checkServer(client, modelName); err != nil {
		if verbose {
			fmt.Printf("%v\n\n", err)
		}
		return fmt.Errorf("Unable to chat. Make sure the server has started successfully.")
	}

	fmt.Println("Type your prompt, then ENTER to submit. CTRL-C to quit.")

	rl, err := readline.NewEx(&readline.Config{
		Prompt: color.RedString("» "),
		//HistoryFile:     "/tmp/readline.tmp",
		//AutoComplete:    completer,
		InterruptPrompt: "^C",
		//EOFPrompt:       "exit", // Does not work as expected

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
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
		Model: modelName,
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
			params, err = handlePrompt(client, params, reasoningModel, prompt, verbose)
			if err != nil {
				return fmt.Errorf("error while processing prompt: %v", err)
			}
		}
	}
	fmt.Println("Closing chat")

	return nil
}

func checkServer(client openai.Client, modelName string) error {
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Are you up?"),
		},
		Model:               modelName,
		MaxCompletionTokens: openai.Int(1),
		MaxTokens:           openai.Int(1), // for runtimes that don't yet support MaxCompletionTokens
	}

	stopProgress := common.StartProgressSpinner("Connecting to server ")
	defer stopProgress()

	ctx := context.Background()
	_, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func findModelName(baseUrl string, verbose bool) (string, error) {
	stopProgress := common.StartProgressSpinner("Looking up model name ")
	defer stopProgress()

	modelService := openai.NewModelService(option.WithBaseURL(baseUrl))
	modelPage, err := modelService.List(context.Background())
	if err != nil {
		stopProgress()
		if verbose {
			fmt.Printf("%v\n\n", err)
		}
		return "", fmt.Errorf("Failed to list available models. Make sure the server has started successfully.")
	}

	if len(modelPage.Data) == 0 {
		return "", fmt.Errorf("server returned no models")
	} else if len(modelPage.Data) > 1 {
		names := make([]string, 0, len(modelPage.Data)) // Pre-allocate for efficiency
		for _, model := range modelPage.Data {
			names = append(names, model.ID)
		}
		return "", fmt.Errorf("server returned multiple models: %s", strings.Join(names, ", "))
	}

	stopProgress()
	return modelPage.Data[0].ID, nil
}

func handlePrompt(client openai.Client, params openai.ChatCompletionNewParams, reasoningModel bool, prompt string, verbose bool) (openai.ChatCompletionNewParams, error) {
	params.Messages = append(params.Messages, openai.UserMessage(prompt))

	paramDebugString, _ := json.Marshal(params)

	if verbose {
		fmt.Printf("Sending request: %s\n", paramDebugString)
	}

	stopProgress := common.StartProgressSpinner("Waiting for a response ")
	stream := client.Chat.Completions.NewStreaming(context.Background(), params)
	stopProgress()

	appendParam, err := processStream(stream, reasoningModel)
	if err != nil {
		return params, fmt.Errorf("error processing stream: %v", err)
	}

	// Store previous prompts for context
	if appendParam != nil {
		params.Messages = append(params.Messages, *appendParam)
	}
	fmt.Println()

	return params, nil
}

func processStream(stream *ssestream.Stream[openai.ChatCompletionChunk], printThinking bool) (*openai.ChatCompletionMessageParamUnion, error) {
	// optionally, an accumulator helper can be used
	acc := openai.ChatCompletionAccumulator{}

	// For reasoning models we assume the first output is them thinking, because the opening <think> tag is not always present.
	thinking := printThinking

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

	if stream.Err() != nil {
		return nil, fmt.Errorf("error reading response stream: %v", stream.Err())
	}

	// After the stream is finished, acc can be used like a ChatCompletion
	appendParam := acc.Choices[0].Message.ToParam()
	if acc.Choices[0].Message.Content == "" {
		return nil, nil
	}
	return &appendParam, nil
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
