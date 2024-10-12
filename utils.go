package main

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func debugLog(messages ...string) {
	log.Print(messages)
}

func getEnvValue(envName string) (string, bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file: ", err)
		return "", false
	}
	return os.Getenv(envName), true
}

func removeDigits(str string) string {
	regex := regexp.MustCompile("[0-9]+")
	return regex.ReplaceAllString(str, "")
}
