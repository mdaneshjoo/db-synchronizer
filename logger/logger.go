package logger

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/mdaneshjoo/db-synchronizer/config"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

func toFile(filename string) *os.File {
	rootDir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	logFilePath := filepath.Join(rootDir, filename)

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return logFile
}

func NewLogger() *Logger {
	cfg := config.NewConfig()
	logCfg := cfg.Logging

	infoFlag := flag.Bool("info", false, "Enable info logging")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	filenameFlag := flag.String("file", "", "File path (If pass instead of stderr will write in file)")
	flag.Parse()

	filename := logCfg.File
	if *filenameFlag != "" {
		filename = *filenameFlag
	}

	out := os.Stderr
	if filename != "" {
		out = toFile(filename)
	}

	errorLogger := log.New(out, "ERROR: ", log.Ldate|log.Ltime)

	var infoLogger *log.Logger
	var debugLogger *log.Logger

	if *infoFlag || *logCfg.Info {
		infoLogger = log.New(out, "INFO: ", log.Ldate|log.Ltime)
	}
	if *debugFlag || *logCfg.Debug {
		debugLogger = log.New(out, "DEBUG: ", log.Ldate|log.Ltime)
	}

	return &Logger{
		Info:  infoLogger,
		Error: errorLogger,
		Debug: debugLogger,
	}
}

func (l *Logger) Infoln(msg string) {
	if l.Info != nil {
		l.Info.Println(msg)
	}
}
func (l *Logger) Infof(format string, err ...any) {
	if l.Info != nil {
		l.Info.Printf(format, err...)
	}
}
func (l *Logger) Debugf(format string, err ...any) {
	if l.Debug != nil {
		l.Debug.Printf(format, err...)
	}
}

func (l *Logger) Debugln(msg string) {
	if l.Debug != nil {
		l.Debug.Println(msg)
	}
}

func (l *Logger) Fatalf(format string, err ...any) {
	l.Error.Fatalf(format, err...)
}
