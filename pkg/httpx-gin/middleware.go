package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/auth"
	"github.com/hiendaovinh/toolkit/pkg/errorx"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/unrolled/secure"
)

type Guard interface {
	AuthenticateJWT(ctx context.Context, tokenStr string) (*jwt.Token, *jwtx.JWTClaims, error)
}

func Authn(guard Guard) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.Next()
			return
		}

		parts := strings.Split(header, "Bearer")
		if len(parts) != 2 {
			c.Next()
			return
		}

		token := strings.TrimSpace(parts[1])
		if len(token) == 0 {
			c.Next()
			return
		}

		_, claims, err := guard.AuthenticateJWT(c, token)
		if err != nil {
			// although it's a client error, we don't want to detailed information
			Abort(c, errorx.Wrap(errors.New("invalid access token"), errorx.Authn), -1)
			return
		}

		ctx := c.Request.Context()
		ctx = auth.WithAuthJWT(ctx, token)
		ctx = auth.WithAuthClaims(ctx, claims)
		c.Request = c.Request.WithContext(ctx) // we don't want to use gin.Context.Set here for universal context.Context usages
		c.Next()
	}
}

type CaptchaPayload struct {
	Captcha         string `json:"captcha"`
	CaptchaFallback string `json:"captcha_fallback"`
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

func CaptchaValid(requirement CaptchaRequirement) gin.HandlerFunc {
	return func(c *gin.Context) {
		if requirement.Action == "" {
			Abort(c, errorx.Wrap(errors.New("missing captcha action"), errorx.Service))
			return
		}

		v := c.Request.Header.Get("x-captcha-internal")
		if requirement.Bypass != "" && requirement.Bypass == v {
			c.Next()
			return
		}

		if requirement.Secret == "" {
			Abort(c, errorx.Wrap(errors.New("missing captcha secret"), errorx.Service))
			return
		}

		var payload CaptchaPayload
		if err := c.ShouldBindBodyWith(&payload, binding.JSON); err != nil {
			Abort(c, errorx.Wrap(err, errorx.Invalid))
			return
		}

		if payload.CaptchaFallback != "" {
			resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
				"secret":   {requirement.FallbackSecret},
				"response": {payload.CaptchaFallback},
			})
			if err != nil {
				Abort(c, errorx.Wrap(err, errorx.Service))
				return
			}
			defer resp.Body.Close()

			var body CaptchaVerifyResponse
			if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
				Abort(c, errorx.Wrap(err, errorx.Service))
				return
			}

			if !body.Success {
				Abort(c, errorx.Wrap(errors.New("fallback captcha failed"), errorx.Captcha))
				return
			}

			// TODO check body.ChallengeTS
			c.Next()
			return
		}

		resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify", url.Values{
			"secret":   {requirement.Secret},
			"response": {payload.Captcha},
		})
		if err != nil {
			Abort(c, errorx.Wrap(err, errorx.Service))
			return
		}
		defer resp.Body.Close()

		var body CaptchaVerifyResponse
		if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
			Abort(c, errorx.Wrap(err, errorx.Service))
			return
		}

		if body.Action == "" || (requirement.Action != "" && body.Action != requirement.Action) {
			Abort(c, errorx.Wrap(errors.New("captcha failed"), errorx.Captcha))
			return
		}

		if body.Score <= 0 || !body.Success || body.Score < requirement.Score {
			Abort(c, errorx.Wrap(errors.New("captcha failed"), errorx.Captcha))
			return
		}

		// TODO check body.ChallengeTS
		c.Next()
	}
}

type ctxKey string

const (
	ctxKeyCSPNonce ctxKey = "csp-nonce"
)

func CSP(secureMiddleware *secure.Secure) gin.HandlerFunc {
	return func(c *gin.Context) {
		nonce, err := secureMiddleware.ProcessAndReturnNonce(c.Writer, c.Request)
		// If there was an error, do not continue.
		if err != nil {
			c.Abort()
			return
		}

		// Avoid header rewrite if response is a redirection.
		if status := c.Writer.Status(); status > 300 && status < 399 {
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), ctxKeyCSPNonce, nonce))
		c.Next()
	}
}

func CSPNonce(ctx context.Context) any {
	return ctx.Value(ctxKeyCSPNonce)
}
