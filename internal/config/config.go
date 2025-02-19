package config

import (
	"flag"
	"fmt"
	"os"
)

var BotToken string
var Logging bool
var DatabasePath string

func LoadAndValidate() {
	flag.StringVar(&BotToken, "token", "", "Bot's secret token")
	flag.StringVar(&DatabasePath, "db", "", "Path to a sqlite database")
	flag.BoolVar(&Logging, "enable-logging", true, "Enable extensive logging")
	flag.Parse()
	if BotToken == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing requried -token parameter\n")
		flag.Usage()
		os.Exit(1)
	}

	if DatabasePath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing requried -db parameter\n")
		flag.Usage()
		os.Exit(1)
	}
}
