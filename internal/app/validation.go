package app

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func ConfigureValidation() {
	validate = validator.New()
	validate.RegisterValidation("not_future", dateIsNotInFuture)
}

func dateIsNotInFuture(fl validator.FieldLevel) bool {
	switch v := fl.Field().Interface().(type) {
	case *time.Time:
		{
			return v == nil || !v.After(time.Now())
		}
	case time.Time:
		{
			return !v.After(time.Now())
		}
	case string:
		{
			date, _ := time.Parse("2006-01-02", v)
			return !date.After(time.Now())
		}
	default:
		{
			panic("Checked value's type is not supported")
		}
	}
}
