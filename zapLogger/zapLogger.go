package zapLogger

import (
	"awesomeProject10/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Log *zap.Logger

func Init() {
	lumberjackRotate := lumberjack.Logger{
		Filename:   "calc.log",
		MaxSize:    5,
		MaxBackups: 2,
		MaxAge:     5,
		Compress:   true,
	}
	writeSyncer := zapcore.AddSync(&lumberjackRotate)

	consoleConfig := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		MessageKey:  "msg",
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime:  zapcore.ISO8601TimeEncoder,
	}

	jsonConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		NameKey:      "zapLogger",
		CallerKey:    "caller",
		MessageKey:   "msg",
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeTime:   zapcore.ISO8601TimeEncoder,
	}
	jsonCore := zapcore.NewCore(zapcore.NewJSONEncoder(jsonConfig), writeSyncer, zapcore.DebugLevel)
	consoleCore := zapcore.NewCore(zapcore.NewConsoleEncoder(consoleConfig), os.Stdout, zapcore.InfoLevel)

	core := zapcore.NewTee(jsonCore, consoleCore)
	sampledCore := zapcore.NewSamplerWithOptions(core, time.Second, 50, 1)

	Log = zap.New(sampledCore)

}

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		ip := c.ClientIP()
		after := time.Since(start)

		entry := models.OkLogs{
			Path:     path,
			Raw:      raw,
			IP:       ip,
			Duration: after.String(),
		}

		Log.Sugar().Debugw("request info",
			"path", path,
			"query", raw,
			"ip", ip,
			"duration", after,
		)
		models.LogChan <- entry
	}
}
