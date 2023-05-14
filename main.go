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

func main() {
	// Create MongoDB connection
	const dburi = "mongodb://localhost:27017"
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	app.Listen(":5000")
}
