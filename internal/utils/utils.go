package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func CreateStructFromSchema(schema ct.Schema) interface{} {
	fields := []reflect.StructField{
		{
			Name:      "Model",
			Type:      reflect.TypeOf(ct.BaseModel{}),
			Anonymous: true,
		},
	}

	// Add fields to the struct type
	for _, field := range schema.Fields {
		fields = append(fields, reflect.StructField{
			Name: StringToUpperCamelCase(field.Name),
			Type: getReflectType(field.Type),
			Tag:  generateTag(field),
		})
	}

	newType := reflect.StructOf(fields)

	// Create a new instance of the struct
	return reflect.New(newType).Interface()
}

func CreateSchema(db *gorm.DB, schema ct.Schema) error {
	// Create a new struct type dynamically
	newStruct := CreateStructFromSchema(schema)

	// Create the table
	err := db.Table(schema.TableName).AutoMigrate(newStruct)
	if err != nil {
		return fmt.Errorf("failed to create or update table: %v", err)
	}

	return nil
}

func StringToUpperCamelCase(s string) string {
	var result strings.Builder

	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for _, word := range words {
		if len(word) > 0 {
			// Convert the first letter to uppercase and the rest to lowercase
			result.WriteString(strings.ToUpper(string(word[0])) + strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

func getReflectType(fieldType string) reflect.Type {
	switch fieldType {
	case "email":
		return reflect.TypeOf("")
	case "url":
		return reflect.TypeOf("")
	case "string":
		return reflect.TypeOf("")
	case "integer":
		return reflect.TypeOf(0)
	case "float":
		return reflect.TypeOf(float64(0))
	case "boolean":
		return reflect.TypeOf(false)
	case "json":
		return reflect.TypeOf(datatypes.JSON{})
	case "singleRealtion":
		return reflect.TypeOf("")
	case "multipleRealtion":
		return reflect.TypeOf(pq.StringArray{})
	case "datetime":
		return reflect.TypeOf(time.Time{})
	case "singleFile":
		return reflect.TypeOf("")
	case "multipleFiles":
		return reflect.TypeOf(pq.StringArray{})
	case "singleSelect":
		return reflect.TypeOf("")
	case "multipleSelect":
		return reflect.TypeOf(pq.StringArray{})
	default:
		return reflect.TypeOf("")
	}
}

func generateTag(field ct.SchemaField) reflect.StructTag {
	// Start with the JSON tag.
	tag := fmt.Sprintf(`json:"%s"`, field.Name)

	// Build the GORM tag based on field properties.
	var gormTag string

	gormTag = appendGormOption(gormTag, fmt.Sprintf("column:%s", field.Name))

	// Check if the field type is "multipleRelation" to add the text[] tag.
	if field.Type == "multipleRelation" {
		gormTag = `gorm:"type:text[]"`
	}

	if field.Type == "multipleSelect" {
		gormTag = `gorm:"type:text[]"`
	}

	if field.Type == "multipleFiles" {
		gormTag = `gorm:"type:text[]"`
	}

	// Add unique and not null constraints if applicable.
	if field.Unique {
		gormTag = appendGormOption(gormTag, "uniqueIndex")
	}
	if field.Required {
		gormTag = appendGormOption(gormTag, "not null")
	}

	// Add GORM tag if needed.
	if gormTag != "" {
		tag = fmt.Sprintf(`%s %s`, tag, gormTag)
	}

	if field.Type == "email" {
		tag = fmt.Sprintf(`%s validate:"email"`, tag)
	}

	if field.Type == "url" {
		tag = fmt.Sprintf(`%s validate:"url"`, tag)
	}

	return reflect.StructTag(tag)
}

// appendGormOption appends a new GORM option to the existing tag with proper formatting.
func appendGormOption(tag, option string) string {
	if tag == "" {
		return fmt.Sprintf(`gorm:"%s"`, option)
	}
	// If tag already has gorm options, append the new option with a semicolon.
	return tag[:len(tag)-1] + ";" + option + `"`
}
