package utility

import (
	"net/url"
	"strings"
)

func URLDecode(encodedString string) (string, error) {
	decoded, err := url.QueryUnescape(encodedString)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func UrlHasQuery(urlString string) (bool, error) {
	urlS, err := URLDecode(urlString)
	if err != nil {
		return false, err
	}

	u, err := url.Parse(urlS)
	if err != nil {
		panic(err)
	}

	queryParameters := u.Query()
	if len(queryParameters) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func AddQueryParam(urlStr *string, paramKey string, paramValue string) error {
	// Parse the URL
	u, err := url.Parse(*urlStr)
	if err != nil {
		return err
	}

	// Get the query parameters as a map
	queryParams, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}

	// Add or update the parameter with the given key and value
	queryParams.Set(paramKey, paramValue)

	// Encode the query parameters and rebuild the URL
	u.RawQuery = queryParams.Encode()

	*urlStr = u.String()
	return nil
}

func Stripslashes(s string) string {
	return strings.ReplaceAll(s, "\\", "")
}
