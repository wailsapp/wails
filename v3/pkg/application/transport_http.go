package application

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"

	json "github.com/goccy/go-json"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

// bufferPool reduces allocations for reading request bodies
var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

type HTTPTransport struct {
	messageProcessor *MessageProcessor
	logger           *slog.Logger
}

func NewHTTPTransport(opts ...HTTPTransportOption) *HTTPTransport {
	t := &HTTPTransport{
		logger: slog.Default(),
	}

	// Apply options
	for _, opt := range opts {
		opt(t)
	}

	return t
}

// HTTPTransportOption is a functional option for configuring HTTPTransport
type HTTPTransportOption func(*HTTPTransport)

// HTTPTransportWithLogger is a functional option to set the logger for HTTPTransport.
func HTTPTransportWithLogger(logger *slog.Logger) HTTPTransportOption {
	return func(t *HTTPTransport) {
		t.logger = logger
	}
}

func (t *HTTPTransport) Start(ctx context.Context, processor *MessageProcessor) error {
	t.messageProcessor = processor

	return nil
}

func (t *HTTPTransport) JSClient() []byte {
	return nil
}

func (t *HTTPTransport) Stop() error {
	return nil
}

type request struct {
	Object *int            `json:"object"`
	Method *int            `json:"method"`
	Args   json.RawMessage `json:"args"`
}

func (t *HTTPTransport) Handler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			path := req.URL.Path
			switch path {
			case "/wails/runtime":
				t.handleRuntimeRequest(rw, req)
			default:
				next.ServeHTTP(rw, req)
			}
		})
	}
}

func (t *HTTPTransport) handleRuntimeRequest(rw http.ResponseWriter, r *http.Request) {
	// Use pooled buffer to reduce allocations
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	_, err := io.Copy(buf, r.Body)
	if err != nil {
		t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "Unable to read request body"))
		return
	}

	var body request
	err = json.Unmarshal(buf.Bytes(), &body)
	if err != nil {
		t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "Unable to parse request body as JSON"))
		return
	}

	if body.Object == nil {
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("missing object value"))
		return
	}

	if body.Method == nil {
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("missing method value"))
		return
	}

	windowIdStr := r.Header.Get(webViewRequestHeaderWindowId)
	windowId := 0
	if windowIdStr != "" {
		windowId, err = strconv.Atoi(windowIdStr)
		if err != nil {
			t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "error decoding windowId value"))
			return
		}
	}

	windowName := r.Header.Get(webViewRequestHeaderWindowName)
	clientId := r.Header.Get("x-wails-client-id")

	resp, err := t.messageProcessor.HandleRuntimeCallWithIDs(r.Context(), &RuntimeRequest{
		Object:            *body.Object,
		Method:            *body.Method,
		Args:              &Args{body.Args},
		WebviewWindowID:   uint32(windowId),
		WebviewWindowName: windowName,
		ClientID:          clientId,
	})

	if err != nil {
		t.httpError(rw, err)
		return
	}

	if stringResp, ok := resp.(string); ok {
		t.text(rw, stringResp)
		return
	}

	t.json(rw, resp)
}

func (t *HTTPTransport) text(rw http.ResponseWriter, data string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	_, err := rw.Write([]byte(data))
	if err != nil {
		t.error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
}

func (t *HTTPTransport) json(rw http.ResponseWriter, data any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	// convert data to json
	var jsonPayload = []byte("{}")
	var err error
	if data != nil {
		jsonPayload, err = json.Marshal(data)
		if err != nil {
			t.error("Unable to convert data to JSON. Please report this to the Wails team!", "error", err)
			return
		}
	}
	_, err = rw.Write(jsonPayload)
	if err != nil {
		t.error("Unable to write json payload. Please report this to the Wails team!", "error", err)
		return
	}
}

func (t *HTTPTransport) httpError(rw http.ResponseWriter, err error) {
	t.error(err.Error())
	// return JSON error if it's a CallError
	var bytes []byte
	if cerr := (*CallError)(nil); errors.As(err, &cerr) {
		if data, jsonErr := json.Marshal(cerr); jsonErr == nil {
			rw.Header().Set("Content-Type", "application/json")
			bytes = data
		} else {
			rw.Header().Set("Content-Type", "text/plain")
			bytes = []byte(err.Error())
		}
	} else {
		rw.Header().Set("Content-Type", "text/plain")
		bytes = []byte(err.Error())
	}
	rw.WriteHeader(http.StatusUnprocessableEntity)

	_, err = rw.Write(bytes)
	if err != nil {
		t.error("Unable to write error response:", "error", err)
	}
}

func (t *HTTPTransport) error(message string, args ...any) {
	t.logger.Error(message, args...)
}
