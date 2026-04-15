package webui

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func Serve(conf Config, staticDir string, port int, bindAddress string) error {

	if err := verifyStaticContent(staticDir); err != nil {
		return fmt.Errorf("unexpected static files: %w", err)
	}

	mux := serverMux(conf, staticDir)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", bindAddress, port), mux)
}

func serverMux(conf Config, staticDir string) *http.ServeMux {
	mux := http.NewServeMux()

	// Serve configuration for the frontend
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(conf); err != nil {
			http.Error(w, "encode config", http.StatusInternalServerError)
		}
	})

	// Serve the frontend static files
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))

	return mux
}

func verifyStaticContent(staticDir string) error {
	indexFile := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexFile); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("checking %q: %s", indexFile, err)
	}
	return nil
}
