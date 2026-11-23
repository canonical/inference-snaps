package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/canonical/inference-snaps/modelctl/v2/commands/common"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/fallbackserver"
	"github.com/canonical/inference-snaps/modelctl/v2/internal/storage"
)

const diagnoseRuntime = "name: test-runtime\nservers:\n  openai:\n    protocol: http\n    base-path: /v1\n"

func TestDiagnose(t *testing.T) {
	tests := []struct {
		name       string
		noEngine   bool
		noModel    bool
		args       []string
		wantErr    error
		wantErrMsg string
	}{
		{name: "ready"},
		{name: "ready with fallback", args: []string{"--serve-fallback=openai"}},
		{name: "ready rejects unsupported protocol", args: []string{"--serve-fallback=kserve"}, wantErrMsg: `unsupported fallback protocol "kserve"`},
		{name: "ready with empty protocol", args: []string{"--serve-fallback="}},
		{name: "no engine", noEngine: true, wantErr: common.ErrNoActiveEngine},
		{name: "no model", noModel: true, wantErr: common.ErrNoActiveModel},
		{name: "neither active", noEngine: true, noModel: true, wantErr: common.ErrNoActiveEngine},
		{name: "no engine with unsupported protocol", noEngine: true, args: []string{"--serve-fallback=kserve"}, wantErrMsg: `unsupported fallback protocol "kserve"`},
		{name: "no model with unsupported protocol", noModel: true, args: []string{"--serve-fallback=kserve"}, wantErrMsg: `unsupported fallback protocol "kserve"`},
		{name: "no engine with empty protocol", noEngine: true, args: []string{"--serve-fallback="}, wantErr: common.ErrNoActiveEngine},
		{name: "no model with empty protocol", noModel: true, args: []string{"--serve-fallback="}, wantErr: common.ErrNoActiveModel},
		{name: "missing protocol", args: []string{"--serve-fallback"}, wantErrMsg: "flag needs an argument"},
		{name: "unexpected argument", args: []string{"extra"}, wantErrMsg: "unknown command"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := testRunContext(t, diagnoseRuntime)
			activateTestModel(t, ctx)
			if tt.noEngine {
				if err := ctx.Cache.SetActiveEngine(""); err != nil {
					t.Fatal(err)
				}
			}
			if tt.noModel {
				if err := ctx.Cache.SetActiveModel(""); err != nil {
					t.Fatal(err)
				}
			}
			cmd := Diagnose(ctx)
			cmd.SetArgs(append([]string{}, tt.args...))
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			err := cmd.Execute()
			switch {
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
			case tt.wantErrMsg != "":
				if err == nil || !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("error = %v, want containing %q", err, tt.wantErrMsg)
				}
			case err != nil:
				t.Fatal(err)
			}
		})
	}
}

func TestDiagnoseMissingManifest(t *testing.T) {
	for _, manifest := range []string{"engine", "model"} {
		t.Run(manifest, func(t *testing.T) {
			ctx := testRunContext(t, diagnoseRuntime)
			activateTestModel(t, ctx)
			path := filepath.Join(ctx.ModelsDir, "test-model", "model.yaml")
			if manifest == "engine" {
				path = filepath.Join(ctx.EnginesDir, "test-engine", "engine.yaml")
			}
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			cmd := Diagnose(ctx)
			cmd.SetArgs([]string{})
			cmd.SetErr(&bytes.Buffer{})
			if err := cmd.Execute(); err == nil {
				t.Fatal("expected missing manifest to fail diagnosis")
			}
		})
	}
}

func TestDiagnoseFallbackListenAddress(t *testing.T) {
	t.Run("ignores namespaced runtime endpoint", func(t *testing.T) {
		ctx := testRunContext(t, diagnoseRuntime+"    namespace: api\n")
		for key, value := range map[string]string{"api.http.host": "::1", "api.http.port": "9090"} {
			if err := ctx.Config.Set(key, value, storage.UserConfig); err != nil {
				t.Fatal(err)
			}
		}
		address, err := (&diagnoseCommand{Context: ctx}).fallbackListenAddress()
		if err != nil || address != "127.0.0.1:8080" {
			t.Fatalf("address = %q, error = %v", address, err)
		}
	})
	t.Run("no OpenAI endpoint", func(t *testing.T) {
		ctx := testRunContext(t, "name: test-runtime\nservers:\n  kserve:\n    protocol: http\n    base-path: /v2\n")
		address, err := (&diagnoseCommand{Context: ctx}).fallbackListenAddress()
		if err != nil || address != "127.0.0.1:8080" {
			t.Fatalf("address = %q, error = %v", address, err)
		}
	})
	t.Run("IPv6 configuration", func(t *testing.T) {
		ctx := testRunContext(t, diagnoseRuntime)
		if err := ctx.Config.Set("http.host", "::1", storage.UserConfig); err != nil {
			t.Fatal(err)
		}
		address, err := (&diagnoseCommand{Context: ctx}).fallbackListenAddress()
		if err != nil || address != "[::1]:8080" {
			t.Fatalf("address = %q, error = %v", address, err)
		}
	})
	for _, key := range []string{"http.host", "http.port"} {
		t.Run("missing "+key, func(t *testing.T) {
			ctx := testRunContext(t, diagnoseRuntime)
			if err := ctx.Config.Set(key, "", storage.UserConfig); err != nil {
				t.Fatal(err)
			}
			_, err := (&diagnoseCommand{Context: ctx}).fallbackListenAddress()
			if err == nil || !strings.Contains(err.Error(), key) {
				t.Fatalf("error = %v, want missing %s", err, key)
			}
		})
	}
}

func TestDiagnoseFallbackBindFailure(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	ctx := testRunContext(t, diagnoseRuntime)
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	if err := ctx.Config.Set("http.port", port, storage.UserConfig); err != nil {
		t.Fatal(err)
	}
	cmd := Diagnose(ctx)
	cmd.SetArgs([]string{"--serve-fallback=openai"})
	cmd.SetErr(&bytes.Buffer{})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "running fallback server") {
		t.Fatalf("error = %v, want fallback startup failure", err)
	}
}

func TestDiagnoseFallbackServer(t *testing.T) {
	for _, mode := range []string{"no-engine", "no-model", "neither"} {
		t.Run(mode, func(t *testing.T) {
			port := availableTCPPort(t)
			cmd := exec.Command(os.Args[0], "-test.run=^TestDiagnoseFallbackHelper$")
			cmd.Env = append(os.Environ(),
				"GO_WANT_DIAGNOSE_HELPER="+mode,
				"DIAGNOSE_HELPER_PORT="+strconv.Itoa(port),
			)
			output := &bytes.Buffer{}
			cmd.Stdout = output
			cmd.Stderr = output
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			exited := false
			t.Cleanup(func() {
				if !exited {
					_ = cmd.Process.Kill()
					<-done
				}
			})

			client := &http.Client{Timeout: 100 * time.Millisecond}
			url := fmt.Sprintf("http://127.0.0.1:%d/v1/chat/completions", port)
			deadline := time.Now().Add(5 * time.Second)
			for {
				response, err := client.Get(url)
				if err == nil {
					defer response.Body.Close()
					var body struct {
						Error struct {
							Message string `json:"message"`
							Code    string `json:"code"`
						} `json:"error"`
					}
					if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
						t.Fatal(err)
					}
					want := "No active engine is set. Please set an active engine and try again."
					if mode == "no-model" {
						want = "No active model is set. Please set an active model and try again."
					}
					if response.StatusCode != http.StatusServiceUnavailable || body.Error.Message != want || body.Error.Code != fallbackserver.ErrorCode {
						t.Fatalf("unexpected fallback response: status %d, body %+v", response.StatusCode, body)
					}
					chatErr := common.ChatClient(fmt.Sprintf("http://127.0.0.1:%d/v1", port), "", false).WaitChatServerReady()
					if wantErr := want + "\n\n" + common.SuggestServerLogs(); chatErr == nil || chatErr.Error() != wantErr {
						t.Fatalf("chat error = %v, want %q", chatErr, wantErr)
					}
					return
				}
				select {
				case err := <-done:
					exited = true
					t.Fatalf("fallback exited before serving: %v\n%s", err, output.String())
				default:
				}
				if time.Now().After(deadline) {
					t.Fatalf("fallback server did not start: %v", err)
				}
				time.Sleep(25 * time.Millisecond)
			}
		})
	}
}

func TestDiagnoseFallbackHelper(t *testing.T) {
	mode := os.Getenv("GO_WANT_DIAGNOSE_HELPER")
	if mode == "" {
		return
	}
	ctx := testRunContext(t, diagnoseRuntime)
	if err := ctx.Config.Set("http.port", os.Getenv("DIAGNOSE_HELPER_PORT"), storage.UserConfig); err != nil {
		t.Fatal(err)
	}
	if mode != "no-model" {
		if err := ctx.Cache.SetActiveEngine(""); err != nil {
			t.Fatal(err)
		}
	}
	if mode == "no-engine" {
		activateTestModel(t, ctx)
	}
	cmd := Diagnose(ctx)
	cmd.SetArgs([]string{"--serve-fallback=openai"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func availableTCPPort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserving TCP port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("releasing TCP port: %v", err)
	}
	return port
}
