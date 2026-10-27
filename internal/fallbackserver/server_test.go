package fallbackserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerReportsServiceUnavailable(t *testing.T) {
	type apiError struct {
		Error struct {
			Message string  `json:"message"`
			Type    string  `json:"type"`
			Param   *string `json:"param"`
			Code    string  `json:"code"`
		} `json:"error"`
	}

	errorMessage := strings.Join([]string{
		"Not enough disk space.",
		`Model "example" could not be installed.`,
	}, "\n")
	responseHandler, err := handler(errorMessage)
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	for _, request := range []struct {
		method string
		path   string
	}{
		{method: http.MethodGet, path: "/"},
		{method: http.MethodGet, path: "/v1/models"},
		{method: http.MethodPost, path: "/v3/chat/completions"},
		{method: http.MethodGet, path: "/unknown"},
	} {
		t.Run(request.method+" "+request.path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			responseHandler.ServeHTTP(recorder, httptest.NewRequest(request.method, request.path, nil))

			if recorder.Code != http.StatusServiceUnavailable {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want %q", got, "application/json")
			}
			if got := recorder.Header().Get("Cache-Control"); got != "no-store" {
				t.Fatalf("Cache-Control = %q, want %q", got, "no-store")
			}

			var response apiError
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decoding response: %v", err)
			}
			if response.Error.Message != errorMessage {
				t.Fatalf("error message = %q, want %q", response.Error.Message, errorMessage)
			}
			if response.Error.Type != "unavailable_error" {
				t.Fatalf("error type = %q, want %q", response.Error.Type, "unavailable_error")
			}
			if response.Error.Param != nil {
				t.Fatalf("error param = %q, want nil", *response.Error.Param)
			}
			if response.Error.Code != "service_unavailable" {
				t.Fatalf("error code = %q, want %q", response.Error.Code, "service_unavailable")
			}
		})
	}
}

func TestHandlerWithoutErrorMessages(t *testing.T) {
	responseHandler, err := handler("")
	if err != nil {
		t.Fatalf("creating handler: %v", err)
	}

	recorder := httptest.NewRecorder()
	responseHandler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	var response errorResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if response.Error.Message != defaultErrorMessage {
		t.Fatalf("error message = %q, want %q", response.Error.Message, defaultErrorMessage)
	}
}
