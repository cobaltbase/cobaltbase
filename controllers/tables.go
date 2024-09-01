package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/cobaltbase/cobaltbase/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type js = ct.Js

func GetAllTables() http.HandlerFunc {
	//var db = config.DB
	return func(w http.ResponseWriter, r *http.Request) {

		var tables []string
		config.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public';").Scan(&tables)

		render.JSON(w, r, tables)

	}
}

func GetSchema() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")

		var schema ct.Schema

		err := config.DB.Preload("Fields").First(&schema, ct.Schema{
			TableName: tableName,
		}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{"error": err.Error()})
			return
		}
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

		err := utils.CreateSchema(config.DB, schema)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, js{"error": err.Error()})
			return
		}

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

		tableName := chi.URLParam(r, "table")

		var schema ct.Schema

		err = config.DB.Preload("Fields").First(&schema, ct.Schema{
			TableName: tableName,
		}).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{"error": err.Error()})
			return
		}

		dynamicStruct := utils.CreateStructFromSchema(schema)
		err = config.DB.Table(tableName).AutoMigrate(dynamicStruct)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{"error": err.Error()})
			return
		}

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
		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Field Deleted Successfully"})
	}
}
