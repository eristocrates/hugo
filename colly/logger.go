package main

import (
	"os"
	"runtime/debug"

	"github.com/rs/zerolog"
)

var logFile *os.File
var jsonLogFile *os.File

func InitializeLogger() *zerolog.Logger {
	buildInfo, _ := debug.ReadBuildInfo()

	logFilePath := "../logs/colly.log"
	jsonLogFilePath := "../logs/colly.json"
	var err error
	logFile, err = os.OpenFile(
		logFilePath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	jsonLogFile, err = os.OpenFile(
		jsonLogFilePath,
		os.O_TRUNC|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	lcw := &CommaWriter{file: logFile}
	jcw := &CommaWriter{file: jsonLogFile}
	multi := zerolog.MultiLevelWriter(lcw, jcw)
	logger := zerolog.New(multi).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	return &logger
}
