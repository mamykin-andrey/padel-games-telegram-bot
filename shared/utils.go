package shared

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
)

func DebugLog(messages ...string) {
	log.Print(messages)
}

func GetEnvValue(envName string) (string, bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file: ", err)
		return "", false
	}
	return os.Getenv(envName), true
}

func RemoveDigits(str string) string {
	regex := regexp.MustCompile("[0-9]+")
	return regex.ReplaceAllString(str, "")
}

func TryToParseDate(dateStr string) (time.Time, bool) {
	var date time.Time
	var err error

	date, err = time.Parse("02.01.2006", dateStr)
	if err == nil {
		return date, true
	}

	date, err = time.Parse("02.01.06", dateStr)
	if err == nil {
		if date.Year() > time.Now().Year() {
			date = date.AddDate(-100, 0, 0)
		}
		return date, true
	}

	return date, false
}

func IsDatePassed(date time.Time) bool {
	today := time.Now()
	return date.Before(today)
}
