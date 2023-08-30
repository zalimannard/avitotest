package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var (
	Instance  *validator.Validate
	slugRegex = regexp.MustCompile("^[a-zA-Z0-9_]*$")
)

func init() {
	Instance = validator.New()
	RegisterCustomValidators(Instance)
}

func ValidateSlug(fl validator.FieldLevel) bool {
	slug := fl.Field().String()
	return slugRegex.MatchString(slug)
}

func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("validateslug", ValidateSlug)
}
