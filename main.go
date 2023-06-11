package main

import (
	"context"
	"log"

	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	// Create MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client, db.DBNAME))
	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	app.Listen(":5000")
}
