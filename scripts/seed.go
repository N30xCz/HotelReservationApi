package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/db/fixtures"
	"github.com/N30xCz/HotelReservationApi/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func SeedUser(fname, lname, email, password string, isAdmin bool) *types.User {

	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     email,
		FirstName: fname,
		LastName:  lname,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s --> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}

func seedHotel(name string, location string, raiting int) *types.Hotel {

	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Raiting:  raiting,
	}
	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel

}
func seedRoom(size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := roomStore.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}
func seedBooking(userID, roomID primitive.ObjectID, fDate, tillDate time.Time, nPersons int, canceled bool) *types.Booking {
	booking := &types.Booking{
		UserID:     userID,
		RoomID:     roomID,
		FromDate:   fDate,
		TillDate:   tillDate,
		NumPersons: nPersons,
		Canceled:   canceled,
	}
	insertedBookings, err := bookingStore.InsertBooking(ctx, booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nBooking ID is %s\n", insertedBookings.ID.Hex())
	return insertedBookings
}

func main() {
	var err error
	ctx := context.Background()
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewHotelStore(client)
	store := db.Store{
		User:  db.NewMongoUserStore(client),
		Hotel: hotelStore,
		Room:  db.NewRoomStore(client, hotelStore),
	}
	user := fixtures.AddUser(&store, "Pedro", "Escobar", false)
	fmt.Println(user)

	return
	seedHotel("Hilton", "GB", 10)
	seedHotel("Hotel-Grand", "CZ", 10)
	seedHotel("Biskup", "SK", 2)
	pedroUser := SeedUser("Pedro", "Escobar", "pedroEscobar@email.cz", "SuperSecretPassword69", false)
	SeedUser("admin", "admin", "admin@admin.cz", "admin", true)
	RammsteinHotel := seedHotel("Ramstein", "GR", 2)
	seededRoom1 := seedRoom("small", true, 89.99, RammsteinHotel.ID)
	seededRoom2 := seedRoom("medium", true, 189.99, RammsteinHotel.ID)
	seededRoom3 := seedRoom("large", true, 289.99, RammsteinHotel.ID)
	seedBooking(pedroUser.ID, seededRoom1.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2, false)
	seedBooking(pedroUser.ID, seededRoom2.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2, true)
	seedBooking(pedroUser.ID, seededRoom3.ID, time.Now(), time.Now().AddDate(0, 0, 2), 2, false)
}

func init() {
	var err error
	ctx := context.Background()
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewHotelStore(client)
	roomStore = db.NewRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

}
