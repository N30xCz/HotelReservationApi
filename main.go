package main

import (
	"context"
	"log"

	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/N30xCz/HotelReservationApi/api/middleware"
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
	// handlers initialization
	var (
		userHandler = api.NewUserHandler(db.NewMongoUserStore(client))
		hotelStore  = db.NewHotelStore(client)
		roomStore   = db.NewRoomStore(client, hotelStore)
		UserStore   = db.NewMongoUserStore(client)
		store       = &db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  UserStore,
		}
		hotelHandler = api.NewHotelHandler(store)
		AuthHandler  = api.NewAuthHandler(UserStore)
		app          = fiber.New(config)
		apiv1        = app.Group("/api/v1", middleware.JWTAuthentication)
		auth         = app.Group("/api")
	)
	// Auth
	auth.Post("/auth", AuthHandler.HandleAuthenticate)
	// Versioned API
	// User Handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	// Hotel Handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	app.Listen(":5000")
}
