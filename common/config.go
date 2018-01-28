package common

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/op/go-logging"
)

type logSection struct {
	File      string
	Level     logging.Level
	Formatter string
}

// Config Initialization
type BasicConfig struct {
	Log logSection
}

func InitConfig(config interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(os.Args[0])
	buffer.WriteString(".toml")

	if _, err := toml.DecodeFile(buffer.String(), config); err != nil {
		fmt.Printf("Error ocurred while decoding config file. Error: %s\n", err)
		os.Exit(-1)

	}
}
