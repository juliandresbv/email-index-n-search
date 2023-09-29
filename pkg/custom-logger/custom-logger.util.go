package customlogger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func NewLogger() *log.Logger {
	_, file, _, ok := runtime.Caller(1)

	if !ok {
		file = "main.go"
	} else {
		file = filepath.Base(file)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMsg := fmt.Sprintf("[%s] %s -- ", file, timestamp)

	return log.New(os.Stdout, formattedMsg, 0)
}
