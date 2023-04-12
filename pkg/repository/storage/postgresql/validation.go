package postgresql

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/vesicash/payment-ms/external/external_models"
	"github.com/vesicash/payment-ms/external/request"
	"github.com/vesicash/payment-ms/utility"
	"gorm.io/gorm"
)

type Validation interface {
	LogModelData(*utility.Logger)
}

type ValidationError struct {
	Field string
	Error string
}

var (
	regexpEmail      = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ErrEmptyField    = errors.New("Field cannot be empty")
	ErrInvalidEmail  = errors.New("Email is invalid")
	ErrInvalidPass   = errors.New("Field length should be greater than 8")
	ErrNil           = errors.New("Nil")
	ValidationNeeded = "Input validation failed on some fields"
)

type ValidateRequestM struct {
	Logger *utility.Logger
	Test   bool
}

func (vr ValidateRequestM) ValidateRequest(V interface{}) error {

	var err []ValidationError
	if reflect.ValueOf(V).Kind() == reflect.Struct {
		t := reflect.TypeOf(V)
		v := reflect.ValueOf(V)

		for i := 0; i < t.NumField(); i++ {
			FieldT := t.Field(i)
			FieldV := v.Field(i)
			// reflect.ValueOf(V).Field(i).Type()

			validateFields := FieldT.Tag.Get("pgvalidate")
			splitFields := strings.Split(validateFields, ",")
			if validateFields == "_" || validateFields == "" {
				continue
			}

			for j := 0; j < len(splitFields); j++ {
				splitFieldsStr := strings.ToLower(splitFields[j])
				if strings.Contains(splitFieldsStr, "notexists") {
					value, status := ValidateNext(FieldV)
					if status {
						firstSplit := strings.Split(splitFieldsStr, "=")
						if len(firstSplit) == 2 {
							secondSplit := strings.Split(firstSplit[1], "$")
							if len(secondSplit) == 3 {
								dbName := secondSplit[0]
								tableName := secondSplit[1]
								columnName := secondSplit[2]
								if !vr.ValidationCheck(dbName, tableName, "notexists", fmt.Sprintf("%v = ?", columnName), value) {
									err = append(err, ValidationError{
										Field: FieldT.Name,
										Error: fmt.Sprintf("%v exists in %v table", columnName, tableName),
									})
								}

							}

						}
					}
				} else if strings.Contains(splitFieldsStr, "exists") {
					value, status := ValidateNext(FieldV)
					if status {
						firstSplit := strings.Split(splitFieldsStr, "=")
						if len(firstSplit) == 2 {
							secondSplit := strings.Split(firstSplit[1], "$")
							if len(secondSplit) == 3 {
								dbName := secondSplit[0]
								tableName := secondSplit[1]
								columnName := secondSplit[2]
								if !vr.ValidationCheck(dbName, tableName, "exists", fmt.Sprintf("%v = ?", columnName), value) {
									err = append(err, ValidationError{
										Field: FieldT.Name,
										Error: fmt.Sprintf("%v does not exist in %v table", columnName, tableName),
									})
								}

							}

						}
					}
				} else if strings.Contains(splitFieldsStr, "email") {
					if FieldV.String() != "" {
						if !regexpEmail.Match([]byte(FieldV.String())) {
							err = append(err, ValidationError{
								Field: FieldT.Name,
								Error: ErrInvalidEmail.Error(),
							})
						}
					}
				}
			}
		}
	}

	errString := ""
	if len(err) < 1 {
		return nil
	} else {
		for _, v := range err {
			errString += v.Field + ": " + v.Error + " ;"
		}
	}
	return fmt.Errorf(errString)
}

func (vr ValidateRequestM) ValidationCheck(dbName string, table, checkType string, query interface{}, args ...interface{}) bool {
	er := request.ExternalRequest{
		Logger: vr.Logger,
		Test:   vr.Test,
	}
	db := ReturnDatabase(dbName)
	switch dbName {
	case "admin":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "auth":
		status, err := er.SendExternalRequest(request.ValidateOnAuth, external_models.ValidateOnDBReq{
			Table: table,
			Type:  checkType,
			Query: fmt.Sprintf("%v", query),
			Value: args[0],
		})
		if err != nil {
			vr.Logger.Error("error occurred in validation", err.Error())
			return false
		}
		return status.(bool)
	case "notifications":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "payment":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "reminder":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "subscription":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "transaction":
		status, err := er.SendExternalRequest(request.ValidateOnTransactions, external_models.ValidateOnDBReq{
			Table: table,
			Type:  checkType,
			Query: fmt.Sprintf("%v", query),
			Value: args[0],
		})
		if err != nil {
			vr.Logger.Error("error occurred in validation", err.Error())
			return false
		}
		return status.(bool)
	case "verification":
		return checkForConnectedDB(db, table, checkType, query, args...)
	case "cron":
		return checkForConnectedDB(db, table, checkType, query, args...)
	default:
		return false
	}
}

func checkForConnectedDB(db *gorm.DB, table, checkType string, query interface{}, args ...interface{}) bool {
	if checkType == "notexists" {
		return !CheckExistsInTable(db, table, query, args...)
	} else if checkType == "exists" {
		return CheckExistsInTable(db, table, query, args...)
	} else {
		return false
	}
}

func ValidateNext(value reflect.Value) (interface{}, bool) {
	if value.Type().Kind() == reflect.Int {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int8 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int16 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int32 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Int64 {
		return value.Int(), value.Int() != 0
	} else if value.Type().Kind() == reflect.Uint {
		return value.Uint(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint8 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint16 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint32 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uint64 {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Uintptr {
		return value.Int(), value.Uint() != 0
	} else if value.Type().Kind() == reflect.Float32 {
		return value.Float(), value.Float() != 0
	} else if value.Type().Kind() == reflect.Float64 {
		return value.Float(), value.Float() != 0
	} else if value.Type().Kind() == reflect.Bool {
		return value.Bool(), true
	} else if value.Type().Kind() == reflect.String {
		return value.String(), value.String() != ""
	} else {
		return value.String(), value.String() != ""
	}
}
