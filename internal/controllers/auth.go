package controllers

import (
	"net/http"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user ct.Auth
		var requestBody ct.AuthRequest

		if err := render.DecodeJSON(r.Body, &requestBody); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
		}

		if err := config.Validate.Var(requestBody.Email, `email,required`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
		}

		if err := config.Validate.Var(requestBody.Email, `required,min=8`); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "invalid body"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		user.Email = requestBody.Email
		user.Password = string(hashedPassword)

		if err := config.DB.Create(&user).Error; err != nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		user.Password = ""
		render.JSON(w, r, ct.Js{"message": "User Created Succesfully", "user": user})
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
