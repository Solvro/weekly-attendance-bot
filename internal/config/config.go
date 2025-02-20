package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var BotToken string
var Logging bool
var DatabasePath string

func LoadAndValidate() {
	if err := godotenv.Load(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, ".env configuration file not present")
		os.Exit(1)
	}

	BotToken = os.Getenv("TOKEN")
	DatabasePath = os.Getenv("DB")
	Logging = os.Getenv("ENABLE_LOGGING") != ""

	if BotToken == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing requried environment variable: 'TOKEN'\n")
		os.Exit(1)
	}

	if DatabasePath == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Missing requried environment variable: 'DB'\n")
		os.Exit(1)
	}
}
