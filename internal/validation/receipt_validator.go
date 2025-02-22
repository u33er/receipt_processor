package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
)

var validate *validator.Validate

func dateFormatValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse(time.DateOnly, fl.Field().String())
	return err == nil
}

func timeFormatValidator(fl validator.FieldLevel) bool {
	_, err := time.Parse("15:04", fl.Field().String())
	return err == nil
}

func retailerValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[\w\s\-&]+$`)
	return re.MatchString(fl.Field().String())
}

func totalValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^\d+\.\d{2}$`)
	return re.MatchString(fl.Field().String())
}

func shortDescValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^[\w\s\-]+$`)
	return re.MatchString(fl.Field().String())
}

func priceValidator(fl validator.FieldLevel) bool {
	re := regexp.MustCompile(`^\d+\.\d{2}$`)
	return re.MatchString(fl.Field().String())
}

func notBlankValidator(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

func init() {
	validate = validator.New()
	validations := []struct {
		tag string
		fn  validator.Func
	}{
		{"date", dateFormatValidator},
		{"time", timeFormatValidator},
		{"retailer", retailerValidator},
		{"notblank", notBlankValidator},
		{"total", totalValidator},
		{"shortDesc", shortDescValidator},
		{"price", priceValidator},
	}
	for _, v := range validations {
		if err := validate.RegisterValidation(v.tag, v.fn); err != nil {
			panic(fmt.Errorf("failed to register validation: %w", err))
		}
	}
}

func ValidateReceipt(r interface{}) error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	return nil
}
