package node

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	Logger = slog.New(jsonHandler)
}
