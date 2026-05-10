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
	"time"

	"encoding/json"

	"github.com/wailsapp/wails/v3/pkg/errs"
)

// bufferPool reduces allocations for reading request bodies.
// Buffers larger than maxPooledBufferSize are not returned to the pool
// to prevent memory bloat from occasional large requests (e.g., images).
const maxPooledBufferSize = 512 * 1024 // 512KB

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(make([]byte, 0, 4096))
	},
}

const (
	chunkIDHeader    = "x-wails-chunk-id"
	chunkIndexHeader = "x-wails-chunk-index"
	chunkTotalHeader = "x-wails-chunk-total"
	chunkTTL         = 30 * time.Second

	maxChunkTotal     = 1024
	maxChunkBodyBytes = 1024 * 1024
	maxAssembledBytes = 64 * 1024 * 1024
)

// pendingChunks accumulates request body chunks sent by the JS runtime
// to work around WebView2's ~2MB limit on request body content delivery
// via the WebResourceRequested event.
type pendingChunks struct {
	mu        sync.Mutex
	chunks    map[int][]byte
	total     int
	size      int
	createdAt time.Time
}

type HTTPTransport struct {
	messageProcessor *MessageProcessor
	logger           *slog.Logger
	chunkStore       sync.Map
	stopCleanup      chan struct{}
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
	t.stopCleanup = make(chan struct{})
	go t.cleanupChunks()
	return nil
}

func (t *HTTPTransport) cleanupChunks() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-t.stopCleanup:
			return
		case <-ticker.C:
			now := time.Now()
			t.chunkStore.Range(func(k, v any) bool {
				pc := v.(*pendingChunks)
				pc.mu.Lock()
				expired := now.Sub(pc.createdAt) > chunkTTL
				pc.mu.Unlock()
				if expired {
					t.chunkStore.Delete(k)
				}
				return true
			})
		}
	}
}

func (t *HTTPTransport) JSClient() []byte {
	return nil
}

func (t *HTTPTransport) Stop() error {
	if t.stopCleanup != nil {
		close(t.stopCleanup)
		t.stopCleanup = nil
	}
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
	// Chunked upload: JS splits large bodies into smaller pieces to work
	// around WebView2's ~2MB request body delivery limit in WebResourceRequested.
	if chunkID := r.Header.Get(chunkIDHeader); chunkID != "" {
		t.handleChunkedRequest(rw, r, chunkID)
		return
	}

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		if buf.Cap() <= maxPooledBufferSize {
			bufferPool.Put(buf)
		}
	}()

	if _, err := io.Copy(buf, r.Body); err != nil {
		t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "Unable to read request body"))
		return
	}

	t.processBody(rw, r, buf.Bytes())
}

// handleChunkedRequest stores an incoming chunk and, once all chunks are
// received, assembles them and delegates to processBody.
func (t *HTTPTransport) handleChunkedRequest(rw http.ResponseWriter, r *http.Request, chunkID string) {
	indexStr := r.Header.Get(chunkIndexHeader)
	totalStr := r.Header.Get(chunkTotalHeader)

	total, err := strconv.Atoi(totalStr)
	if err != nil || total <= 0 || total > maxChunkTotal {
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("invalid chunk total: %s", totalStr))
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= total {
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("invalid chunk index: %s", indexStr))
		return
	}

	r.Body = http.MaxBytesReader(rw, r.Body, int64(maxChunkBodyBytes))

	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		if buf.Cap() <= maxPooledBufferSize {
			bufferPool.Put(buf)
		}
	}()

	if _, err := io.Copy(buf, r.Body); err != nil {
		t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "unable to read chunk body"))
		return
	}

	chunk := make([]byte, buf.Len())
	copy(chunk, buf.Bytes())

	actual, _ := t.chunkStore.LoadOrStore(chunkID, &pendingChunks{
		chunks:    make(map[int][]byte),
		total:     total,
		createdAt: time.Now(),
	})
	pc := actual.(*pendingChunks)

	pc.mu.Lock()
	if pc.total != total {
		pc.mu.Unlock()
		t.chunkStore.Delete(chunkID)
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("inconsistent chunk total"))
		return
	}
	pc.chunks[index] = chunk
	pc.size += len(chunk)
	if pc.size > maxAssembledBytes {
		pc.mu.Unlock()
		t.chunkStore.Delete(chunkID)
		t.httpError(rw, errs.NewInvalidRuntimeCallErrorf("assembled body too large"))
		return
	}
	received := len(pc.chunks)
	pc.mu.Unlock()

	if received < total {
		rw.WriteHeader(http.StatusOK)
		return
	}

	// All chunks received — assemble in order and process.
	t.chunkStore.Delete(chunkID)

	pc.mu.Lock()
	var assembled []byte
	for i := 0; i < pc.total; i++ {
		assembled = append(assembled, pc.chunks[i]...)
	}
	pc.mu.Unlock()

	t.processBody(rw, r, assembled)
}

func (t *HTTPTransport) processBody(rw http.ResponseWriter, r *http.Request, bodyBytes []byte) {
	var body request
	var err error

	if len(bodyBytes) > 0 {
		if err = json.Unmarshal(bodyBytes, &body); err != nil {
			t.httpError(rw, errs.WrapInvalidRuntimeCallErrorf(err, "Unable to parse request body as JSON"))
			return
		}
	} else {
		// Fallback: WebKitGTK 6.0 may send POST data as query params for custom URI schemes
		query := r.URL.Query()
		if objStr := query.Get("object"); objStr != "" {
			obj, parseErr := strconv.Atoi(objStr)
			if parseErr == nil {
				body.Object = &obj
			}
		}
		if methStr := query.Get("method"); methStr != "" {
			meth, parseErr := strconv.Atoi(methStr)
			if parseErr == nil {
				body.Method = &meth
			}
		}
		if argsStr := query.Get("args"); argsStr != "" {
			var args json.RawMessage
			if json.Unmarshal([]byte(argsStr), &args) == nil {
				body.Args = args
			}
		}
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
