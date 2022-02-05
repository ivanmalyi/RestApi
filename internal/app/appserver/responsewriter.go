package appserver

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	code int
}

func (w *ResponseWriter)writeHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
