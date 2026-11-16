package assetserver

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssetServerWithoutLogger(t *testing.T) {
	options := &Options{Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("response"))
	})}
	server, err := NewAssetServer(options)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	server.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/asset.txt", nil))
	if response.Code != http.StatusAccepted {
		t.Errorf("status = %d, want 202", response.Code)
	}
	server.ServeHTTP(&failingWriter{ResponseWriter: httptest.NewRecorder(), failOn: 1},
		httptest.NewRequest(http.MethodGet, "/asset.txt", nil))
	if options.Logger != nil {
		t.Error("constructor modified the caller's options")
	}
	if server.options.Logger.Enabled(context.Background(), slog.LevelError) {
		t.Error("missing logger enabled error logging")
	}
}

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

type requestLogContextKey struct{}

type requestContextLogHandler struct {
	slog.Handler
	contexts []context.Context
	levels   []slog.Level
}

func (h *requestContextLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return ctx.Value(requestLogContextKey{}) == "trace-42" && h.Handler.Enabled(ctx, level)
}

func (h *requestContextLogHandler) Handle(ctx context.Context, record slog.Record) error {
	h.contexts = append(h.contexts, ctx)
	h.levels = append(h.levels, record.Level)
	return h.Handler.Handle(ctx, record)
}

func TestAssetRequestLogsPreserveContext(t *testing.T) {
	for _, writeFailure := range []bool{false, true} {
		name := "debug"
		if writeFailure {
			name = "debug_and_write_error"
		}
		t.Run(name, func(t *testing.T) {
			logger := &requestContextLogHandler{Handler: slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})}
			server, err := NewAssetServer(&Options{Logger: slog.New(logger), Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if writeFailure {
					_, _ = io.WriteString(w, "response")
					return
				}
				w.WriteHeader(http.StatusAccepted)
			})})
			if err != nil {
				t.Fatal(err)
			}
			var response http.ResponseWriter = httptest.NewRecorder()
			if writeFailure {
				response = &failingWriter{ResponseWriter: response, failOn: 1}
			}
			request := httptest.NewRequest(http.MethodGet, "/asset.txt", nil)
			request = request.WithContext(context.WithValue(request.Context(), requestLogContextKey{}, "trace-42"))
			server.ServeHTTP(response, request)
			want := []slog.Level{slog.LevelDebug}
			if writeFailure {
				want = append(want, slog.LevelError)
			}
			if len(logger.contexts) != len(want) {
				t.Fatalf("logged levels=%v, want %v", logger.levels, want)
			}
			for i, ctx := range logger.contexts {
				if ctx.Value(requestLogContextKey{}) != "trace-42" || logger.levels[i] != want[i] {
					t.Errorf("record %d: context=%v, level=%v", i, ctx, logger.levels[i])
				}
			}
		})
	}
}

func TestAssetLogHelpersPreserveContext(t *testing.T) {
	handler := &requestContextLogHandler{Handler: slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})}
	ctx := context.WithValue(context.Background(), requestLogContextKey{}, "trace-42")
	ctx = contextWithLogger(ctx, slog.New(handler))
	logInfo(ctx, "asset info")
	logError(ctx, "asset error")
	want := []slog.Level{slog.LevelInfo, slog.LevelError}
	if len(handler.contexts) != len(want) {
		t.Fatalf("logged levels=%v, want %v", handler.levels, want)
	}
	for i, ctx := range handler.contexts {
		if ctx.Value(requestLogContextKey{}) != "trace-42" || handler.levels[i] != want[i] {
			t.Errorf("record %d lost context or level", i)
		}
	}
}
