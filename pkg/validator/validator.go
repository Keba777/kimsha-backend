package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func Validate(s interface{}) map[string]string {
	err := v.Struct(s)
	if err == nil {
		return nil
	}
	errs := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		errs[field] = fmt.Sprintf("failed on '%s' rule", e.Tag())
	}
	return errs
}
