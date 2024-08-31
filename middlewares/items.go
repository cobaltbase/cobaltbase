package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cobaltbase/cobaltbase/ct"
	"github.com/go-chi/render"
)

func PreProcessingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024); err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Json{"message": "Invalid Form Data"})
		}

		var jsondata = make(map[string]interface{})

		for key, value := range r.Form {
			if len(value) == 1 {
				jsondata[key] = convertToType(value[0])
			} else {
				jsondata[key] = value
			}
		}

		jsonDataJson, err := json.Marshal(jsondata)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Json{"error": "Invalid Json"})
		}

		var person Person

		// Convert JSON to struct
		err = json.Unmarshal(jsonDataJson, &person)
		if err != nil {
			render.Status(r, 400)
			render.JSON(w, r, ct.Json{"error": err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), ct.JsonDataKey, person)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// func stripOmitempty(tag string) string {
// 	parts := strings.Split(tag, ",")
// 	return parts[0] // Return only the tag name, effectively stripping ,omitempty
// }

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Gay  bool   `json:"gay"`
}

func convertToType(s string) interface{} {

	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}

	// Try to convert to int
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}

	// If not int, try to convert to float64
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	return s
}
