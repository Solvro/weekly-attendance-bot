package config

import (
	"flag"
	"fmt"
	"os"
)

var BotToken string
var Logging bool

func LoadAndValidate() {
	flag.StringVar(&BotToken, "token", "", "Bot's secret token")
	flag.BoolVar(&Logging, "enable-logging", true, "Enable extensive logging")
	flag.Parse()
	if BotToken == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing requried -token parameter\n")
		flag.Usage()
		os.Exit(1)
	}
}
