package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/shortuuid/v3"

	"jolly/backend/common/log"
)

const (
	TestNameHeader          = "TestName"
	CorrelationIDHttpHeader = "Correlation-ID"
)

func useMiddlewares(e *echo.Echo) {
	e.Use(
		middleware.ContextTimeout(10*time.Second),
		middleware.Recover(),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				ctx := req.Context()

				reqCorrelationID := req.Header.Get(CorrelationIDHttpHeader)
				if reqCorrelationID == "" {
					reqCorrelationID = shortuuid.New()
				}

				logger := slog.With("correlation_id", reqCorrelationID)

				if testName := c.Request().Header.Get(TestNameHeader); testName != "" {
					logger = logger.With("test_name", testName)
				}

				ctx = log.ToContext(ctx, logger)
				ctx = log.ContextWithCorrelationID(ctx, reqCorrelationID)
				c.SetRequest(req.WithContext(ctx))
				c.Response().Header().Set(CorrelationIDHttpHeader, reqCorrelationID)

				return next(c)
			}
		},
		requestLogMiddleware,
	)
}

type bodyCapturingWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyCapturingWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyCapturingWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyCapturingWriter) Flush() {
	err := http.NewResponseController(w.ResponseWriter).Flush()
	if err != nil && !errors.Is(err, http.ErrNotSupported) {
		slog.Warn("response writer flush failed", "error", err)
	}
}

func (w *bodyCapturingWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *bodyCapturingWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func requestLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqBody []byte
		if c.Request().Body != nil {
			reqBody, _ = io.ReadAll(c.Request().Body)
		}
		c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody))

		resBody := new(bytes.Buffer)
		mw := io.MultiWriter(c.Response().Writer, resBody)
		c.Response().Writer = &bodyCapturingWriter{Writer: mw, ResponseWriter: c.Response().Writer}

		start := time.Now()
		err := next(c)
		if err != nil {
			c.Error(err)
			err = nil
		}
		duration := time.Since(start)

		ctx := c.Request().Context()

		logger := log.FromContext(ctx).With(
			"URI", c.Request().RequestURI,
			"status", c.Response().Status,
			"method", c.Request().Method,
			"duration", duration.String(),
		)
		if err != nil {
			logger = logger.With("error", err)
		}
		logger = logger.With("request_body", truncateBodyForLog(string(reqBody)))

		body := resBody.String()
		if utf8.ValidString(body) {
			if isDebug := log.FromContext(ctx).Enabled(ctx, slog.LevelDebug); !isDebug {
				body = truncateBodyForLog(body)
			}
			logger = logger.With("response_body", body)
		} else {
			logger = logger.With("response_body", "<binary data>")
		}

		logger.Info("Request done")
		return err
	}
}
