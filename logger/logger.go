package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info *log.Logger
	Err  *log.Logger
}

func New() *Logger {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Logger{
		Info: infoLog,
		Err:  errLog,
	}
}
