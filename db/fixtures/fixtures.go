package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/N30xCz/HotelReservationApi/api"
	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {

	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s_%s@email.cz", fname, lname),
		FirstName: fname,
		LastName:  lname,
		Password:  fmt.Sprintf("%s_%s", fname, lname),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s --> %s\n", user.Email, api.CreateTokenFromUser(user))
	return insertedUser
}
