package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

type CustomValidator struct {
	v         *validator.Validate
	passwdErr error
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.v.Struct(i)
	if err != nil {
		fieldErr := err.(validator.ValidationErrors)[0]

		return cv.newValidationError(fieldErr.Field(), fieldErr.Value(), fieldErr.Tag(), fieldErr.Param())
	}
	return nil
}

func (cv *CustomValidator) newValidationError(field string, value interface{}, tag string, param string) error {
	switch tag {
	case "required":
		return fmt.Errorf("field %s is required", field)
	case "email":
		return fmt.Errorf("field %s must be a valid email address", field)
	case "password":
		return cv.passwdErr
	case "min":
		return fmt.Errorf("field %s must be at least %s characters", field, param)
	case "max":
		return fmt.Errorf("field %s must be at most %s characters", field, param)
	default:
		return fmt.Errorf("field %s is invalid", field)
	}
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	cv := &CustomValidator{v: v}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	err := v.RegisterValidation("password", cv.passwordValidate)
	if err != nil {
		panic(err)
	}

	return cv
}

func (cv *CustomValidator) passwordValidate(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	lower := regexp.MustCompile(`[a-z]`)
	upper := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`\d`)
	length := regexp.MustCompile(`.{8,}`)

	return lower.MatchString(password) && upper.MatchString(password) && digit.MatchString(password) && length.MatchString(password)
}
