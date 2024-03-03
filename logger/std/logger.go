package std

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/RustGrub/FunnyGoService/config"
)

const (
	dtMask     = "2006-01-02 15:04:05.0000"
	basePrefix = ""
	baseFlag   = 0
)

const defaultLevel = InfoLevel

const (
	DebugLevel   = "DEBUG"
	InfoLevel    = "INFO"
	WarnignLevel = "WARNING"
	ErrLevel     = "ERROR"
	FatalLevel   = "FATAL"
)

const (
	DEBUG   = 40
	INFO    = 30
	WARNING = 20
	ERROR   = 10
	FATAL   = 0
)

// Украдено

type Logger struct {
	env    string
	level  int
	logger *log.Logger
}

func New(conf *config.Config) *Logger {
	level := conf.Logger.Level
	if level == "" {
		level = defaultLevel
	}

	var logLevelMap = map[string]int{
		DebugLevel:   DEBUG,
		InfoLevel:    INFO,
		WarnignLevel: WARNING,
		ErrLevel:     ERROR,
		FatalLevel:   FATAL,
	}
	logger := log.New(os.Stdout, basePrefix, baseFlag)

	l := &Logger{
		env:    conf.Environment,
		level:  logLevelMap[level],
		logger: logger,
	}

	return l
}

func (l *Logger) Close() {}

func (l *Logger) Debug(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= DEBUG {
		l.logger.Print(now+" DEBUG ", l.getFuncName(), " ", v)
	}
}

func (l *Logger) Info(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= INFO {
		l.logger.Print(now+" INFO ", v)
	}
}

func (l *Logger) Warning(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= WARNING {
		l.logger.Print(now+" WARNING ", l.getFuncName(), " ", v)
	}
}

func (l *Logger) Error(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= ERROR {
		l.logger.Print(now+" ERROR ", l.getFuncName(), " ", v)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	now := time.Now().Format(dtMask)
	l.logger.Fatal(now+" FATAL ", l.getFuncName(), " ", v)
}

func (l *Logger) getFuncName() string {
	var skip int = 3 // nolint: revive

	var buffer bytes.Buffer
	pc := make([]uintptr, 10)
	runtime.Callers(skip, pc)
	frame, _ := runtime.CallersFrames(pc).Next()
	function := frame.Function
	line := frame.Line
	buffer.WriteString(function)
	buffer.WriteString(fmt.Sprintf(":%d", line))

	return filepath.Base(buffer.String())
}
