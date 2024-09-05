package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cobaltbase/cobaltbase/internal/config"
	"github.com/cobaltbase/cobaltbase/internal/ct"
	"github.com/cobaltbase/cobaltbase/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lib/pq"
)

func CheckTableExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tableName := chi.URLParam(r, "table")
		schema := utils.Schemas[tableName]
		if schema.Table == "" {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{"message": "Invalid table name"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func PreProcessingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		contentType := r.Header.Get("Content-Type")
		tableName := chi.URLParam(r, "table")

		schema := utils.Schemas[tableName]

		if strings.HasPrefix(contentType, "multipart/form-data") {

			if err := r.ParseMultipartForm(1024 * 1024); err != nil {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{"message": "Invalid Form Data"})
			}

			var data = make(ct.Js)

			for key, value := range r.Form {
				if len(value) == 1 {
					data[key] = value[0]
				} else {
					data[key] = value
				}
			}

			filesForm := r.MultipartForm

			for _, v := range schema.Fields {
				if v.Type == "singleFile" {
					if len(filesForm.File[v.Name]) > 1 {
						render.Status(r, 404)
						render.JSON(w, r, ct.Js{
							"error": fmt.Sprintf("Expected one file got multiple for key '%s'", v.Name),
						})
						return
					}
					var err error
					data[v.Name], err = utils.UploadSingleFileLocally(filesForm.File[v.Name][0])
					if err != nil {
						render.Status(r, 400)
						render.JSON(w, r, ct.Js{
							"error": err.Error(),
						})
						return
					}

				}
				if v.Type == "multipleFiles" {
					var err error
					data[v.Name], err = utils.UploadMultipleFilesLocally(filesForm.File[v.Name])
					if err != nil {
						render.Status(r, 400)
						render.JSON(w, r, ct.Js{
							"error": err.Error(),
						})
						return
					}
				}
			}

			data, errors := ValidataBody(data, schema)
			if len(errors) > 0 {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{
					"message": "validation error",
					"errors":  errors,
				})
				return
			}

			ctx := context.WithValue(r.Context(), ct.JsonDataKey, data)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)

		} else if strings.HasPrefix(contentType, "application/json") {

			var data ct.Js
			render.DecodeJSON(r.Body, &data)
			data, errors := ValidataBody(data, schema)
			if len(errors) > 0 {
				render.Status(r, 400)
				render.JSON(w, r, ct.Js{
					"message": "validation error",
					"errors":  errors,
				})
				return
			}

			ctx := context.WithValue(r.Context(), ct.JsonDataKey, data)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			render.Status(r, 400)
			render.JSON(w, r, ct.Js{
				"message": "invalid request",
			})
			return
		}

	})
}

// func stripOmitempty(tag string) string {
// 	parts := strings.Split(tag, ",")
// 	return parts[0] // Return only the tag name, effectively stripping ,omitempty
// }

func ValidataBody(body ct.Js, schema ct.Schema) (ct.Js, []string) {
	validate := config.Validate
	var errors []string
	for _, v := range schema.Fields {
		value := body[v.Name]

		if v.Required && value == nil {
			errors = append(errors, fmt.Sprintf("required field '%s' is not present", v.Name))
		}

		if (v.Type == "email" || v.Type == "url") && value != nil {
			err := validate.Var(value, v.Type)
			if err != nil {
				errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
			}
		}

		if (v.Type == "integer" || v.Type == "float") && value != nil {
			str, ok := value.(string)
			if ok {
				i, err := strconv.ParseFloat(str, 64)
				if err != nil {
					errors = append(errors, fmt.Sprintf("validation failed for '%s'", v.Name))
				} else {
					body[v.Name] = i
				}
			} else {
				_, ok = value.(float64)
				if !ok {
					errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
				}
			}

		}

		if v.Type == "boolean" {
			str, ok := value.(string)
			if ok {
				if str == "true" {
					body[v.Name] = true
				}
				if str == "false" {
					body[v.Name] = false
				}
				if str == "yes" {
					body[v.Name] = true
				}
				if str == "no" {
					body[v.Name] = false
				}
				if str == "t" {
					body[v.Name] = true
				}
				if str == "f" {
					body[v.Name] = false
				}
			} else {
				_, ok = value.(bool)
				if !ok {
					errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
				}
			}
		}

		if v.Type == "singleSelect" && value != nil {
			str, ok := value.(string)
			if !ok {
				errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
			} else {
				if !contains(v.SelectOptions, str) {
					errors = append(errors, fmt.Sprintf("'%s' is not an option for field '%s'", value, v.Name))
				}
			}
		}

		if v.Type == "multipleSelect" && value != nil {
			var array pq.StringArray
			valueInterface, ok := value.([]interface{})
			if ok {
				for _, val := range valueInterface {
					val, ok := val.(string)
					if !ok {
						errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
					}
					if !contains(v.SelectOptions, val) {
						errors = append(errors, fmt.Sprintf("'%s' is not an option for field '%s'", val, v.Name))
					}
					array = append(array, val)
				}
			} else {
				valueString, ok := value.([]string)
				if ok {
					for _, val := range valueString {
						if !contains(v.SelectOptions, val) {
							errors = append(errors, fmt.Sprintf("'%s' is not an option for field '%s'", val, v.Name))
						}
						array = append(array, val)
					}
				} else {
					errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
				}
			}

			body[v.Name] = array
		}

		if (v.Type == "singleRelation" || v.Type == "singleFile") && value != nil {
			_, ok := value.(string)
			if !ok {
				errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
			}
		}

		if (v.Type == "multipleRelation" || v.Type == "multipleFiles") && value != nil {
			var array pq.StringArray
			valueInterface, ok := value.([]interface{})
			if ok {
				for _, val := range valueInterface {
					val, ok := val.(string)
					if !ok {
						errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
					}
					array = append(array, val)
				}
			} else {
				valueStrings, ok := value.([]string)
				if ok {
					for _, val := range valueStrings {
						array = append(array, val)
					}
				} else {
					errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
				}
			}
			body[v.Name] = array
		}

		if v.Type == "json" && value != nil {
			val, ok := value.(string)
			if ok {
				var jsonMap ct.Js
				err := json.Unmarshal([]byte(val), &jsonMap)
				if err != nil {
					errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
				}
				body[v.Name] = jsonMap

			}
		}

		if v.Type == "string" && value != nil {
			_, ok := value.(string)
			if !ok {
				errors = append(errors, fmt.Sprintf("validation failed for field '%s'", v.Name))
			}
		}

	}
	return body, errors
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
