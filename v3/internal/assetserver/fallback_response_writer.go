package assetserver

import (
	"maps"
	"net/http"
)

// fallbackResponseWriter wraps a [http.ResponseWriter].
// If the main handler returns status code 404,
// its response is discarded
// and the request is forwarded to the fallback handler.
type fallbackResponseWriter struct {
	rw       http.ResponseWriter
	req      *http.Request
	fallback http.Handler

	header        http.Header
	headerWritten bool
	complete      bool
}

// Unwrap returns the wrapped [http.ResponseWriter] for use with [http.ResponseController].
func (fw *fallbackResponseWriter) Unwrap() http.ResponseWriter {
	return fw.rw
}

func (fw *fallbackResponseWriter) Header() http.Header {
	if fw.header == nil {
		// Preserve original header in case we get a 404 response.
		fw.header = fw.rw.Header().Clone()
	}
	return fw.header
}

func (fw *fallbackResponseWriter) Write(chunk []byte) (int, error) {
	if fw.complete {
		// Fallback triggered, discard further writes.
		return len(chunk), nil
	}

	if !fw.headerWritten {
		fw.WriteHeader(http.StatusOK)
	}

	return fw.rw.Write(chunk)
}

func (fw *fallbackResponseWriter) WriteHeader(statusCode int) {
	if fw.headerWritten {
		return
	}
	fw.headerWritten = true

	if statusCode == http.StatusNotFound {
		// Protect fallback header from external modifications.
		if fw.header == nil {
			fw.header = fw.rw.Header().Clone()
		}

		// Invoke fallback handler.
		fw.complete = true
		fw.fallback.ServeHTTP(fw.rw, fw.req)
		return
	}

	if fw.header != nil {
		// Apply headers and forward original map to the main handler.
		maps.Copy(fw.rw.Header(), fw.header)
		fw.header = fw.rw.Header()
	}

	fw.rw.WriteHeader(statusCode)
}
