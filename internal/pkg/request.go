package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidateBody[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var req T

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Message: "invalid request body"})
		return req, false
	}

	if err := validate.Struct(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, Response{Message: err.Error()})
		return req, false
	}

	return req, true
}

// func DecodeAndValidate[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
// 	var req T

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	if err := decoder.Decode(&req); err != nil {
// 		WriteJSON(w, http.StatusBadRequest, Response{Message: "invalid request body"})
// 		return req, false
// 	}

// 	if err := validate.Struct(req); err != nil {
// 		WriteJSON(w, http.StatusBadRequest, Response{Message: err.Error()})
// 		return req, false
// 	}

// 	return req, true
// }
