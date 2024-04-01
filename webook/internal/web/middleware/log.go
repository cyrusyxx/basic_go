package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type LogMiddlewareBuilder struct {
	logFn         func(ctx context.Context, al AccessLog)
	allowReqBody  bool
	allowRespBody bool
}

type AccessLog struct {
	Path     string        `json:"path"`
	Method   string        `json:"method"`
	State    int           `json:"state"`
	ReqBody  string        `json:"req_body"`
	RespBody string        `json:"resp_body"`
	Duration time.Duration `json:"duration"`
}

func NewLogMiddlewareBuilder(logFn func(ctx context.Context, al AccessLog)) *LogMiddlewareBuilder {
	return &LogMiddlewareBuilder{
		logFn: logFn,
	}
}

func (l *LogMiddlewareBuilder) AllowReqBody() *LogMiddlewareBuilder {
	l.allowReqBody = true
	return l
}

func (l *LogMiddlewareBuilder) AllowRespBody() *LogMiddlewareBuilder {
	l.allowRespBody = true
	return l
}

func (l *LogMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if len(path) > 1024 {
			path = path[:1024]
		}
		method := ctx.Request.Method
		al := AccessLog{
			Path:   path,
			Method: method,
		}
		if l.allowReqBody {
			body, _ := ctx.GetRawData()
			if len(body) > 2048 {
				body = body[:2048]
			}
			al.ReqBody = string(body)
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		}

		start := time.Now()

		if l.allowRespBody {
			ctx.Writer = &responseWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}

		defer func() {
			al.Duration = time.Since(start)
			// Get response body
			l.logFn(ctx, al)
		}()

		ctx.Next()
	}
}

type responseWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

func (r *responseWriter) Write(b []byte) (int, error) {
	r.al.RespBody = string(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseWriter) WriteString(s string) (int, error) {
	r.al.RespBody = s
	return r.ResponseWriter.WriteString(s)
}

func (r *responseWriter) WriteHeader(code int) {
	r.al.State = code
	r.ResponseWriter.WriteHeader(code)
}
