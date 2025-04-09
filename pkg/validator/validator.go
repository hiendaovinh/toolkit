package valideitor

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var alphaNumericRegex = regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)

type HasGetOrZero[T any] interface {
	GetOrZero() T
	IsSet() bool
}

type ValidatorStruct interface {
	Struct(s any) error
}

func NewDefaultValidator() (ValidatorStruct, error) {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if v, ok := field.Interface().(HasGetOrZero[string]); ok {
			return v.GetOrZero()
		}

		if v, ok := field.Interface().(HasGetOrZero[int]); ok {
			return v.GetOrZero()
		}

		if v, ok := field.Interface().(HasGetOrZero[uuid.UUID]); ok {
			if !v.IsSet() {
				return nil
			}
			return v.GetOrZero().String()
		}

		return nil
	}, omit.Val[string]{}, omitnull.Val[string]{}, omitnull.Val[uuid.UUID]{}, omit.Val[uuid.UUID]{}, omit.Val[int]{})

	err := v.RegisterValidation("slug", func(fl validator.FieldLevel) bool {
		return alphaNumericRegex.MatchString(fl.Field().String())
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}

func ValidateStruct(c context.Context, v ValidatorStruct, s any) error {
	err := v.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fields := []string{}
	for _, f := range validationErrors {
		fields = append(fields, f.StructField())
	}

	return fmt.Errorf("invalid %s", strings.Join(fields, ", "))
}
