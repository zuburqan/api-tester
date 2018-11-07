package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func Inject(idValue, text string) string {
	return strings.Replace(text, "#{id}", idValue, -1)
}

func GenerateID(id int, destination string) string {
	return strings.ToUpper(destination) + "-" + "API-TESTER-" + strconv.Itoa(id)
}

func RemoveIDFrom(text, destination string) string {
	regexString := fmt.Sprintf("%s-API-TESTER-[[:digit:]]+", strings.ToUpper(destination))
	re := regexp.MustCompile(regexString)
	return re.ReplaceAllString(text, "ID")
}
