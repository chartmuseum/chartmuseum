package chartmuseum

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger handles all logging from application
	Logger struct {
		*zap.SugaredLogger
	}

	logLevel string

	loggingFn func(level logLevel, msg string, keysAndValues ...interface{})
)

const (
	debugLevel logLevel = "DEBUG"
	infoLevel  logLevel = "INFO"
	warnLevel  logLevel = "WARN"
	errorLevel logLevel = "ERROR"
)

// NewLogger creates a new Logger instance
func NewLogger(json bool, debug bool) (*Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	config.Development = false
	config.DisableCaller = true
	if json {
		config.Encoding = "json"
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	if !debug {
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	logger, err := config.Build()
	if err != nil {
		return new(Logger), err
	}
	defer logger.Sync()
	return &Logger{logger.Sugar()}, nil
}

// Debugc wraps Debugw provided by zap, adding data from gin request context
func (logger *Logger) Debugc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Debugw(msg, keysAndValues...)
}

// Infoc wraps Infow provided by zap, adding data from gin request context
func (logger *Logger) Infoc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Infow(msg, keysAndValues...)
}

// Warnc wraps Warnw provided by zap, adding data from gin request context
func (logger *Logger) Warnc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Warnw(msg, keysAndValues...)
}

// Errorc wraps Errorw provided by zap, adding data from gin request context
func (logger *Logger) Errorc(c *gin.Context, msg string, keysAndValues ...interface{}) {
	msg, keysAndValues = transformLogcArgs(c, msg, keysAndValues)
	logger.Errorw(msg, keysAndValues...)
}

// transformLogcArgs prefixes msg with RequestCount and adds RequestId to keysAndValues
func transformLogcArgs(c *gin.Context, msg string, keysAndValues []interface{}) (string, []interface{}) {
	if reqCount, exists := c.Get("RequestCount"); exists {
		msg = fmt.Sprintf("[%s] %s", reqCount, msg)
		keysAndValues = append(keysAndValues, "reqID", c.MustGet("RequestId"))
	}
	return msg, keysAndValues
}

func loggingMiddleware(logger *Logger) gin.HandlerFunc {
	var requestCount int64
	return func(c *gin.Context) {
		reqCount := strconv.FormatInt(atomic.AddInt64(&requestCount, 1), 10)
		c.Set("RequestCount", reqCount)

		reqPath := c.Request.URL.Path
		logger.Debugc(c, fmt.Sprintf("Incoming request: %s", reqPath))
		start := time.Now()
		c.Next()

		msg := "Request served"
		status := c.Writer.Status()

		meta := []interface{}{
			"path", reqPath,
			"comment", c.Errors.ByType(gin.ErrorTypePrivate).String(),
			"latency", time.Now().Sub(start),
			"clientIP", c.ClientIP(),
			"method", c.Request.Method,
			"statusCode", status,
		}

		switch {
		case status == 200 || status == 201 || status == 301:
			logger.Infoc(c, msg, meta...)
		case status == 404:
			logger.Warnc(c, msg, meta...)
		default:
			logger.Errorc(c, msg, meta...)
		}
	}
}

/*
contextLoggingFn creates a loggingFn to be used in
places that do not necessarily need access to the gin context
*/
func (server *Server) contextLoggingFn(c *gin.Context) loggingFn {
	return func(level logLevel, msg string, keysAndValues ...interface{}) {
		switch level {
		case debugLevel:
			server.Logger.Debugc(c, msg, keysAndValues...)
		case infoLevel:
			server.Logger.Infoc(c, msg, keysAndValues...)
		case warnLevel:
			server.Logger.Warnc(c, msg, keysAndValues...)
		case errorLevel:
			server.Logger.Errorc(c, msg, keysAndValues...)
		}
	}
}

func init() {
	logrus.SetLevel(logrus.WarnLevel) // silence logs from zsais/go-gin-prometheus
}
