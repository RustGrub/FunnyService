package logger

import (
	"fmt"
	"github.com/RustGrub/FunnyGoService/config"
	"github.com/RustGrub/FunnyGoService/logger/std"
)

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	// Fatal writes log message with fatal level and os.Exit(1) after
	Fatal(v ...interface{})
	Close()
}

func NewLogger(cfg *config.Config) Logger {
	switch cfg.Environment {
	case "prod":
		// Типа есть
		return nil
	default:
		fmt.Println("Local")
		return std.New(cfg)
	}
}
