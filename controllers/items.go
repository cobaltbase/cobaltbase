package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/cobaltbase/cobaltbase/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GetAllItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		var items []js
		result := config.DB.Table(table).Find(&items)
		if result.Error != nil {
			render.JSON(w, r, js{"error": result.Error.Error()})
			return
		}
		render.JSON(w, r, items)
	}
}

func CreateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")
		var schema ct.Schema

		result := config.DB.Preload("Fields").First(&schema, &ct.Schema{
			TableName: tableName,
		})
		if result.Error != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{
				"error": result.Error.Error(),
			})
			return
		}

		data, ok := r.Context().Value(ct.JsonDataKey).(ct.Js)

		if !ok {
			http.Error(w, "Middleware data not found", http.StatusInternalServerError)
			return
		}

		data["id"], _ = gonanoid.New(10)
		data["created_at"] = time.Now()
		data["updated_at"] = time.Now()

		err := config.DB.Table(tableName).Create(&data).Error
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, js{
				"error": err.Error(),
			})
			return
		}

		render.JSON(w, r, js{
			"message": "Create Item Endpoint",
			"data":    data,
		})
	}
}

func GetItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		filterObjString := r.URL.Query().Get("filterObj")
		filterQueryString := r.URL.Query().Get("filterQuery")

		var filterObj ct.Js
		if filterObjString != "" {
			if err := json.Unmarshal([]byte(filterObjString), &filterObj); err != nil {
				render.JSON(w, r, js{"error": err.Error()})
				return
			}
		}

		var filterQueryObj struct {
			Field    string   `json:"field"`
			Operator string   `json:"operator"`
			Value1   string   `json:"value1"`
			Value2   string   `json:"value2"`
			List     []string `json:"list"`
		}
		if filterQueryString != "" {
			if err := json.Unmarshal([]byte(filterQueryString), &filterQueryObj); err != nil {
				render.JSON(w, r, js{"error": err.Error()})
				return
			}
		}

		if !utils.CheckForFieldInSchema(table, filterQueryObj.Field) || !utils.OperatorIndex[filterQueryObj.Operator] {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{
				"message": "invalid query",
			})
		}

		filterQuery := fmt.Sprintf("%s %s ?", filterQueryObj.Field, filterQueryObj.Operator)

		var item js
		result := config.DB.Table(table).Where(filterQuery, filterQueryObj.Value1).Find(&item, filterObj)
		if result.Error != nil {
			render.JSON(w, r, js{"error": result.Error.Error()})
			return
		}
		render.JSON(w, r, item)

	}
}

func UpdateItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
