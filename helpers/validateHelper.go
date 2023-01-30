package helpers

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	validate.RegisterValidation("dateformat", dateFormat)
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Error()
			errors = append(errors, &element)
		}
	}
	return errors
}

func dateFormat(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, date)
	if match {
		_, err := time.Parse("2006-01-02", date)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
