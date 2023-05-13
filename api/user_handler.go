package api

import (
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Anthony",
		LastName:  "GG",
	}
	return c.JSON(u)
}
