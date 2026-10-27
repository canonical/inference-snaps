package fallbackserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const defaultErrorMessage = "Inference server failed to start."

type errorResponse struct {
	Error struct {
		Message string  `json:"message"`
		Type    string  `json:"type"`
		Param   *string `json:"param"`
		Code    string  `json:"code"`
	} `json:"error"`
}

// Run starts an OpenAI-compatible error server at address and blocks until the
// server stops.
func Run(address string, errorMessage string) error {
	responseHandler, err := handler(errorMessage)
	if err != nil {
		return fmt.Errorf("creating error response: %w", err)
	}

	server := &http.Server{
		Addr:              address,
		Handler:           responseHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	fmt.Printf("[fallback] Listening on %s\n", address)
	return server.ListenAndServe()
}

func handler(errorMessage string) (http.Handler, error) {
	message := defaultErrorMessage
	if strings.TrimSpace(errorMessage) != "" {
		message = errorMessage
	}

	response := errorResponse{}
	response.Error.Message = message
	response.Error.Type = "unavailable_error"
	response.Error.Code = "service_unavailable"

	body, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("encoding OpenAI error: %w", err)
	}
	body = append(body, '\n')

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientAddress := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientAddress = host
		}
		status := http.StatusServiceUnavailable
		if r.Method == http.MethodOptions {
			status = http.StatusNoContent
		}
		fmt.Printf("[fallback] %s - \"%s %s %s\" %d -\n",
			clientAddress, r.Method, r.URL.RequestURI(), r.Proto, status)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(status)
		_, _ = w.Write(body)
	}), nil
}
