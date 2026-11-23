package snap_store

import (
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

const (
	testSnapId   = "nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA"
	infoRoute    = "GET /v2/snaps/info/deepseek-r1"
	refreshRoute = "POST /v2/snaps/refresh"
)

const infoFixture = `{
	"name": "deepseek-r1",
	"snap": {
		"name": "deepseek-r1",
		"snap-id": "nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA",
		"title": "DeepSeek R1"
	},
	"snap-id": "nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA"
}`

const refreshFixture = `{
	"error-list": [],
	"results": [
		{
			"instance-key": "nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA",
			"name": "deepseek-r1",
			"result": "refresh",
			"snap": {
				"resources": [
					{"name": "llama-aio", "revision": 17, "type": "component/standard", "download": {"size": 4947968}},
					{"name": "llama-aio-ampereone", "revision": 10, "type": "component/standard", "download": {"size": 5001216}}
				]
			},
			"snap-id": "nr1Yeg25CSgcFmDpHD448ngXkZwSXPFA"
		}
	]
}`

func serveStore(t *testing.T, pattern string, status int, body string) {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	})
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	original := storeURL
	storeURL = server.URL
	t.Cleanup(func() { storeURL = original })
}

func TestSnapInfo(t *testing.T) {
	serveStore(t, infoRoute, http.StatusOK, infoFixture)

	info, err := snapInfo("deepseek-r1")
	if err != nil {
		t.Fatalf("error fetching info: %v", err)
	}
	if info.SnapId != testSnapId {
		t.Fatalf("snap id = %q, want %q", info.SnapId, testSnapId)
	}
}

func TestSnapInfoErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"non ok status", http.StatusServiceUnavailable, "", "http status not OK: 503"},
		{"malformed json", http.StatusOK, `{"snap-id":`, "decoding json"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			serveStore(t, infoRoute, tc.status, tc.body)

			_, err := snapInfo("deepseek-r1")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

func TestGetComponents(t *testing.T) {
	serveStore(t, refreshRoute, http.StatusOK, refreshFixture)

	components, err := snapComponents(testSnapId, 53, "amd64")
	if err != nil {
		t.Fatalf("error fetching components: %v", err)
	}

	got := make(map[string]int64)
	for _, component := range components {
		got[component.Name] = component.Download.Size
	}
	want := map[string]int64{"llama-aio": 4947968, "llama-aio-ampereone": 5001216}
	if len(components) != len(want) || !maps.Equal(got, want) {
		t.Fatalf("component sizes = %v, want %v", got, want)
	}
}

func TestGetComponentsErrors(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"non ok status", http.StatusServiceUnavailable, "", "http status not OK: 503"},
		{"malformed json", http.StatusOK, `{"results": [`, "decoding json"},
		{"no results", http.StatusOK, `{"error-list": [], "results": []}`, "store returned no refresh results"},
		{"other snap", http.StatusOK, `{"results": [{"snap-id": "other"}]}`, "no refresh results found for snap id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			serveStore(t, refreshRoute, tc.status, tc.body)

			_, err := snapComponents(testSnapId, 53, "amd64")
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want it to contain %q", err, tc.want)
			}
		})
	}
}

func TestComponentSizesCached(t *testing.T) {
	for _, tc := range []struct {
		name          string
		status        int
		refreshStatus int
		wantErr       bool
		wantRefresh   int32
	}{
		{name: "success", status: http.StatusOK, refreshStatus: http.StatusOK, wantRefresh: 1},
		{name: "failure", status: http.StatusServiceUnavailable, wantErr: true},
		{name: "refresh failure", status: http.StatusOK, refreshStatus: http.StatusServiceUnavailable, wantErr: true, wantRefresh: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("SNAP_NAME", "deepseek-r1")
			t.Setenv("SNAP_REVISION", "53")
			t.Setenv("SNAP_ARCH", "amd64")
			var infoRequests, refreshRequests atomic.Int32
			mux := http.NewServeMux()
			mux.HandleFunc(infoRoute, func(w http.ResponseWriter, r *http.Request) {
				infoRequests.Add(1)
				w.WriteHeader(tc.status)
				w.Write([]byte(infoFixture))
			})
			mux.HandleFunc(refreshRoute, func(w http.ResponseWriter, r *http.Request) {
				refreshRequests.Add(1)
				w.WriteHeader(tc.refreshStatus)
				w.Write([]byte(refreshFixture))
			})
			server := httptest.NewServer(mux)
			t.Cleanup(server.Close)
			originalURL, originalCache := storeURL, cachedComponentSizes
			storeURL = server.URL
			cachedComponentSizes = sync.OnceValues(fetchComponentSizes)
			t.Cleanup(func() {
				storeURL, cachedComponentSizes = originalURL, originalCache
			})

			var calls sync.WaitGroup
			for range 10 {
				calls.Go(func() {
					sizes, err := ComponentSizes()
					if (err != nil) != tc.wantErr {
						t.Errorf("ComponentSizes error = %v, wantErr %v", err, tc.wantErr)
					}
					if !tc.wantErr {
						want := map[string]int64{"llama-aio": 4947968, "llama-aio-ampereone": 5001216}
						if !maps.Equal(sizes, want) {
							t.Errorf("component sizes = %v, want %v", sizes, want)
						}
						sizes["llama-aio"] = 0
					}
				})
			}
			calls.Wait()
			if infoRequests.Load() != 1 || refreshRequests.Load() != tc.wantRefresh {
				t.Fatalf("requests: info=%d refresh=%d; want 1 and %d", infoRequests.Load(), refreshRequests.Load(), tc.wantRefresh)
			}
		})
	}
}
