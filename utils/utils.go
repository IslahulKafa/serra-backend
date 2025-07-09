package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	var wrapped map[string]any

	switch status {
	case http.StatusOK, http.StatusCreated:
		wrapped = map[string]any{
			"status": "success",
			"data":   payload,
		}
	default:
		wrapped = map[string]any{
			"status":  "error",
			"message": payload,
		}
	}

	return json.NewEncoder(w).Encode(wrapped)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, err.Error())
}
