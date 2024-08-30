package customTypes

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type ContextKey string

const JsonDataKey = ContextKey("jsondata")

type Json = map[string]interface{}

type BaseModel struct {
	ID string `gorm:"primarykey"`
	gorm.Model
}

type SchemaField struct {
	BaseModel
	Name          string `json:"name"`
	Type          string `json:"type"`
	Required      bool   `json:"required"`
	Unique        bool   `json:"unique"`
	SchemaID      string
	RelatedTable  string         `json:"relatedTable,omitempty"`
	SelectOptions pq.StringArray `json:"selectOptions,omitempty" gorm:"type:text[]"`
}

type Schema struct {
	BaseModel
	TableName string        `json:"tableName" gorm:"uniqueIndex"`
	Fields    []SchemaField `json:"fields" gorm:"foreignKey:SchemaID"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID, err = gonanoid.New(10) // Generate a 10-character NanoID
	return
}

type ValidatorHook struct{}

func validateStruct(data interface{}) error {
	validate := validator.New()
	// Perform validation
	err := validate.Struct(data)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}

func (v ValidatorHook) BeforeCreate(tx *gorm.DB) (err error) {
	return validateStruct(tx.Statement.Model)
}

func (v ValidatorHook) BeforeUpdate(tx *gorm.DB) (err error) {
	return validateStruct(tx.Statement.Model)
}
