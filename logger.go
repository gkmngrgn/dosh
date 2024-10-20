package main

import (
	"log"
	"sync"
)

var (
	instance *Logger
	once     sync.Once
)

type Logger struct {
	verbose int
}

func GetLogger() *Logger {
	once.Do(func() {
		log.SetFlags(log.Ldate | log.Ltime)
		instance = &Logger{}
	})

	return instance
}

func (l *Logger) setVerbosity(v int) {
	if v > 3 {
		v = 0
	}

	l.verbose = v
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

func (l *Logger) logError(message string) {
	if l.verbose >= 3 {
		log.Println(message)
	}
}
