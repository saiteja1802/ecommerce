package logger

import (
	"log/slog"
	"os"
)

var L *slog.Logger

func Init() {
	L = slog.New(slog.NewJSONHandler(os.Stderr, nil))
}


