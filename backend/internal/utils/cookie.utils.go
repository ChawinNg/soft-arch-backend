package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateSessionCookie(value string, expireInSecond int) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "session",
		Value:    value,
		Expires:  time.Now().Add(time.Duration(expireInSecond) * time.Second),
		HTTPOnly: true,
	}
}
