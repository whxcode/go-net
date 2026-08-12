// 日志模块
package logs

import (
	"io"
	"os"
	"time"

	config "go-net/im/src"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func initLogsConfig() io.Writer {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	err := os.MkdirAll(
		config.ConfigData.Log.PathDir,
		0o755,
	)
	if err != nil {
		panic(err)
	}

	if config.ConfigData.Log.Level&config.Log != 0 {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if config.ConfigData.Log.Level&config.Info != 0 {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if config.ConfigData.Log.Level&config.WARN != 0 {
		logrus.SetLevel(logrus.WarnLevel)
	}

	if config.ConfigData.Log.Level&config.ERROR != 0 {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	// 30MB
	logger := MakeDateWrite(config.ConfigData.Log.PathDir, 30)

	logger.rotate()

	logrus.SetOutput(io.MultiWriter(logger, os.Stdout))

	return logrus.StandardLogger().Out
}

func LoggerMiddleware() gin.HandlerFunc {
	initLogsConfig()

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()

		filed := logrus.Fields{
			"time":      start.Format("2006-01-02 15:04:05"),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"status":    c.Writer.Status(),
			"latency":   latency.String(),
			"clientIP":  c.ClientIP(),
			"userAgent": c.Request.UserAgent(),
		}

		if status >= 500 {
			logrus.WithFields(filed).Error("HTTP Request")
		} else if status >= 400 {
			logrus.WithFields(filed).Warn("HTTP Request")
		} else {
			logrus.WithFields(filed).Info("HTTP Request")
		}
		// 记录到 logrus（自动 JSON）
	}
}
