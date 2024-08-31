package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/cobaltbase/cobaltbase/utils"
	"github.com/go-chi/render"
)

type js = ct.Json

func GetAllTables() http.HandlerFunc {
	//var db = config.DB
	return func(w http.ResponseWriter, r *http.Request) {

		var tables []string
		config.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname='public';").Scan(&tables)

		render.JSON(w, r, tables)

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
		//config.DB.Create(&schema)

		err := utils.CreateSchema(config.DB, schema)
		if err != nil {
			render.Status(r, 500)
			render.JSON(w, r, js{"message": "Schema Created Successfully", "error": err.Error()})
			return
		}

		render.Status(r, 201)
		render.JSON(w, r, js{"message": "Schema Created Successfully", "table": schema.TableName})
	}
}
