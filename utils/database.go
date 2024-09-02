package utils

import (
	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
)

var OperatorIndex = map[string]bool{
	">":           true,
	"<=":          true,
	">=":          true,
	"<>":          true,
	"BETWEEN":     true,
	"NOT BETWEEN": true,
	"IN":          true,
	"NOT IN":      true,
	"LIKE":        true,
	"NOT LIKE":    true,
}

func CheckForFieldInSchema(tableName string, field string) bool {
	var schema ct.Schema
	result := config.DB.Preload("Fields").First(&schema, &ct.Schema{
		TableName: tableName,
	})
	if result.Error != nil {
		return false
	}

	for _, v := range schema.Fields {
		if v.Name == field {
			return true
		}
	}

	return false
}
