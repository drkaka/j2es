package logger

import (
	"time"

	"github.com/uber-go/zap"
)

var (
	// Log the zap logger
	Log zap.Logger

	debug string
)

func rfc3339NanoFormatter(key string) zap.TimeFormatter {
	return zap.TimeFormatter(func(t time.Time) zap.Field {
		return zap.String(key, t.Format(time.RFC3339Nano))
	})
}

func init() {
	Log = zap.New(
		zap.NewJSONEncoder(
			rfc3339NanoFormatter("timestamp"), // human-readable timestamps
			zap.MessageKey("msg"),             // customize the message key
			zap.LevelString("lvl"),            // stringify the log level
		),
	)

	if debug == "debug" {
		Log.SetLevel(zap.DebugLevel)
	}
}
