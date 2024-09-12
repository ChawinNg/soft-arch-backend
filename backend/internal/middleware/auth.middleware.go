package middleware

import (
	"net/http"

	"backend/internal/model"
	"backend/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Middleware struct {
	secret string
}

func NewMiddleware(secret string) Middleware {
	return Middleware{
		secret: secret,
	}
}

func (m *Middleware) SessionMiddleware(c *fiber.Ctx) error {
	cookie := new(model.Cookies)

	err := c.CookieParser(cookie)
	if err != nil {
		return c.Next()
	}

	claim, err := utils.ParseToken(cookie.Session, m.secret)
	if err == nil {
		c.Locals("session", claim.Session)
	}

	return c.Next()
}

func (m *Middleware) WithAuthentication(next func(*fiber.Ctx) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		_, ok := c.Locals("session").(model.Sessions)
		if !ok {
			return c.Status(http.StatusUnauthorized).SendString("Unauthorized")
		}

		return next(c)
	}
}
