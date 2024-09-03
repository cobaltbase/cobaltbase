package utils

import (
	"encoding/json"
	"strings"

	"github.com/cobaltbase/cobaltbase/ct"
)

var OperatorIndex = map[string]bool{
	"=":           true,
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

func CheckForFieldInSchema(tableName string, field string, schema ct.Schema) bool {

	for _, v := range schema.Fields {
		if v.Name == field {
			return true
		}
	}

	return false
}

func ConvertPgArrayToSlice(pgArray string) []string {
	// Remove the curly braces
	trimmed := strings.Trim(pgArray, "{}")

	// Split the string by comma
	elements := strings.Split(trimmed, ",")

	// Trim whitespace from each element
	for i, elem := range elements {
		elements[i] = strings.TrimSpace(elem)
	}

	return elements
}

func ProcessOutputData(tableName string, data ct.Js, schema ct.Schema) ct.Js {

	for _, v := range schema.Fields {
		if data[v.Name] != nil {
			if v.Type == "multipleSelect" || v.Type == "multipleFiles" || v.Type == "multipleRelations" {
				data[v.Name] = ConvertPgArrayToSlice(data[v.Name].(string))
			}
			if v.Type == "json" {
				var unmarshalledData ct.Js
				err := json.Unmarshal([]byte(data[v.Name].(string)), &unmarshalledData)
				if err != nil {
					data[v.Name] = "invalid json"
				}
				data[v.Name] = unmarshalledData
			}
		}
	}

	return data
}

func ProcessOutputDataList(tableName string, items []ct.Js, schema ct.Schema) []ct.Js {
	for index, item := range items {
		items[index] = ProcessOutputData(tableName, item, schema)
	}
	return items
}
