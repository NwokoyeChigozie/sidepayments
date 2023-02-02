package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Tabler interface {
	TableName() string
}

type jsonmap map[string]interface{}

// Value Marshal
func (a jsonmap) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *jsonmap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
