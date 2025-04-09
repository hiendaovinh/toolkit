package httpx

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/hiendaovinh/toolkit/pkg/auth"
	"github.com/hiendaovinh/toolkit/pkg/errorx"
	"github.com/hiendaovinh/toolkit/pkg/limiter"
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

func RestAbort(c echo.Context, v any, err error) error {
	if errors.Is(err, auth.ErrInvalidSession) {
		return Abort(c, errorx.Wrap(err, errorx.Authn))
	}

	if errors.Is(err, limiter.ErrRateLimited) {
		return Abort(c, errorx.Wrap(err, errorx.RateLimiting))
	}

	if _, ok := err.(*errorx.Error); ok {
		return Abort(c, err)
	}

	if err != nil {
		return Abort(c, errorx.Wrap(err, errorx.Service))
	}

	return Abort(c, v)
}

func QueryParamInt(c echo.Context, name string, val int) int {
	v := c.QueryParam(name)
	if v == "" {
		return val
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return val
	}

	return i
}
