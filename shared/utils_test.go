package shared

import (
	"testing"
	"time"
)

func TestRemoveDigitsFromEmptyString(t *testing.T) {
	result := RemoveDigits("")

	if result != "" {
		t.Fatalf(`The expected result is "", but the actual result is %q`, result)
	}
}

func TestRemoveDigitsFromStringWithNoDigits(t *testing.T) {
	input := "hello"

	result := RemoveDigits(input)

	if result != input {
		t.Fatalf(`The expected result is %q", but the actual result is %q`, input, result)
	}
}

func TestRemoveDigitsFromStringWithDigitsOnly(t *testing.T) {
	result := RemoveDigits("12345")

	if result != "" {
		t.Fatalf(`The expected result is "", but the actual result is %q`, result)
	}
}

func TestRemoveDigitsFromStringWithDigits(t *testing.T) {
	result := RemoveDigits("1he2llo3")

	if result != "hello" {
		t.Fatalf(`The expected result is "hello", but the actual result is %q`, result)
	}
}

func TestTryToParseDateFromInvalidString(t *testing.T) {
	result, ok := TryToParseDate("fregre")

	if ok {
		t.Fatalf(`The expected result is "fail", but the actual result is %q`, result)
	}
}

func TestTryToParseDateFromFullLayout(t *testing.T) {
	result, ok := TryToParseDate("10.10.2022")

	expectedDate := time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC)
	if !ok || result != expectedDate {
		t.Fatalf(`The expected result is "10.10.2022", but the actual result is %q`, result)
	}
}

func TestTryToParseDateFromShortLayout(t *testing.T) {
	result, ok := TryToParseDate("10.10.2022")

	expectedDate := time.Date(2022, 10, 10, 0, 0, 0, 0, time.UTC)
	if !ok || result != expectedDate {
		t.Fatalf(`The expected result is "10.10.22", but the actual result is %q`, result)
	}
}

func TestIsDatePassedForFutureDate(t *testing.T) {
	date := time.Now()
	date = date.AddDate(0, 0, 1)

	passed := IsDatePassed(date)

	if passed {
		t.Fatalf(`The expected result is "false", but the actual result is "true"`)
	}
}

func TestIsDatePassedForPastDate(t *testing.T) {
	date := time.Now()
	date = date.AddDate(0, 0, -1)

	passed := IsDatePassed(date)

	if !passed {
		t.Fatalf(`The expected result is "true", but the actual result is "false"`)
	}
}
