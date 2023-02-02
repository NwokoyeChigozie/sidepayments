package utility

import (
	"time"
)

func FormatDate(date, currentISOFormat, newISOFormat string) (string, error) {
	t, err := time.Parse(currentISOFormat, date)
	if err != nil {
		return date, err
	}
	return t.Format(newISOFormat), nil
}
