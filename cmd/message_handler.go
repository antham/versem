package cmd

import (
	"io"
	"log"
	"os"

	"github.com/fatih/color"
)

func failOnFprintError(c int, err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type messageHandler struct {
	exit         func(int)
	stdoutWriter io.Writer
	stderrWriter io.Writer
}

func newMessageHandler() messageHandler {
	return messageHandler{os.Exit, os.Stdout, os.Stderr}
}

func (m messageHandler) error(err error) {
	failOnFprintError(color.New(color.FgRed).Fprint(m.stderrWriter, err.Error()+"\n"))
}

func (m messageHandler) errorFatal(err error) {
	failOnFprintError(color.New(color.FgRed).Fprint(m.stderrWriter, err.Error()+"\n"))
	m.exit(1)
}

func (m messageHandler) errorFatalStr(err string, args ...interface{}) {
	failOnFprintError(color.New(color.FgRed).Fprintf(m.stderrWriter, err+"\n", args...))
	m.exit(1)
}

func (m messageHandler) success(str string, args ...interface{}) {
	failOnFprintError(color.New(color.FgGreen).Fprintf(m.stdoutWriter, str+"\n", args...))
}
