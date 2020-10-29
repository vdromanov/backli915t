package main

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

type Logger struct {
	*log.Logger
	outputs []io.Writer
}

var (
	DebugLog *Logger
	InfoLog  *Logger
)

func InitLogger(dest io.Writer, prefix string, flag int) *Logger {
	return &Logger{
		log.New(dest, prefix, flag),
		[]io.Writer{dest},
	}
}

func (logger *Logger) AddDestination(destinations ...io.Writer) {
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
	DebugLog = InitLogger(mainlogSource, "DEBUG: ", log.LstdFlags|log.Lshortfile)
	InfoLog = InitLogger(mainlogSource, "", log.LstdFlags)
}
