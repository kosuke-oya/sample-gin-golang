package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	zap "go.uber.org/zap"
)

type LogInfo struct {
	RemoteAddr  string
	ContentType string
	Path        string
	Query       string
	Method      string
	Body        string
	HttpStatus  int
	Err         error
	ReturnJson  string
}

type MyResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *MyResponseWriter) Write(data []byte) (int, error) {
	// Capture the response body
	w.body.Write(data)
	// Call the original Write method to write to the actual response writer
	return w.ResponseWriter.Write(data)
}

func Logger(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		fmt.Println("starting logger", start)
		w := &MyResponseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = w
		c.Next()

		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// CloudLoggingに送信用のメッセージを定義
		heaserBytes, err := json.Marshal(c.Request.Header)
		if err != nil {
			fmt.Println("failed to marshal header")
		}
		headerStr := string(heaserBytes)

		// log here
		l.Info("Request",
			zap.String("uuid", xid.New().String()),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("content_length", c.Request.ContentLength),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("request_body", string(bodyBytes)),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("elapsed", time.Since(start)),
			zap.String("header", headerStr),
			zap.Int("response_size(bytes)", c.Writer.Size()),
		)
	}

}
