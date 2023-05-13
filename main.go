package main

import (
	"context"
	"fmt"
	"log"

	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func userHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"welcome msg": "Hello From server :)"})

}

const dbUri = "mongodb://localhost:27017"
const dbName = "Hotel-Reservation"
const userColl = "users"

func main() {
	// Create MongoDB connection

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	// Insert User to the MongoDB
	ctx := context.Background()
	coll := client.Database(dbName).Collection(userColl)
	user := types.User{
		FirstName: "Pepek",
		LastName:  "Namornik",
	}
	res, err := coll.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	// Start Fiber framework
	app := fiber.New()
	appv1 := app.Group("/api/v1")
	app.Get("/user", userHandler)
	appv1.Get("/users", api.GetUser)

	app.Listen(":5000")

}
