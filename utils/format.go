package utils

import (
	"regexp"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer *message.Printer
var startingZeroRegex *regexp.Regexp

func init() {
	printer = message.NewPrinter(language.Spanish)
	startingZeroRegex = regexp.MustCompile("^0+")
}

func FormatCurrency(amount interface{}) string {
	return printer.Sprintf("$%v", amount)
}

func RemoveStartingZeros(raw string) string {
	return startingZeroRegex.ReplaceAllString(raw, "")
}
