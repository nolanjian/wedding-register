package logger

import (
	"os"

	"github.com/CodisLabs/codis/pkg/utils/log"
)

var mylogger *log.Logger

func init() {
	fileName := "ll.log"
	logFile, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	mylogger = log.New(logFile, `Wedding `)
}

func GetLogger() *log.Logger {
	return mylogger
}
