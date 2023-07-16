package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     "pedroEscobar@email.cz",
		FirstName: "Pedro",
		LastName:  "Escobar",
		Password:  "SuperSecretPassword69",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}
func TestAuthenticateWrongPassword(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)
	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "pedroEscobar@email.cz",
		Password: "SuperSecrdwwdwdwdwdetPassword69",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {

		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status of 400 but got : %d", resp.StatusCode)
	}

	var genResp genericResp

	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}
	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}
	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be : invalid credentials but got %s", genResp.Msg)
	}
}
func TestAuthenticateSuccess(t *testing.T) {

	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := insertTestUser(t, tdb.UserStore)
	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "pedroEscobar@email.cz",
		Password: "SuperSecretPassword69",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))

	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {

		t.Error(err)
	}
	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		fmt.Println(resp.Status)
		t.Fatal(err)
	}
	if authResponse.Token == "" {
		t.Fatal("expected the JWT token to be present in the auth response.")
	}
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		fmt.Println(insertedUser)
		fmt.Println(authResponse.User)
		t.Fatal("expected the user to be inserted user ")
	}
}
