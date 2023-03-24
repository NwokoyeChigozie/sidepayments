package utility

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hash(str string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashed), err
}

func CompareHash(str string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str)) == nil
}

func ShaHash(str string) (string, error) {
	passSha1 := sha1.New()
	_, err := passSha1.Write([]byte(str))
	if err != nil {
		return str, err
	}

	getSha1 := passSha1.Sum(nil)
	return fmt.Sprintf("%x", getSha1), nil
}

func Sha512Hmac(secret string, body []byte) string {
	hmacSecret := []byte(secret)
	h := hmac.New(sha512.New, hmacSecret)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}
