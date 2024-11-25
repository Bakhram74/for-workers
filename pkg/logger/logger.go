package logger

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLog(level string) {
	var l zerolog.Level
	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)

	zerolog.TimeFieldFormat = time.RFC1123Z

	// this is equivalent of Lshortfile
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}

	Log = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()
}

func logCallerDesable() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
