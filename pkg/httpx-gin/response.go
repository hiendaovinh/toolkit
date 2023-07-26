package httpx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hiendaovinh/toolkit/pkg/errorx"
)

type body struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Abort(c *gin.Context, v any, codes ...int) {
	code := -1
	if len(codes) >= 1 {
		code = codes[0]
	}

	err, ok := v.(error)
	if ok {
		abortErrorWithStatusJSON(c, err, code)
		return
	}

	c.AbortWithStatusJSON(code, &body{Data: v})
}

func abortErrorWithStatusJSON(c *gin.Context, err error, code int) {
	var target *errorx.Error

	message := errorx.MaskErrorMessage(err)

	if !errors.As(err, &target) {
		//nolint:errcheck
		c.Error(err)
		if code == -1 {
			code = http.StatusInternalServerError
		}
		c.AbortWithStatusJSON(code, &body{Code: "error", Message: message})
		return
	}

	if code == -1 {
		code = target.Status()
	}

	if target.Of(errorx.Database) || target.Of(errorx.Service) {
		//nolint:errcheck
		c.Error(err)
		c.AbortWithStatusJSON(code, &body{Code: target.Code(), Message: message})
		return
	}

	c.AbortWithStatusJSON(code, &body{Code: target.Code(), Message: message})
}

type Validator interface {
	Struct(s any) error
	Var(v any, tag string) error
}

func ValidateStruct(c context.Context, checker Validator, s any) error {
	err := checker.Struct(s)
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

func ValidateVar(c context.Context, checker Validator, v any, tag string) error {
	return checker.Var(v, tag)
}
