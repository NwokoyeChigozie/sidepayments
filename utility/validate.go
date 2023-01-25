package utility

import (
	"net/mail"
	"os"
	"regexp"
	"strings"
)

func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func PhoneValid(phone string) (string, bool) {
	if phone == "" || len(phone) < 5 {
		return phone, false
	}
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	if strings.Contains(phone, "+234") {
		phone = strings.Replace(phone, "+234", "0", 1)
	} else if strings.Contains(phone, "234") {
		phone = strings.Replace(phone, "234", "0", 1)
	}

	return phone, re.MatchString(phone)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
