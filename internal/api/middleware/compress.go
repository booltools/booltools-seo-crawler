package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} {
		writer, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return writer
	},
}

type compressedResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
}

func (crw *compressedResponseWriter) Write(data []byte) (int, error) {
	return crw.gzipWriter.Write(data)
}

func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}

		if strings.Contains(request.URL.Path, "/progress") {
			next.ServeHTTP(writer, request)
			return
		}

		gzWriter := gzipWriterPool.Get().(*gzip.Writer)
		gzWriter.Reset(writer)
		defer func() {
			gzWriter.Close()
			gzipWriterPool.Put(gzWriter)
		}()

		writer.Header().Set("Content-Encoding", "gzip")
		writer.Header().Set("Vary", "Accept-Encoding")
		writer.Header().Del("Content-Length")

		next.ServeHTTP(&compressedResponseWriter{ResponseWriter: writer, gzipWriter: gzWriter}, request)
	})
}
