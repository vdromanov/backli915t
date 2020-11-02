package multilog

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

// Logger is a log package's logger with a slice of it's outputs
type Logger struct {
	*log.Logger
	outputs []io.Writer
}

var (
	// Debug messages are a developer's ones. Actually, they are visible in syslog
	Debug *Logger
	// Info messages are supposed to be visible for user.
	Info *Logger
)

// New makes a logger (a log package's one) with a slice of logging sources
func New(dest io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		log.New(dest, prefix, flag),
		[]io.Writer{dest},
	}
}

// AddOutput appends a io.Writer to logger's outputs
func (logger *Logger) AddOutput(destinations ...io.Writer) {
	logger.outputs = append(logger.outputs, destinations...)
	writers := io.MultiWriter(logger.outputs...)
	logger.SetOutput(writers)
}

func init() {
	var mainlogSource io.Writer
	sysLog, err := syslog.New(syslog.LOG_NOTICE, "backli915t")
	if err != nil {
		mainlogSource = os.Stdout
	} else {
		mainlogSource = sysLog
	}
	Debug = New(mainlogSource, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	Info = New(mainlogSource, "", log.LstdFlags)
}
