package middleware

import (
	"fmt"

	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("user is not Admin")
	}
	return c.Next()
}
