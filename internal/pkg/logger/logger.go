package logger

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	Info  *log.Logger
	Warn  *log.Logger
	Debug *log.Logger
	Error *log.Logger
)

func Setup(logLevel string) error {
	logLevel = strings.ToLower(logLevel)

	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	debugHandle := ioutil.Discard
	errorHandle := os.Stderr

	if logLevel == "info" {
		infoHandle = os.Stdout
	}
	if logLevel == "info" || logLevel == "warn" {
		warnHandle = os.Stdout
	}
	if logLevel == "info" || logLevel == "warn" || logLevel == "debug" {
		debugHandle = os.Stdout
	}

	Debug = log.New(debugHandle, "\u001b[46;1m\u001b[37mDEB:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	Info = log.New(infoHandle, "\u001b[42;1m\u001b[37mINF:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	Warn = log.New(warnHandle, "\u001b[43;1m\u001b[38mWARN:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)
	Error = log.New(errorHandle, "\u001b[41;1m\u001b[37mERR:\u001b[0m ", log.Ldate|log.Lmicroseconds|log.LUTC|log.Llongfile)

	return nil
}
