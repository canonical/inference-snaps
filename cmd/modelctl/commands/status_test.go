package commands

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/canonical/inference-snaps-cli/v2/cmd/modelctl/common"
	"github.com/canonical/inference-snaps-cli/v2/pkg/snap"
	"github.com/canonical/inference-snaps-cli/v2/pkg/storage"
)

// createTestContextForStatus creates a test context for Example tests
func createTestContextForStatus() *common.Context {
	testDataDir := "../../../test_data"
	enginesDir := filepath.Join(testDataDir, "engines")
	runtimesDir := filepath.Join(testDataDir, "runtimes")
	modelsDir := filepath.Join(testDataDir, "models")

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("cpu"); err != nil {
		log.Fatalf("failed to set active engine: %v", err)
	}
	if err := cache.SetActiveModel("4b-it-int4-fq-ov"); err != nil {
		log.Fatalf("failed to set active model: %v", err)
	}

	configs := map[string]string{
		"http.port":      "8080",
		"http.host":      "0.0.0.0",
		"ws.unix-socket": "/run/whisper.sock",
		// namespaced configurations
		"logger.http.port": "8081",
		"logger.http.host": "localhost",
	}

	config := storage.NewMockConfig()
	for key, value := range configs {
		if err := config.Set(key, value, storage.UserConfig); err != nil {
			log.Fatalf("failed to set %s: %v", key, err)
		}
	}

	return &common.Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      config,
		Snap:        snap.MockWithServiceStatuses(map[string]string{"llama-server": "active"}),
	}
}

func Example_statusCommand_printStatusYaml() {
	ctx := createTestContextForStatus()
	cmd := statusCommand{Context: ctx}

	statusText, err := cmd.statusYaml()
	if err != nil {
		log.Fatalf("failed to get status in yaml format: %v", err)
	}
	fmt.Print(statusText)

	// Output:
	// engine: cpu
	// services:
	//     llama-server: active
	// entrypoints:
	//     logger:
	//         url: http://localhost:8081/
	//     openai:
	//         url: http://0.0.0.0:8080/v1
	//     whisperlive:
	//         unix-socket: /run/whisper.sock (ws://unix/realtime)
	// model:
	//     name: gemma-4-4b-it-int4-fq-ov
}

func Example_statusCommand_printStatusJson() {
	ctx := createTestContextForStatus()
	cmd := statusCommand{Context: ctx}

	statusText, err := cmd.statusJson()
	if err != nil {
		log.Fatalf("failed to get status in yaml format: %v", err)
	}
	fmt.Print(statusText)

	// Output:
	// {
	//   "engine": "cpu",
	//   "services": {
	//     "llama-server": "active"
	//   },
	//   "entrypoints": {
	//     "logger": {
	//       "url": "http://localhost:8081/"
	//     },
	//     "openai": {
	//       "url": "http://0.0.0.0:8080/v1"
	//     },
	//     "whisperlive": {
	//       "unix-socket": "/run/whisper.sock",
	//       "unix-socket-url": "ws://unix/realtime"
	//     }
	//   },
	//   "model": {
	//     "name": "gemma-4-4b-it-int4-fq-ov"
	//   }
	// }
}
