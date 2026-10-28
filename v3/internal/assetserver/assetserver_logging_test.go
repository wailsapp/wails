package assetserver

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssetServerRequestDebugLogging(t *testing.T) {
	var output bytes.Buffer
	server, err := NewAssetServer(&Options{
		Logger: slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusAccepted)
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/asset%20name.txt", nil)
	request.Header.Set(webViewRequestHeaderWindowId, "42")
	request.Header.Set(webViewRequestHeaderWindowName, "main")
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	if response.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", response.Code)
	}
	var record map[string]any
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]any{
		"msg": "Asset Request:", "windowName": "main", "windowID": "42",
		"code": float64(http.StatusAccepted), "method": http.MethodGet, "path": "/asset%20name.txt",
	} {
		if record[key] != want {
			t.Errorf("log %s = %v, want %v", key, record[key], want)
		}
	}
	if _, ok := record["duration"].(float64); !ok {
		t.Error("debug log omitted request duration")
	}
}
