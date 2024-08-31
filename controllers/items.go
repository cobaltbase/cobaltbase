package controllers

import (
	"net/http"
	"time"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/cobaltbase/cobaltbase/middlewares"
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

		data, ok := r.Context().Value(ct.JsonDataKey).(middlewares.Person)

		if !ok {
			http.Error(w, "Middleware data not found", http.StatusInternalServerError)
			return
		}

		render.JSON(w, r, js{
			"message": "Create Item Endpoint",
			"data":    data,
		})
	}
}
