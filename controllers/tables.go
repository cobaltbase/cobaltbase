package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/config"
	"github.com/cobaltbase/cobaltbase/customTypes"
	"github.com/go-chi/render"
)

type js = customTypes.Json

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

	}
}
