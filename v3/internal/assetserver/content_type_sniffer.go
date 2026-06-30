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
}

// Unwrap returns the wrapped [http.ResponseWriter] for use with [http.ResponseController].
func (rw *contentTypeSniffer) Unwrap() http.ResponseWriter {
	return rw.rw
}

func (rw *contentTypeSniffer) Header() http.Header {
	return rw.rw.Header()
}

func (rw *contentTypeSniffer) Write(chunk []byte) (int, error) {
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
	rw.sniff()

	if rw.headerWritten && len(rw.prefix) > 0 {
		n, err = rw.rw.Write(rw.prefix)
		rw.prefix = nil
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
func (rw *contentTypeSniffer) Flush() {
	if f, ok := rw.rw.(http.Flusher); ok {
		f.Flush()
	}
}