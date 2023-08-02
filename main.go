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
		userHandler  = api.NewUserHandler(db.NewMongoUserStore(client))
		hotelStore   = db.NewHotelStore(client)
		roomStore    = db.NewRoomStore(client, hotelStore)
		UserStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    UserStore,
			Booking: bookingStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		AuthHandler    = api.NewAuthHandler(UserStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		apiv1          = app.Group("/api/v1", middleware.JWTAuthentication(UserStore))
		auth           = app.Group("/api")
		admin          = apiv1.Group("/admin", middleware.AdminAuth)
	)
	// Auth
	auth.Post("/auth", AuthHandler.HandleAuthenticate)
	// Admin route

	// Versioned API
	// User Handlers
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/users", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	// Hotel Handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	// Room Handlers
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/rooms", roomHandler.HandleGetRooms)
	// Book Handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	app.Listen(":5000")
}
