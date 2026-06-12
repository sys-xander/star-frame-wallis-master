package log

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(p []byte) (int, error) {
	return rr.ResponseWriter.Write(p)
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var reqBody []byte
		if r.Body != nil && r.Body != http.NoBody {
			reqBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		rr := newResponseRecorder(w)
		next(rr, r)

		duration := time.Since(start)

		logx.WithContext(r.Context()).Infow(
			"API request",
			logx.Field("time", start.Format(time.DateTime)),
			logx.Field("method", r.Method),
			logx.Field("path", r.URL.Path),
			logx.Field("query", r.URL.RawQuery),
			logx.Field("reqBody", string(reqBody)),
			logx.Field("status", rr.statusCode),
			logx.Field("duration", duration.Milliseconds()),
		)
	}
}