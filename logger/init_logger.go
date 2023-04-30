package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "Jan _2 15:04:05.00000",
	}

	zlog.Logger = zerolog.
		New(writer).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}
