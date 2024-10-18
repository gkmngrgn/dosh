package main

import (
	"log"
	"os"
)

type Logger struct {
	verbose int
}

func NewLogger(verbose int) *Logger {
	log.SetOutput(os.Stdout)

	if verbose > 3 {
		verbose = 0
	}

	return &Logger{
		verbose: verbose,
	}
}

func (l *Logger) logDefault(message string) {
	if l.verbose >= 1 {
		log.Println(message)
	}
}

func (l *Logger) logDetailed(message string) {
	if l.verbose >= 2 {
		log.Println(message)
	}
}

func (l *Logger) logDebug(message string) {
	if l.verbose >= 3 {
		log.Println(message)
	}
}

func initLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
