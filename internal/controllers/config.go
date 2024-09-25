package controllers

import (
	"fmt"
	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
)

func ListOauthConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var configsindb []ct.OauthConfig

		if err := config.DB.Find(&configsindb).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"configs": configsindb})
	}
}

func CreateOauthConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqconfig ct.OauthConfig
		err := render.DecodeJSON(r.Body, &reqconfig)
		if err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		if err := config.DB.Create(&reqconfig).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"message": fmt.Sprintf("config saved for %s", reqconfig.Provider)})
	}
}

func RetrieveOauthConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		var reqconfig struct {
			Provider string `json:"provider"`
		}
		reqconfig.Provider = provider
		var configindb ct.OauthConfig
		if err := config.DB.First(&configindb, reqconfig).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"config": configindb})
	}
}

func UpdateOauthConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqconfig ct.OauthConfig
		err := render.DecodeJSON(r.Body, &reqconfig)
		if err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		if err := config.DB.Save(&reqconfig).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"message": fmt.Sprintf("config saved for %s", reqconfig.Provider)})
	}
}

func RemoveOauthConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqconfig struct {
			Provider string `json:"provider"`
		}
		err := render.DecodeJSON(r.Body, &reqconfig)
		if err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		if err := config.DB.Unscoped().Delete(&ct.OauthConfig{}, reqconfig).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"message": fmt.Sprintf("config deleted for %s", reqconfig.Provider)})
	}
}

func UpdateSMTPConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqconfig ct.SMTPConfig
		err := render.DecodeJSON(r.Body, &reqconfig)
		if err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		var savedConfig ct.SMTPConfig
		if err := config.DB.First(&savedConfig).Error; err != nil {
			savedConfig = reqconfig
		}
		savedConfig.Host = reqconfig.Host
		savedConfig.Port = reqconfig.Port
		savedConfig.Username = reqconfig.Username
		savedConfig.Password = reqconfig.Password
		savedConfig.From = reqconfig.From
		savedConfig.FromName = reqconfig.FromName
		if err := config.DB.Save(&savedConfig).Error; err != nil {
			render.JSON(w, r, ct.Js{"error": err.Error()})
			return
		}
		render.JSON(w, r, ct.Js{"message": "config saved for " + savedConfig.From})
	}
}
