package controllers

import (
	"net/http"
	"time"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
		for i := 0; i < 10; i++ {
			render.JSON(w, r, "multi reponse")
			time.Sleep(5 * 1000 * 1000)
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
		// if err := r.ParseMultipartForm(1024 * 1024); err != nil {
		// 	fmt.Println(err.Error())
		// }

		// v := reflect.ValueOf(&ct.SchemaField{}).Elem()
		// t := v.Type()

		// for i := 0; i < t.NumField(); i++ {
		// 	jsontag := t.Field(i).Tag.Get("json")
		// 	if jsontag != "" {
		// 		fmt.Println(stripOmitempty(jsontag))
		// 	}
		// }

		data, ok := r.Context().Value(ct.JsonDataKey).(ct.Js)

		if !ok {
			http.Error(w, "Middleware data not found", http.StatusInternalServerError)
			return
		}

		// data["id"], _ = gonanoid.New(10)
		// data["created_at"] = time.Now()
		// data["updated_at"] = time.Now()

		// err := config.DB.Table(tableName).Create(&data).Error
		// if err != nil {
		// 	render.Status(r, 400)
		// 	render.JSON(w, r, js{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		render.JSON(w, r, js{
			"message": "Create Item Endpoint",
			"data":    data,
		})
	}
}
