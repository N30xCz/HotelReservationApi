package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	app.Get("/user", userHandler)
	app.Listen(":5000")
}
func userHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"welcome msg": "Hello From server :)"})

}
