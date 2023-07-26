package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hiendaovinh/toolkit/pkg/errorx"
	"github.com/labstack/echo/v4"
)

type body struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Abort(c echo.Context, v any, codes ...int) error {
	code := -1
	if len(codes) >= 1 {
		code = codes[0]
	}

	err, ok := v.(error)
	if ok {
		return abortErrorWithStatusJSON(c, err, code)
	}

	if code == -1 {
		code = 200
	}
	return c.JSON(code, &body{Data: v})
}

func abortErrorWithStatusJSON(c echo.Context, err error, code int) error {
	var target *errorx.Error

	message := errorx.MaskErrorMessage(err)

	if !errors.As(err, &target) {
		c.Logger().Error(err)
		if code == -1 {
			code = http.StatusInternalServerError
		}
		return c.JSON(code, &body{Code: "error", Message: message})
	}

	if code == -1 {
		code = target.Status()
	}

	if target.Of(errorx.Database) || target.Of(errorx.Service) {
		c.Logger().Error(err)
		return c.JSON(code, &body{Code: target.Code(), Message: message})
	}

	return c.JSON(code, &body{Code: target.Code(), Message: message})
}

type ValidatorStruct interface {
	Struct(s any) error
}

func ValidateStruct(c echo.Context, v ValidatorStruct, s any) error {
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
