package main

import (
	"context"
	"log"

	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func SeedUser(fname, lname, email, password string, isAdmin bool) {

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
	_, err = userStore.InsertUser(ctx, user)
}

func seedHotel(name string, location string, raiting int) {

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
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "normal",
			Price: 199.9,
		},
		{
			Size:  "kingsize",
			Price: 122.9,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {

	seedHotel("Hilton", "GB", 10)
	seedHotel("Hotel-Grand", "CZ", 10)
	seedHotel("Biskup", "SK", 2)
	SeedUser("Pedro", "Escobar", "pedroEscobar@email.cz", "SuperSecretPassword69", false)
	SeedUser("admin", "admin", "admin@admin.cz", "admin", true)
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

}
