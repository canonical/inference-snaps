package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/canonical/inference-snaps/modelctl/v2/internal/engines"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/selector"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/snap"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/testfixtures"
	"go.yaml.in/yaml/v4"
)

func TestStatusStruct(t *testing.T) {
	enginesDir := t.TempDir()
	runtimesDir := t.TempDir()
	modelsDir := t.TempDir()

	writeEngineYAML(t, enginesDir, "my-engine", `name: my-engine
runtime: my-runtime
`)
	writeRuntimeYAML(t, runtimesDir, "my-runtime", `name: my-runtime
servers:
  openai:
    protocol: http
    base-path: /v1
`)
	writeModelYAML(t, modelsDir, "my-model", `name: my-model`)

	cache := storage.NewMockCache()
	if err := cache.SetActiveEngine("my-engine"); err != nil {
		t.Fatalf("SetActiveEngine: %v", err)
	}
	if err := cache.SetActiveModel("my-model"); err != nil {
		t.Fatalf("SetActiveModel: %v", err)
	}

	config := storage.NewMockConfig()
	config.Set("http.port", "8080", storage.UserConfig)
	config.Set("http.host", "127.0.0.1", storage.UserConfig)

	ctx := &Context{
		EnginesDir:  enginesDir,
		RuntimesDir: runtimesDir,
		ModelsDir:   modelsDir,
		Cache:       cache,
		Config:      config,
		Snap:        snap.MockWithServiceStatuses(map[string]string{"llama-server": "active"}),
	}

	status, err := SnapStatus(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if status.Engine != "my-engine" {
		t.Errorf("engine: got %q, want %q", status.Engine, "my-engine")
	}
	if status.Services["llama-server"] != "active" {
		t.Errorf("service status: got %q, want %q", status.Services["llama-server"], "active")
	}
	if status.Entrypoints["openai"].Url != "http://localhost:8080/v1" {
		t.Errorf("openai entrypoint URL: got %q, want %q", status.Entrypoints["openai"].Url, "http://localhost:8080/v1")
	}
	baseURL, err := OpenAiBaseUrl(ctx)
	if err != nil {
		t.Fatalf("OpenAiBaseUrl: %v", err)
	}
	if baseURL != "http://127.0.0.1:8080/v1" {
		t.Errorf("original OpenAI URL changed: got %q", baseURL)
	}
	if status.Model == nil || (*status.Model)["name"] != "my-model" {
		t.Errorf("model name: got %v, want %q", status.Model, "my-model")
	}
}

func TestStatusStruct_NoActiveEngine(t *testing.T) {
	cache := storage.NewMockCache()
	ctx := &Context{
		Cache:  cache,
		Config: storage.NewMockConfig(),
		Snap:   snap.Mock(),
	}

	_, err := SnapStatus(ctx)
	if err == nil {
		t.Fatal("expected error for no active engine, got nil")
	}
	if !strings.Contains(err.Error(), "no active engine") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStatusBetterEngineNotice(t *testing.T) {
	t.Setenv("SNAP_INSTANCE_NAME", "test-snap")
	notice := "A better engine is available for your machine. Run test-snap use-engine --auto to switch."
	tests := []struct {
		name         string
		machine      string
		active       string
		engines      []string
		experimental string
		equalCPU     bool
		wantNotice   bool
		wantErr      string
	}{
		{name: "better engine", machine: "i7-1165G7", active: "cpu", engines: []string{"cpu", "intel-gpu"}, wantNotice: true},
		{name: "active is best", machine: "i7-1165G7", active: "intel-gpu", engines: []string{"cpu", "intel-gpu"}},
		{name: "equal scores", machine: "i7-1165G7", active: "cpu", engines: []string{"cpu"}, equalCPU: true},
		{name: "experimental not selectable", machine: "i7-1165G7", active: "cpu", engines: []string{"cpu", "intel-gpu"}, experimental: "intel-gpu"},
		{name: "no compatible engines", machine: "dummy-machine", active: "cpu-avx512", engines: []string{"cpu-avx512", "cuda-generic"}},
		{name: "active incompatible", machine: "dummy-machine", active: "cpu-avx512", engines: []string{"cpu-avx512", "cpu"}, wantNotice: true},
		{name: "active missing", machine: "i7-1165G7", active: "missing", engines: []string{"cpu"}, wantErr: `active engine "missing" not found`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			machine, err := testfixtures.Machine(tt.machine)
			if err != nil {
				t.Fatal(err)
			}
			var manifests []engines.Manifest
			for _, name := range tt.engines {
				manifest, err := engines.LoadManifest(filepath.Join("..", "..", "test_data", "engines"), name)
				if err != nil {
					t.Fatal(err)
				}
				if name == tt.experimental {
					experimental := true
					manifest.Experimental = &experimental
				}
				manifests = append(manifests, *manifest)
			}
			if tt.equalCPU {
				equal := manifests[0]
				equal.Name = "equal-cpu"
				manifests = append([]engines.Manifest{equal}, manifests...)
			}
			scored, err := selector.ScoreEngines(machine, manifests)
			if err != nil {
				t.Fatal(err)
			}
			status := &Status{Engine: tt.active, Notices: []string{"existing notice"}}
			err = addBetterEngineNotice(status, scored)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			want := []string{"existing notice"}
			if tt.wantNotice {
				want = append(want, notice)
			}
			if !reflect.DeepEqual(status.Notices, want) {
				t.Fatalf("notices = %v, want %v", status.Notices, want)
			}
			for _, format := range []string{"json", "yaml"} {
				var data []byte
				if format == "json" {
					data, err = json.Marshal(status)
				} else {
					data, err = yaml.Marshal(status)
				}
				if err != nil {
					t.Fatal(err)
				}
				var output struct {
					Notices []string `json:"notices" yaml:"notices"`
				}
				if format == "json" {
					err = json.Unmarshal(data, &output)
				} else {
					err = yaml.Unmarshal(data, &output)
				}
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(output.Notices, want) {
					t.Errorf("%s notices = %v, want %v", format, output.Notices, want)
				}
			}
		})
	}
}

func TestStatusHardwareObservability(t *testing.T) {
	tests := []struct {
		name       string
		observable bool
		observeErr error
		wantErr    string
	}{
		{name: "unavailable skips scoring"},
		{name: "observability error", observeErr: errors.New("observe failed"), wantErr: "observe failed"},
		{name: "scoring error", observable: true, wantErr: "scoring engines"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enginesDir, runtimesDir, modelsDir := t.TempDir(), t.TempDir(), t.TempDir()
			writeEngineYAML(t, enginesDir, "cpu", "name: cpu\nruntime: rt\n")
			writeEngineYAML(t, enginesDir, "bad-engine", "name: [\n")
			writeRuntimeYAML(t, runtimesDir, "rt", "name: rt\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n")
			writeModelYAML(t, modelsDir, "my-model", "name: my-model\n")
			cache := storage.NewMockCache()
			if err := cache.SetActiveEngine("cpu"); err != nil {
				t.Fatal(err)
			}
			if err := cache.SetActiveModel("my-model"); err != nil {
				t.Fatal(err)
			}
			config := storage.NewMockConfig()
			if err := config.Set("http.port", "8080", storage.UserConfig); err != nil {
				t.Fatal(err)
			}
			if err := config.Set("http.host", "127.0.0.1", storage.UserConfig); err != nil {
				t.Fatal(err)
			}
			ctx := &Context{
				EnginesDir: enginesDir, RuntimesDir: runtimesDir, ModelsDir: modelsDir,
				Cache: cache, Config: config,
				Snap: snap.MockWithHardwareObservability(tt.observable, tt.observeErr),
			}
			_, err := SnapStatus(ctx)
			if tt.wantErr == "" && err != nil || tt.wantErr != "" && (err == nil || !strings.Contains(err.Error(), tt.wantErr)) {
				t.Fatalf("got error %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func exampleStatusWithBetterEngine() (*Status, error) {
	machine, err := testfixtures.Machine("i7-1165G7")
	if err != nil {
		return nil, err
	}
	var manifests []engines.Manifest
	for _, name := range []string{"cpu", "intel-gpu"} {
		manifest, err := engines.LoadManifest(filepath.Join("..", "..", "test_data", "engines"), name)
		if err != nil {
			return nil, err
		}
		manifests = append(manifests, *manifest)
	}
	scored, err := selector.ScoreEngines(machine, manifests)
	if err != nil {
		return nil, err
	}
	model := map[string]string{"name": "26b-q4-k-m-gguf"}
	status := &Status{
		Engine:   "cpu",
		Services: map[string]string{"server": "active"},
		Model:    &model,
	}
	if err := addBetterEngineNotice(status, scored); err != nil {
		return nil, err
	}
	return status, nil
}

func setExampleSnapInstanceName() func() {
	previous, wasSet := os.LookupEnv("SNAP_INSTANCE_NAME")
	if err := os.Setenv("SNAP_INSTANCE_NAME", "modelctl-test"); err != nil {
		panic(err)
	}
	return func() {
		var err error
		if wasSet {
			err = os.Setenv("SNAP_INSTANCE_NAME", previous)
		} else {
			err = os.Unsetenv("SNAP_INSTANCE_NAME")
		}
		if err != nil {
			panic(err)
		}
	}
}

func ExampleStatus_yamlBetterEngine() {
	defer setExampleSnapInstanceName()()
	status, err := exampleStatusWithBetterEngine()
	if err != nil {
		panic(err)
	}
	output, err := yaml.Marshal(status)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(output))

	// Output:
	// engine: cpu
	// services:
	//     server: active
	// model:
	//     name: 26b-q4-k-m-gguf
	// notices:
	//     - A better engine is available for your machine. Run modelctl-test use-engine --auto to switch.
}

func ExampleStatus_jsonBetterEngine() {
	defer setExampleSnapInstanceName()()
	status, err := exampleStatusWithBetterEngine()
	if err != nil {
		panic(err)
	}
	output, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))

	// Output:
	// {
	//   "engine": "cpu",
	//   "services": {
	//     "server": "active"
	//   },
	//   "model": {
	//     "name": "26b-q4-k-m-gguf"
	//   },
	//   "notices": [
	//     "A better engine is available for your machine. Run modelctl-test use-engine --auto to switch."
	//   ]
	// }
}
