package postgresql

import (
	"gorm.io/gorm"
)

func SaveAllFields(db *gorm.DB, model interface{}) (*gorm.DB, error) {
	result := db.Save(model)
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}
