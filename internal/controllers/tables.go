package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type js = ct.Js

func GetAllTables() http.HandlerFunc {
	//var db = config.DB
	return func(w http.ResponseWriter, r *http.Request) {
		var tables = make([]string, len(utils.Schemas))

		i := 0
		for key := range utils.Schemas {
			tables[i] = key
			i++
		}

		render.JSON(w, r, tables)

	}
}

func GetSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")
		schema := utils.Schemas[tableName]
		render.JSON(w, r, schema)
	}
}

func CreateTable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var schema ct.Schema
		if err := render.DecodeJSON(r.Body, &schema); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   err.Error(),
			})
		}
		result := config.DB.Save(&schema)
		if result.Error != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   result.Error.Error(),
			})
			return
		}

		config.UpdateAndMigrateSchemas()
		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Schema Created Successfully", "table": schema.TableName})
	}
}

func UpdateSingleField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var field ct.SchemaField
		if err := render.DecodeJSON(r.Body, &field); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   err.Error(),
			})
			return
		}
		err := config.DB.Save(&field).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   err.Error(),
			})
			return
		}

		config.UpdateAndMigrateSchemas()
		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Schema Created Field", "field": field})
	}

}

func DeleteSingleField() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var field ct.SchemaField
		if err := render.DecodeJSON(r.Body, &field); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   err.Error(),
			})
			return
		}
		err := config.DB.Delete(&field).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"message": "invalid input",
				"error":   err.Error(),
			})
			return
		}
		config.UpdateAndMigrateSchemas()
		render.Status(r, 200)
		render.JSON(w, r, js{"message": "Field Deleted Successfully"})
	}
}

func DeleteSchemaWithoutData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")
		err := config.DB.Delete(&ct.SchemaField{}, &ct.SchemaField{
			SchemaID: utils.Schemas[tableName].ID,
		}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"error": err.Error(),
			})
		}
		err = config.DB.Delete(&ct.Schema{}, &ct.Schema{
			TableName: tableName,
		}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"error": err.Error(),
			})
		}
		config.UpdateAndMigrateSchemas()
		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Schema Deleted Successfully"})
	}
}

// Danger
func DeleteSchemaWithData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")
		soft := r.URL.Query().Get("soft")
		var err error
		if soft == "no" {

			err = config.DB.Unscoped().Delete(&ct.SchemaField{}, &ct.SchemaField{
				SchemaID: utils.Schemas[tableName].ID,
			}).Error
			if err != nil {
				render.Status(r, 400)
				render.JSON(w, r, js{
					"error": err.Error(),
				})
				return
			}

			err = config.DB.Unscoped().Delete(&ct.Schema{}, &ct.Schema{
				TableName: tableName,
			}).Error

		} else {

			err = config.DB.Delete(&ct.SchemaField{}, &ct.SchemaField{
				SchemaID: utils.Schemas[tableName].ID,
			}).Error
			if err != nil {
				render.Status(r, 400)
				render.JSON(w, r, js{
					"error": err.Error(),
				})
				return
			}

			err = config.DB.Delete(&ct.Schema{}, &ct.Schema{
				TableName: tableName,
			}).Error

		}

		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"error": err.Error(),
			})
			return
		}
		if soft == "no" {
			err = config.DB.Unscoped().Table(tableName).Where("1=1").Delete(ct.Js{}).Error
		} else {
			err = config.DB.Table(tableName).Where("1=1").Delete(ct.Js{}).Error
		}
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"error": err.Error(),
			})
			return
		}
		config.UpdateAndMigrateSchemas()
		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Schema With Data Deleted"})
	}
}
