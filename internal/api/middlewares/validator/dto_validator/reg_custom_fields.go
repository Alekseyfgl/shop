package dto_validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validate = validator.New()

// init registers custom validation rules.
func init() {
	// Register custom email validation.
	validate.RegisterValidation("custom_email", func(fl validator.FieldLevel) bool {
		emailRegex := `^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		re := regexp.MustCompile(emailRegex)
		return re.MatchString(fl.Field().String())
	})
	//validate.RegisterValidation("img_base64_or_null", imgBase64OrNull)
}
