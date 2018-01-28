package common

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/op/go-logging"
)

func InitLogger(config *BasicConfig) {

	file, err := os.OpenFile(config.Log.File,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)

	if err != nil {
		fmt.Printf("Failed to open log file. Error: %s\n", err)
		os.Exit(-1)
	}

	logHandler := logging.NewLogBackend(file, "", 0)
	format := logging.MustStringFormatter(config.Log.Formatter)
	logFormatter := logging.NewBackendFormatter(logHandler, format)
	logging.SetBackend(logFormatter)
	logging.SetLevel(config.Log.Level, "")
}

func InitLoggerTest() {
	logHandler := logging.NewLogBackend(ioutil.Discard, "", 0)
	logging.SetBackend(logHandler)
}
