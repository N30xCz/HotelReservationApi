package main

import (
	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	appv1 := app.Group("/api/v1")
	app.Get("/user", userHandler)
	appv1.Get("/users", api.GetUser)

	app.Listen(":5000")
}
func userHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"welcome msg": "Hello From server :)"})

}
