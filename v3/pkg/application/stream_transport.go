package application

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/assetserver/bundledassets"
)

// The desktop transport. Two endpoints on the asset server the app already
// serves — nothing listens on a port, and nothing goes through
// evaluateJavaScript.
//
//	GET  /wails/stream/poll   held until there is something to deliver
//	POST /wails/stream/send   one frame from the frontend
//
// Control data travels in headers, never in the body or the query string.
// WebKitGTK 6.0 can deliver POST bodies as query params for custom URI schemes
// (transport_http.go carries a fallback for exactly that), and WebView2 caps
// body delivery around 2 MB. Headers are delivered intact on all three engines.
const (
	streamPath     = "/wails/stream/"
	streamPathPoll = streamPath + "poll"
	streamPathSend = streamPath + "send"

	streamHeaderSession    = "x-wails-stream-session"
	streamHeaderGeneration = "x-wails-stream-generation"
	streamHeaderConn       = "x-wails-stream-conn"
	streamHeaderKind       = "x-wails-stream-kind"
	streamHeaderName       = "x-wails-stream-name"

	// Large frames are split by the client, for the same WebView2 body limit
	// the runtime's own chunking works around.
	streamHeaderChunkID    = "x-wails-stream-chunk"
	streamHeaderChunkIndex = "x-wails-stream-chunk-index"
	streamHeaderChunkTotal = "x-wails-stream-chunk-total"

	// Several data frames for one connection carried in a single POST. Each
	// request costs a scheme-handler round trip - about eleven cgo calls on
	// macOS - so batching divides that fixed cost by the batch size. Body
	// layout: count u32, then count x ( len u32 | payload ).
	streamHeaderBatch = "x-wails-stream-batch"

	streamMaxFrameBytes = 64 << 20
	streamMaxChunkTotal = 4096
	streamMaxChunkSets  = 256
	streamMaxChunkIDLen = 64
	streamMaxNameLen    = 256

	// Chunk reassembly is bounded both per session and across the application.
	// Completing a set temporarily retains both its parts and the contiguous
	// assembled frame, so the admitted payload allowance is half the effective
	// 256 MiB memory ceiling. This supports two simultaneous maximum-size uploads
	// without letting their reassembly copies double retained data past the cap.
	streamMaxChunkBytesGlobal = 2 * streamMaxFrameBytes

	// Payload bytes do not account for maps, slice headers, identifiers, and
	// per-part bookkeeping. Tiny or empty chunks could otherwise multiply that
	// metadata across every admitted session without approaching the byte cap.
	// The runtime emits at most 128 parts for a maximum-size frame, so this still
	// permits substantially more concurrent legitimate uploads than the byte
	// allowance can hold.
	streamMaxChunkPartsGlobal = 4096

	// Frame header on the wire: connID(4) + kind(1) + length(4).
	streamFrameHeaderBytes = 9
)

// streamMagic identifies the framing version. A page that has been served a
// stale cached runtime will fail the check loudly instead of misparsing.
var streamMagic = [4]byte{'W', 'S', '1', 0}

// runtimeJSWithPrelude returns the runtime bundle with the build's stream
// prelude in front of it, assembled once. Serving it this way is what makes the
// transport choice available to a stream created at module scope, which is the
// shape generated bindings use.
var runtimeJSWithPrelude = sync.OnceValue(func() []byte {
	prelude := streamPrelude()
	if len(prelude) == 0 {
		return bundledassets.RuntimeJS
	}
	out := make([]byte, 0, len(prelude)+len(bundledassets.RuntimeJS))
	out = append(out, prelude...)
	return append(out, bundledassets.RuntimeJS...)
})

// serveStream routes the two stream endpoints. Registered in the asset server
// middleware chain next to the event payload endpoint.
func (a *App) serveStream(rw http.ResponseWriter, req *http.Request) {
	if a.streams == nil {
		http.NotFound(rw, req)
		return
	}

	switch req.URL.Path {
	case streamPathPoll:
		a.serveStreamPoll(rw, req)
	case streamPathSend:
		a.serveStreamSend(rw, req)
	default:
		http.NotFound(rw, req)
	}
}

// streamRequestWindow reads the window id the platform layer tags onto every
// webview request (see webViewAssetRequest.Header). Zero means the platform did
// not tag it, which skips the binding check the way eventPayloadStore.take
// does.
func streamRequestWindow(req *http.Request) uint {
	raw := req.Header.Get(webViewRequestHeaderWindowId)
	if raw == "" {
		return 0
	}
	parsed, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil {
		return 0
	}
	return uint(parsed)
}

// sessionFor resolves the session named by the request and enforces that it
// belongs to the requesting window. Returns nil after writing a response.
// maySupersede is true only for the poll. A page bootstrapping always polls, so
// the poll is the honest signal that a new page has replaced the old one in this
// window. Letting a send supersede meant any stray or racing send carrying an
// unfamiliar session id tore down the live session and every connection on it —
// which is exactly what killed the harness partway through a run.
func (a *App) sessionFor(rw http.ResponseWriter, req *http.Request, create, maySupersede bool) *streamSession {
	id := req.Header.Get(streamHeaderSession)
	// Session ids are generated by the client, one per page load, exactly like
	// the runtime's clientId. They are not a capability: the binding below is
	// what stops one window reaching another's streams.
	if id == "" || len(id) > 64 {
		http.Error(rw, "missing session", http.StatusBadRequest)
		return nil
	}

	windowID := streamRequestWindow(req)
	rawGeneration := req.Header.Get(streamHeaderGeneration)
	if rawGeneration == "" {
		// A runtime from before page generations cannot participate safely in
		// reload superseding. Gone is deliberate: old poll clients understand it
		// and close their sockets, whereas a generic 400 makes them retry forever.
		http.Error(rw, "stream runtime is out of date", http.StatusGone)
		return nil
	}
	generation, err := strconv.ParseUint(rawGeneration, 10, 64)
	if err != nil || generation == 0 {
		http.Error(rw, "bad session generation", http.StatusBadRequest)
		return nil
	}

	var s *streamSession
	atCapacity := false
	if create {
		s, atCapacity = a.streams.sessionWithAdmission(id, windowID, generation, maySupersede)
	} else {
		s = a.streams.existingSession(id)
	}
	if s == nil {
		if atCapacity {
			// A poll from the current page can supersede older generations and
			// free this allowance. Make an open that raced ahead of that poll
			// retry instead of permanently closing the new socket.
			rw.Header().Set("Retry-After", "0")
			http.Error(rw, "stream session capacity reached", http.StatusTooManyRequests)
			return nil
		}
		// Gone rather than NotFound: it tells the client to stop polling and
		// surface a close, instead of retrying forever.
		http.Error(rw, "session closed", http.StatusGone)
		return nil
	}
	if windowID != 0 && s.windowID != windowID {
		http.Error(rw, "session belongs to another window", http.StatusForbidden)
		return nil
	}
	if s.generation != generation {
		http.Error(rw, "session generation mismatch", http.StatusForbidden)
		return nil
	}
	return s
}

func (a *App) serveStreamPoll(rw http.ResponseWriter, req *http.Request) {
	// GET only. A HEAD would run the same drain and then have its body
	// suppressed by net/http, so the frames would be consumed and never
	// delivered. This is a protocol endpoint, not a resource to probe.
	if req.Method != http.MethodGet {
		rw.Header().Set("Allow", "GET")
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	s := a.sessionFor(rw, req, true, true)
	if s == nil {
		return
	}

	frames, more, alive := s.awaitFrames(req.Context(), streamHoldTimeout, streamMaxResponseBytes)
	if !alive {
		a.streams.removeSession(s.id)
		http.Error(rw, "session closed", http.StatusGone)
		return
	}
	if len(frames) == 0 {
		rw.Header().Set("Cache-Control", "no-store")
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	bufp := streamRespBufs.Get().(*[]byte)
	body := encodeStreamFrames((*bufp)[:0], frames, more)

	rw.Header().Set("Content-Type", "application/octet-stream")
	rw.Header().Set("Cache-Control", "no-store")
	rw.Header().Set("Content-Length", strconv.Itoa(len(body)))
	// Write copies into the platform response before returning, so the buffer
	// is free to reuse immediately afterwards.
	_, _ = rw.Write(body)

	if cap(body) <= streamMaxResponseBytes {
		*bufp = body
		streamRespBufs.Put(bufp)
	}
}

// encodeStreamFrames lays out a poll response:
//
//	magic(4) | flags(1) | count(4) | count × ( connID(4) | kind(1) | len(4) | payload )
//
// Binary rather than JSON because frames are []byte: base64 inside a JSON
// envelope would cost 33% on every frame, and a megabyte of JSON also costs a
// parse on the UI thread.
// encodeStreamFrames appends the response into dst and returns it. Callers pass
// a pooled buffer: a response is up to streamMaxResponseBytes and there are
// thousands per second under load, so allocating one per poll was among the
// largest sources of garbage in the transport.
func encodeStreamFrames(dst []byte, frames []outFrame, more bool) []byte {
	size := len(streamMagic) + 1 + 4
	for _, f := range frames {
		size += streamFrameHeaderBytes + len(f.data)
	}
	if cap(dst) < size {
		dst = make([]byte, 0, size)
	}
	buf := dst[:0]

	buf = append(buf, streamMagic[:]...)
	var flags byte
	if more {
		flags |= 1
	}
	buf = append(buf, flags)
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(frames)))

	for _, f := range frames {
		buf = binary.BigEndian.AppendUint32(buf, f.connID)
		buf = append(buf, f.kind)
		buf = binary.BigEndian.AppendUint32(buf, uint32(len(f.data)))
		buf = append(buf, f.data...)
	}
	return buf
}

// streamRespBufs recycles poll response buffers. Buffers that grew past the
// response cap (a single oversized frame) are not returned, so one huge frame
// cannot pin megabytes in the pool for the life of the process.
var streamRespBufs = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 64<<10)
		return &b
	},
}

// decodeStreamBatch splits a batched send body into its frames. The payloads
// alias body rather than being copied: body is freshly read per request and is
// not retained anywhere else.
func decodeStreamBatch(body []byte) ([][]byte, error) {
	if len(body) < 4 {
		return nil, errStreamBadBody
	}
	count := binary.BigEndian.Uint32(body[:4])
	if count == 0 || count > streamMaxChunkTotal {
		return nil, errStreamBadBody
	}

	frames := make([][]byte, 0, count)
	off := 4
	for i := uint32(0); i < count; i++ {
		if len(body)-off < 4 {
			return nil, errStreamBadBody
		}
		n := int(binary.BigEndian.Uint32(body[off : off+4]))
		off += 4
		// n < 0 is unreachable where int is 64 bits and looks like dead code
		// there, but it is the guard that holds on a 32-bit build: a wire length
		// above math.MaxInt32 wraps negative, and body[off:off+n] would then
		// panic rather than be rejected. Keep it.
		//
		// The upper bound is written as n > len(body)-off rather than
		// off+n > len(body) so it cannot overflow before it is compared.
		if n < 0 || n > len(body)-off {
			return nil, errStreamBadBody
		}
		frames = append(frames, body[off:off+n])
		off += n
	}
	if off != len(body) {
		return nil, errStreamBadBody
	}
	return frames, nil
}

func (a *App) serveStreamSend(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		rw.Header().Set("Allow", "POST")
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	connID64, err := strconv.ParseUint(req.Header.Get(streamHeaderConn), 10, 32)
	if err != nil {
		http.Error(rw, "bad connection id", http.StatusBadRequest)
		return
	}
	connID := uint32(connID64)

	kind64, err := strconv.ParseUint(req.Header.Get(streamHeaderKind), 10, 8)
	if err != nil {
		http.Error(rw, "bad frame kind", http.StatusBadRequest)
		return
	}
	kind := uint8(kind64)
	switch kind {
	case frameOpen, frameData, frameClose:
	default:
		http.Error(rw, "unknown frame kind", http.StatusBadRequest)
		return
	}

	isBatch := req.Header.Get(streamHeaderBatch) != ""
	hasChunks := req.Header.Get(streamHeaderChunkID) != ""
	if isBatch && (kind != frameData || hasChunks) {
		// Runtime batches are already bounded below the platform body limit and
		// therefore are never chunked. Combining both retry protocols would make
		// a partial-batch acknowledgement ambiguous on a final-chunk retry.
		http.Error(rw, "invalid batched frame", http.StatusBadRequest)
		return
	}
	name := req.Header.Get(streamHeaderName)
	if kind == frameOpen && (name == "" || len(name) > streamMaxNameLen) {
		http.Error(rw, "invalid stream name", http.StatusBadRequest)
		return
	}
	if kind != frameData && (hasChunks || req.ContentLength != 0) {
		// Open and close carry all of their protocol data in headers. Reject a
		// body before session lookup so control requests cannot retain chunks that
		// no later dispatch path consumes.
		http.Error(rw, "control frame must not carry a body", http.StatusBadRequest)
		return
	}

	// Only an open frame may bring a session into being. A data or close frame
	// for a session that no longer exists used to silently create a fresh one:
	// Go then held a live session with no connections while the page believed
	// its streams were open, and every later send succeeded into the void. The
	// client is told the session is gone instead, so it surfaces a close.
	s := a.sessionFor(rw, req, kind == frameOpen, false)
	if s == nil {
		return
	}

	var dataConn *StreamConn
	if kind == frameData {
		dataConn = s.conn(connID)
		if dataConn == nil {
			http.Error(rw, "unknown connection", http.StatusGone)
			return
		}
	}

	body, complete, err := readStreamBody(s, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrStreamFull):
			rw.Header().Set("Retry-After", "0")
			http.Error(rw, "chunk reassembly is at capacity", http.StatusTooManyRequests)
		case errors.Is(err, ErrStreamClosed):
			http.Error(rw, "session closed", http.StatusGone)
		default:
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if !complete {
		// A chunk landed; the frame is not whole yet.
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	// A completed chunk set must survive a 429 from the connection inbox. The
	// client retries only the chunk whose request was rejected; deleting the set
	// before downstream acceptance would turn that retry into a new, incomplete
	// upload that answers 204 and silently loses the frame.
	chunkID := req.Header.Get(streamHeaderChunkID)
	removeChunk := chunkID != ""
	defer func() {
		if removeChunk {
			s.chunks().remove(chunkID)
		}
	}()

	// A batch is several data frames for one connection in one body. Splitting
	// here keeps the rest of the path identical to a single frame.
	if isBatch {
		frames, err := decodeStreamBatch(body)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		for i, f := range frames {
			switch err := dataConn.deliver(f); {
			case errors.Is(err, ErrStreamFull):
				removeChunk = false
				// Partial acceptance: tell the client how many landed so it can
				// resend only the remainder, preserving order.
				rw.Header().Set(streamHeaderBatch, strconv.Itoa(i))
				rw.Header().Set("Retry-After", "0")
				http.Error(rw, "receiver is behind", http.StatusTooManyRequests)
				return
			case err != nil:
				http.Error(rw, "connection closed", http.StatusGone)
				return
			}
		}
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	switch kind {
	case frameOpen:
		switch err := s.open(connID, name); {
		case errors.Is(err, ErrStreamFull):
			rw.Header().Set("Retry-After", "0")
			http.Error(rw, "stream session is at capacity", http.StatusTooManyRequests)
			return
		case errors.Is(err, ErrStreamClosed):
			http.Error(rw, "session closed", http.StatusGone)
			return
		case errors.Is(err, errStreamDuplicateConnection):
			http.Error(rw, "duplicate connection id", http.StatusConflict)
			return
		case err != nil:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	case frameData:
		// Queue before responding. The client does not issue its next send
		// until this response lands, so append-then-respond is what makes
		// frontend send order the order the handler observes.
		switch err := dataConn.deliver(body); {
		case errors.Is(err, ErrStreamFull):
			removeChunk = false
			// Backpressure, signalled rather than held. Waiting here would
			// occupy a request slot and starve the window's poll; the client
			// retries this same frame instead, so order is preserved and the
			// inbox stays bounded.
			rw.Header().Set("Retry-After", "0")
			http.Error(rw, "receiver is behind", http.StatusTooManyRequests)
			return
		case err != nil:
			http.Error(rw, "connection closed", http.StatusGone)
			return
		}

	case frameClose:
		if c := s.conn(connID); c != nil {
			c.closedByPeer()
		}

	}

	rw.WriteHeader(http.StatusNoContent)
}

// readStreamFrameBody reads the whole body in one allocation.
//
// io.ReadAll starts at 512 bytes and doubles, so a 512 KB frame costs about
// eleven allocations and a megabyte of copying — on the hottest path there is,
// once per frame. The webview gives us Content-Length, so size the buffer
// exactly and read straight into it.
func readStreamFrameBody(req *http.Request) ([]byte, error) {
	if n := req.ContentLength; n > 0 {
		if n > streamMaxFrameBytes {
			return nil, errStreamBadBody
		}
		buf := make([]byte, int(n))
		if _, err := io.ReadFull(req.Body, buf); err != nil {
			return nil, errStreamBadBody
		}
		return buf, nil
	}
	// No Content-Length (some platforms omit it): fall back to growing.
	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errStreamBadBody
	}
	return data, nil
}

// readStreamBody returns the frame body, reassembling it first if the client
// split it. complete is false when this request carried a non-final chunk.
func readStreamBody(s *streamSession, req *http.Request) (body []byte, complete bool, err error) {
	req.Body = http.MaxBytesReader(nil, req.Body, streamMaxFrameBytes)

	chunkID := req.Header.Get(streamHeaderChunkID)
	if chunkID == "" {
		data, err := readStreamFrameBody(req)
		if err != nil {
			return nil, false, err
		}
		return data, true, nil
	}
	if len(chunkID) > streamMaxChunkIDLen {
		return nil, false, errStreamBadChunk
	}

	total, err := strconv.Atoi(req.Header.Get(streamHeaderChunkTotal))
	if err != nil || total <= 0 || total > streamMaxChunkTotal {
		return nil, false, errStreamBadChunk
	}
	index, err := strconv.Atoi(req.Header.Get(streamHeaderChunkIndex))
	if err != nil || index < 0 || index >= total {
		return nil, false, errStreamBadChunk
	}

	data, err := readStreamFrameBody(req)
	if err != nil {
		return nil, false, err
	}

	assembled, done, err := s.chunks().add(chunkID, index, total, data)
	if err != nil {
		return nil, false, err
	}
	if !done {
		return nil, false, nil
	}
	return assembled, true, nil
}

type streamError string

func (e streamError) Error() string { return string(e) }

const (
	errStreamBadBody       streamError = "unable to read frame body"
	errStreamBadChunk      streamError = "invalid chunk headers"
	errStreamChunkTooBig   streamError = "assembled frame too large"
	errStreamChunkConflict streamError = "conflicting chunk total"
	errStreamChunkDup      streamError = "duplicate chunk index"
)

// ---------------------------------------------------------------------------
// Chunk reassembly for oversized frontend frames
// ---------------------------------------------------------------------------

type streamChunkStore struct {
	mu     sync.Mutex
	items  map[string]*streamChunkSet
	bytes  int
	parts  int
	closed bool
	mgr    *streamManager
}

type streamChunkSet struct {
	parts      map[int][]byte
	total      int
	size       int
	partCount  int
	created    time.Time
	assembled  []byte
	retryIndex int
	retryPart  []byte
}

func (s *streamSession) chunks() *streamChunkStore {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.chunkStore == nil {
		s.chunkStore = &streamChunkStore{
			items:  make(map[string]*streamChunkSet),
			closed: s.closed,
			mgr:    s.mgr,
		}
	}
	return s.chunkStore
}

// add stores one chunk and returns the assembled frame once the last one lands.
//
// A rejection is reported as an error rather than as "not done yet". Returning
// the same not-done signal for both made the endpoint answer 204, so the
// client's send chain completed successfully while the frame was silently
// dropped.
func (c *streamChunkStore) add(id string, index, total int, data []byte) ([]byte, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, false, ErrStreamClosed
	}

	// Reap anything a disappearing page left behind before taking more.
	now := time.Now()
	for k, v := range c.items {
		if now.Sub(v.created) > streamHoldTimeout {
			c.bytes -= v.size
			c.parts -= v.partCount
			c.releaseResources(v.size, v.partCount)
			delete(c.items, k)
		}
	}

	set, ok := c.items[id]
	createdSet := !ok
	if !ok {
		if len(c.items) >= streamMaxChunkSets {
			return nil, false, ErrStreamFull
		}
		set = &streamChunkSet{parts: make(map[int][]byte, total), total: total, created: now}
		c.items[id] = set
	}
	if set.total != total {
		// Discard the set, and give back what it was holding — not doing so
		// consumed that capacity for the life of the session.
		c.bytes -= set.size
		c.parts -= set.partCount
		c.releaseResources(set.size, set.partCount)
		delete(c.items, id)
		return nil, false, errStreamChunkConflict
	}
	if set.assembled != nil {
		// Only the request that completed this set can legitimately be retried:
		// it is the one that may have received 429 after assembly. Return the
		// same immutable frame until the endpoint acknowledges downstream
		// delivery and removes the set.
		if index == set.retryIndex && bytes.Equal(data, set.retryPart) {
			return set.assembled, true, nil
		}
		c.bytes -= set.size
		c.parts -= set.partCount
		c.releaseResources(set.size, set.partCount)
		delete(c.items, id)
		return nil, false, errStreamChunkDup
	}
	if _, dup := set.parts[index]; dup {
		c.bytes -= set.size
		c.parts -= set.partCount
		c.releaseResources(set.size, set.partCount)
		delete(c.items, id)
		return nil, false, errStreamChunkDup
	}
	if c.bytes+len(data) > streamMaxFrameBytes {
		c.bytes -= set.size
		c.parts -= set.partCount
		c.releaseResources(set.size, set.partCount)
		delete(c.items, id)
		return nil, false, errStreamChunkTooBig
	}
	if c.mgr != nil && !c.mgr.reserveChunkResources(len(data), 1) {
		// The shared allowance is transient backpressure. Keep an existing
		// partial set so the client can retry this exact chunk once another
		// session releases capacity; do not retain a newly-created empty set.
		if createdSet {
			delete(c.items, id)
		}
		return nil, false, ErrStreamFull
	}

	set.parts[index] = data
	set.size += len(data)
	set.partCount++
	c.bytes += len(data)
	c.parts++

	if len(set.parts) < total {
		return nil, false, nil
	}

	assembled := make([]byte, 0, set.size)
	retryStart := 0
	retryEnd := 0
	for i := 0; i < total; i++ {
		start := len(assembled)
		assembled = append(assembled, set.parts[i]...)
		if i == index {
			retryStart = start
			retryEnd = len(assembled)
		}
	}
	set.parts = nil
	set.assembled = assembled
	set.retryIndex = index
	set.retryPart = assembled[retryStart:retryEnd]
	return assembled, true, nil
}

// close releases every incomplete or retryable upload when its page session
// ends and prevents an already-dispatched request from recreating state after
// cancellation.
func (c *streamChunkStore) close() {
	c.mu.Lock()
	c.closed = true
	c.items = make(map[string]*streamChunkSet)
	released := c.bytes
	releasedParts := c.parts
	c.bytes = 0
	c.parts = 0
	c.mu.Unlock()
	c.releaseResources(released, releasedParts)
}

// remove acknowledges a completed set after its frame has been accepted by
// the connection inbox (or permanently rejected for some reason other than
// backpressure).
func (c *streamChunkStore) remove(id string) {
	c.mu.Lock()
	var released int
	var releasedParts int
	if set, ok := c.items[id]; ok {
		c.bytes -= set.size
		c.parts -= set.partCount
		released = set.size
		releasedParts = set.partCount
		delete(c.items, id)
	}
	c.mu.Unlock()
	c.releaseResources(released, releasedParts)
}

func (c *streamChunkStore) releaseResources(bytes, parts int) {
	if c.mgr != nil {
		c.mgr.releaseChunkResources(bytes, parts)
	}
}
