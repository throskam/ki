package ki

import (
	"bytes"
	"net/http"
)

// BufferedResponseWriter is a response writer that buffers the response.
type BufferedResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buffer     bytes.Buffer
	flushed    bool
}

// NewBufferedResponseWriter returns a new BufferedResponseWriter.
func NewBufferedResponseWriter(w http.ResponseWriter) *BufferedResponseWriter {
	return &BufferedResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

// Header returns the header of the response.
func (w *BufferedResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// WriteHeader sets the status code of the response.
func (w *BufferedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// Write writes the given bytes to the buffer if the response has not been flushed or the response has already otherwise.
func (w *BufferedResponseWriter) Write(b []byte) (int, error) {
	if w.flushed {
		return w.ResponseWriter.Write(b)
	}

	return w.buffer.Write(b)
}

// Flush writes the buffered response to the underlying response writer.
func (w *BufferedResponseWriter) Flush() (int, error) {
	if w.flushed {
		return 0, nil
	}

	w.ResponseWriter.WriteHeader(w.statusCode)

	w.flushed = true

	if w.buffer.Len() > 0 {
		return w.ResponseWriter.Write(w.buffer.Bytes())
	}

	return 0, nil
}

// StatusCode returns the status code of the response.
func (w *BufferedResponseWriter) StatusCode() int {
	return w.statusCode
}

// Size returns the size of the buffered response.
func (w *BufferedResponseWriter) Size() int {
	return w.buffer.Len()
}

