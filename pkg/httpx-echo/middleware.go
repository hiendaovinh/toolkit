package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/auth"
	"github.com/hiendaovinh/toolkit/pkg/errorx"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/hiendaovinh/toolkit/pkg/limiter"
	"github.com/labstack/echo/v4"
)

type Guard interface {
	AuthenticateJWT(ctx context.Context, tokenStr string) (*jwt.Token, *jwtx.JWTClaims, error)
}

func Authn(guard Guard) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return next(c)
			}

			parts := strings.Split(header, "Bearer")
			if len(parts) != 2 {
				return next(c)
			}

			token := strings.TrimSpace(parts[1])
			if len(token) == 0 {
				return next(c)
			}

			_, claims, err := guard.AuthenticateJWT(c.Request().Context(), token)
			if err != nil {
				// although it's a client error, we don't want to detailed information
				//nolint:errcheck
				Abort(c, errorx.Wrap(errors.New("invalid access token"), errorx.Authn), -1)
				return nil
			}

			ctx := c.Request().Context()
			ctx = auth.WithAuthJWT(ctx, token)
			ctx = auth.WithAuthClaims(ctx, claims)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

type CaptchaVerifyResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type CaptchaRequirement struct {
	Secret         string
	Score          float64
	Action         string
	FallbackSecret string
	Bypass         string
}

type CaptchaPayload struct {
	Captcha         string `json:"captcha"`
	CaptchaFallback string `json:"captcha_fallback"`
}

func CaptchaValid(requirement CaptchaRequirement) echo.MiddlewareFunc {
	recaptchaV3 := func(next echo.HandlerFunc, c echo.Context, payload CaptchaPayload) error {
		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
			"secret":   {requirement.Secret},
			"response": {payload.Captcha},
		})
		if err != nil {
			return Abort(c, errorx.Wrap(err, errorx.Service))
		}
		defer resp.Body.Close()

		var body CaptchaVerifyResponse
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return Abort(c, errorx.Wrap(err, errorx.Service))
		}

		if requirement.Action != "" && body.Action != requirement.Action {
			return Abort(c, errorx.Wrap(errors.New("captcha failed"), errorx.Authz))
		}

		if !body.Success || body.Score < requirement.Score {
			return Abort(c, errorx.Wrap(errors.New("captcha failed"), errorx.Authz))
		}

		return next(c)
	}

	recaptchaV2 := func(next echo.HandlerFunc, c echo.Context, payload CaptchaPayload) error {
		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
			"secret":   {requirement.FallbackSecret},
			"response": {payload.CaptchaFallback},
		})
		if err != nil {
			return Abort(c, errorx.Wrap(err, errorx.Service))
		}
		defer resp.Body.Close()

		var body CaptchaVerifyResponse
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			return Abort(c, errorx.Wrap(err, errorx.Service))
		}

		if !body.Success {
			return Abort(c, errorx.Wrap(errors.New("fallback captcha failed"), errorx.Captcha))
		}

		return next(c)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if requirement.Secret == "" {
				return Abort(c, errorx.Wrap(errors.New("missing captcha secret"), errorx.Service))
			}

			var buf bytes.Buffer
			tee := io.TeeReader(c.Request().Body, &buf)

			c.Request().Body = io.NopCloser(tee)

			var payload CaptchaPayload
			if err := c.Bind(&payload); err != nil {
				return Abort(c, errorx.Wrap(err, errorx.Invalid))
			}

			c.Request().Body = io.NopCloser(&buf)

			if payload.CaptchaFallback != "" {
				return recaptchaV2(next, c, payload)
			}

			return recaptchaV3(next, c, payload)
		}
	}
}

func DisableLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = limiter.Skip(ctx)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
