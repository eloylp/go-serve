package logging

import (
	"fmt"
	"log"
)

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Error(msg string)
	Errorf(msg string, args ...interface{})
}

type consoleLogger struct{}

func NewConsoleLogger() *consoleLogger { //nolint:golint
	return &consoleLogger{}
}

func (c *consoleLogger) Infof(msg string, args ...interface{}) {
	log.Println(fmt.Sprintf(msg, args...))
}

func (c *consoleLogger) Info(msg string) {
	log.Println(msg)
}

func (c *consoleLogger) Errorf(msg string, args ...interface{}) {
	log.Println(fmt.Sprintf(msg, args...))
}

func (c *consoleLogger) Error(msg string) {
	log.Println(msg)
}
