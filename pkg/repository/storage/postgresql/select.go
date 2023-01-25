package postgresql

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func SelectAllFromDb(db *gorm.DB, order string, receiver interface{}, query interface{}, args ...interface{}) error {
	if order == "" {
		order = "desc"
	}
	tx := db.Where(query, args...).Find(receiver)
	return tx.Error
}

func SelectOneFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}
func SelectLatestFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Order("id desc").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func SelectRandomFromDb(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) (error, error) {

	tx := db.Order("rand()").Where(query, args...).First(receiver)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error, tx.Error
	}
	return tx.Error, nil
}

func SelectFirstFromDb(db *gorm.DB, receiver interface{}) error {
	tx := db.First(receiver)
	return tx.Error
}

func CheckExists(db *gorm.DB, receiver interface{}, query interface{}, args ...interface{}) bool {

	tx := db.Where(query, args...).First(receiver)
	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}

func CheckExistsInTable1(db *gorm.DB, table string, query interface{}, args ...interface{}) bool {
	var result interface{}
	tx := db.Table(table).Where(query, args...).Take(&result)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return false
		} else {
			fmt.Println("tx error", tx.Error.Error())
		}
	}

	return true
}

func CheckExistsInTable(db *gorm.DB, table string, query interface{}, args ...interface{}) bool {
	var result map[string]interface{}
	tx := db.Table(table).Where(query, args...).Take(&result)
	return tx.RowsAffected != 0
}
