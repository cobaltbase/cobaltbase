package ct

import (
	"github.com/lib/pq"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type ContextKey string

const JsonDataKey = ContextKey("jsondata")
const AuthMiddlewareKey = ContextKey("auth_middleware_key")

type Js = map[string]interface{}

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
	Table  string        `json:"tableName" gorm:"uniqueIndex"`
	Fields []SchemaField `json:"fields" gorm:"foreignKey:SchemaID"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	base.ID, err = gonanoid.New(10) // Generate a 10-character NanoID
	return
}

func (Schema) TableName() string {
	return "_schemas"
}

func (SchemaField) TableName() string {
	return "_schema_fields"
}
