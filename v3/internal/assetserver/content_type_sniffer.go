package assetserver

import (
	"net/http"
)

// newContentTypeSniffer creates a contentTypeSniffer that wraps the provided http.ResponseWriter.
// The returned sniffer does not allocate a close notification channel; it will be initialized lazily by CloseNotify.
func newContentTypeSniffer(rw http.ResponseWriter) *contentTypeSniffer {
	return &contentTypeSniffer{
		rw: rw,
	}
}

type contentTypeSniffer struct {
	rw              http.ResponseWriter
	prefix          []byte
	closeChannel    chan bool // lazily allocated only if CloseNotify is called
	status          int
	headerCommitted bool
	headerWritten   bool

	// err is sticky. complete may fail partway through emitting the sniffing
	// prefix, and http.Flusher has no way to report that, so the failure is
	// recorded here and surfaced from the next Write or complete instead of
	// being lost.
	err error
}

// Unwrap returns the wrapped [http.ResponseWriter] for use with [http.ResponseController].
func (rw *contentTypeSniffer) Unwrap() http.ResponseWriter {
	return rw.rw
}

func (rw *contentTypeSniffer) Header() http.Header {
	return rw.rw.Header()
}

func (rw *contentTypeSniffer) Write(chunk []byte) (int, error) {
	if rw.err != nil {
		return 0, rw.err
	}

	if !rw.headerCommitted {
		rw.WriteHeader(http.StatusOK)
	}

	if rw.headerWritten {
		return rw.rw.Write(chunk)
	}

	if len(chunk) == 0 {
		return 0, nil
	}

	// Cut away at most 512 bytes from chunk, and not less than 0.
	cut := max(min(len(chunk), 512-len(rw.prefix)), 0)
	if cut >= 512 {
		// Avoid copying data if a full prefix is available on first non-zero write.
		cut = len(chunk)
		rw.prefix = chunk
		chunk = nil
	} else if cut > 0 {
		// First write had less than 512 bytes -- copy data to the prefix buffer.
		if rw.prefix == nil {
			// Preallocate space for the prefix to be used for sniffing.
			rw.prefix = make([]byte, 0, 512)
		}
		rw.prefix = append(rw.prefix, chunk[:cut]...)
		chunk = chunk[cut:]
	}

	if len(rw.prefix) < 512 {
		return cut, nil
	}

	if _, err := rw.complete(); err != nil {
		return cut, err
	}

	n, err := rw.rw.Write(chunk)
	return cut + n, err
}

func (rw *contentTypeSniffer) WriteHeader(code int) {
	if rw.headerCommitted {
		return
	}

	rw.status = code
	rw.headerCommitted = true

	if _, hasType := rw.Header()[HeaderContentType]; hasType {
		rw.rw.WriteHeader(rw.status)
		rw.headerWritten = true
	}
}

// sniff sniffs the content type from the stored prefix if necessary,
// then writes the header.
func (rw *contentTypeSniffer) sniff() {
	if rw.headerWritten || !rw.headerCommitted {
		return
	}

	m := rw.Header()
	if _, hasType := m[HeaderContentType]; !hasType {
		m.Set(HeaderContentType, http.DetectContentType(rw.prefix))
	}

	rw.rw.WriteHeader(rw.status)
	rw.headerWritten = true
}

// complete sniffs the content type if necessary, writes the header
// and sends the data prefix that has been stored for sniffing.
//
// Whoever creates a contentTypeSniffer instance
// is responsible for calling complete after the nested handler has returned.
func (rw *contentTypeSniffer) complete() (n int, err error) {
	if rw.err != nil {
		return 0, rw.err
	}

	rw.sniff()

	if rw.headerWritten && len(rw.prefix) > 0 {
		n, err = rw.rw.Write(rw.prefix)

		// Drop only what actually went out. Clearing the whole prefix on a
		// short or failed write would discard bytes that were never sent, and
		// the caller would have no way to retry them.
		if n < 0 {
			n = 0
		}
		if n > len(rw.prefix) {
			n = len(rw.prefix)
		}
		rw.prefix = rw.prefix[n:]

		if err != nil {
			rw.err = err
		}
	}

	return
}

// CloseNotify implements the http.CloseNotifier interface.
// The channel is lazily allocated to avoid allocation overhead for requests
// that don't use this deprecated interface.
func (rw *contentTypeSniffer) CloseNotify() <-chan bool {
	if rw.closeChannel == nil {
		rw.closeChannel = make(chan bool, 1)
	}
	return rw.closeChannel
}

func (rw *contentTypeSniffer) closeClient() {
	if rw.closeChannel != nil {
		rw.closeChannel <- true
	}
}

// Flush implements the http.Flusher interface.
//
// The prefix has to be resolved first. Until 512 bytes have been seen the
// sniffer is deliberately holding the body back so it can detect a
// Content-Type, so delegating straight to the wrapped writer would flush
// nothing at all — the caller asks for a flush, gets no error, and no bytes
// reach the client. That silently breaks any response that emits less than
// 512 bytes and then waits, which is every streaming format.
//
// Completing here means the Content-Type is sniffed from a short prefix
// instead of a full one. That is the right trade: an explicit flush is the
// caller stating that what has been written so far should be sent now.
func (rw *contentTypeSniffer) Flush() {
	// Errors are dropped deliberately: http.Flusher cannot report one, and the
	// same write error will surface from the next Write or from complete.
	_, _ = rw.complete()

	if f, ok := rw.rw.(http.Flusher); ok {
		f.Flush()
	}
}
