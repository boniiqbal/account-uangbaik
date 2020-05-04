package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/go-playground/validator.v9"
)

// Message as a response function standard
func Message(status bool, message string, data interface{}, includes interface{}) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message, "data": data, "includes": includes}
}

// Respond as a response format standard
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// RoundNumber for rounding of amount
func RoundNumber(x int64) int64 {
	var y int64
	if x == 0 {
		y = 1
	} else {
		y = x
	}
	return y
}

//TokenGenerator random Token
func TokenGenerator(phone string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(phone+time.Now().String()), bcrypt.MinCost)
	return string(b)
}

// NewValidatorError validate error
func NewValidatorError(err error) map[string]interface{} {
	var errors string

	for _, v := range err.(validator.ValidationErrors) {
		switch v.Tag() {
		case "min":
			errors = v.Field() + " Minimum " + v.Param()
		case "max":
			errors = v.Field() + " Maximum " + v.Param()
		case "email":
			errors = v.Field() + " must be a valid email "
		case "required":
			errors = v.Field() + " cannot be empty "
		}
	}

  return Message(false, errors, nil, nil)
}
