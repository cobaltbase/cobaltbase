package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func GetAllItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		table := chi.URLParam(r, "table")
		var items []js
		config.DB.Table(table).Find(&items)
		render.JSON(w, r, items)
	}
}
