package log

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
}

func LogError(msg string, err error) {
	if err != nil {
		errorLogger.Printf("%s %v", msg, err)
	}
}

func LogInfo(msg string) {
	infoLogger.Println(msg)
}
