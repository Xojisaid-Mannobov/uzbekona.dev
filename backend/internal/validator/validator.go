// Package validator kiruvchi ma'lumotlarni struct teglari asosida tekshiradi
// va xatolarni o'zbekcha, maydon nomi bo'yicha qaytaradi.
package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"uzbekona.dev/backend/internal/apperr"
	"uzbekona.dev/backend/pkg/slug"
)

type Validator struct {
	v *validator.Validate
}

func New() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())

	// Xato kalitlari JSON nomlari bilan bir xil bo'lsin (masalan "short_description")
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})

	_ = v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return slug.Valid.MatchString(fl.Field().String())
	})

	return &Validator{v: v}
}

// Struct tekshiradi; xato bo'lsa *apperr.Error (422) qaytaradi.
func (val *Validator) Struct(s any) error {
	err := val.v.Struct(s)
	if err == nil {
		return nil
	}

	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return apperr.BadRequest("Noto‘g‘ri ma’lumot")
	}

	fields := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		key := fieldKey(fe.Namespace())
		if _, exists := fields[key]; !exists {
			fields[key] = message(fe)
		}
	}
	return apperr.Validation(fields)
}

// "ProjectInput.metrics[0].value" → "metrics.0.value"
func fieldKey(ns string) string {
	_, rest, found := strings.Cut(ns, ".")
	if !found {
		rest = ns
	}
	rest = strings.ReplaceAll(rest, "[", ".")
	return strings.ReplaceAll(rest, "]", "")
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required", "required_without":
		return "Majburiy maydon"
	case "email":
		return "Email noto‘g‘ri"
	case "url":
		return "Havola noto‘g‘ri (https://… bilan boshlang)"
	case "slug":
		return "Faqat kichik lotin harflari, raqam va “-” belgisi"
	case "hexcolor":
		return "Rang HEX formatida bo‘lsin (#2350F5)"
	case "oneof":
		return "Ruxsat etilmagan qiymat"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Kamida %s ta belgi", fe.Param())
		}
		if fe.Kind() == reflect.Slice {
			return fmt.Sprintf("Kamida %s ta element", fe.Param())
		}
		return fmt.Sprintf("Qiymat %s dan kichik bo‘lmasin", fe.Param())
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("Ko‘pi bilan %s ta belgi", fe.Param())
		}
		if fe.Kind() == reflect.Slice {
			return fmt.Sprintf("Ko‘pi bilan %s ta element", fe.Param())
		}
		return fmt.Sprintf("Qiymat %s dan katta bo‘lmasin", fe.Param())
	default:
		return "Noto‘g‘ri qiymat"
	}
}
