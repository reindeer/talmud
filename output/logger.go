package output

import (
	"fmt"
	"log"
)

type logWriter struct {
}

func (logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func Fatalf(format string, v ...interface{}) {
	log.Fatalf("\033[31m"+format+"\033[0m", v...)
}

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
