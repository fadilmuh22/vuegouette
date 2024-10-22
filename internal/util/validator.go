package util

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() (*CustomValidator, error) {
	cv := &CustomValidator{}
	err := cv.Init()
	if err != nil {
		return nil, err
	}
	return cv, nil
}

func (cv *CustomValidator) Init() error {
	cv.Validator = validator.New()
	err := cv.Validator.RegisterValidation("uuid", isValidateUUID)
	if err != nil {
		return err
	}

	cv.Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func isValidateUUID(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := uuid.FromString(id)
	return err == nil
}

func IsValidTikTokURL(url string) bool {
	return len(url) > 0 // && strings.Contains(url, "tiktok.com") && strings.Contains(url, "video")
}
