package utility

import (
	"net/mail"
	"os"

	"github.com/nyaruka/phonenumbers"
)

func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func PhoneValid(phone string) (string, bool) {
	parsed, err := phonenumbers.Parse(phone, "")
	if err != nil {
		return phone, false
	}

	if !phonenumbers.IsValidNumber(parsed) {
		return phone, false
	}

	formattedNum := phonenumbers.Format(parsed, phonenumbers.NATIONAL)
	return formattedNum, true
}

func MakePhoneNumber(number string, countryCode string) (string, error) {
	parsed, err := phonenumbers.Parse(number, countryCode)
	if err != nil {
		return number, err
	}

	formattedNum := phonenumbers.Format(parsed, phonenumbers.NATIONAL)
	return formattedNum, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
