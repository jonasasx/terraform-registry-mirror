package server

import (
	"fmt"
	"io"
	"runtime"
)

type logger struct {
	writer io.Writer
}

func (l logger) warn(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	_, _ = l.writer.Write([]byte(fmt.Sprintf("[WARN] %s:%d %s\n", filename, line, msg)))
}

func (l logger) error(msg string) {
	_, filename, line, _ := runtime.Caller(1)
	_, _ = l.writer.Write([]byte(fmt.Sprintf("[ERROR] %s:%d %s\n", filename, line, msg)))
}
