package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func LogMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	raw := c.Request.URL.RawQuery
	if raw != "" {
		path = path + "?" + raw
	}

	c.Next()

	status := c.Writer.Status()

	var msg = "incoming request"
	var attrs = []slog.Attr{
		slog.String("method", c.Request.Method),
		slog.String("path", path),
		slog.Int("status", status),
		slog.Duration("latency", time.Since(start)),
		slog.String("request_id", c.Request.Header.Get("X-Request-ID")),
		slog.String("route", c.FullPath()),
		slog.String("client_ip", c.ClientIP()),
		slog.String("user_agent", c.Request.UserAgent()),
		slog.Int("body_size", c.Writer.Size()),
	}
	var level slog.Level
	switch {
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		level = slog.LevelWarn
		msg = c.Errors.String()
		if msg == "" {
			msg = fmt.Sprintf("client error: %d %s", status, strings.ToLower(http.StatusText(status)))
		}
	case status >= http.StatusInternalServerError:
		level = slog.LevelError
		msg = c.Errors.String()
		if msg == "" {
			msg = fmt.Sprintf("server error: %d %s", status, strings.ToLower(http.StatusText(status)))
		}
	}
	slog.LogAttrs(c.Request.Context(), level, msg, attrs...)
}

func RecoveryMiddleware(c *gin.Context) {
	defer func() {
		if rcv := recover(); rcv != nil {
			slog.Error("recovered from panic", "error", rcv)
			err, isErr := rcv.(error)
			if !isErr {
				err = fmt.Errorf("%v", rcv)
			}
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}()
	c.Next()

}
