package api

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/N30xCz/HotelReservationApi/db"
	"github.com/N30xCz/HotelReservationApi/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (a *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := a.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
	}

	if !types.IsPasswordValid(user.EncryptedPassword, params.Password) {
		fmt.Errorf("invalid credentials")
	}
	resp := AuthResponse{
		User:  user,
		Token: createTokenFromUser(user),
	}

	return c.JSON(resp)
}
func createTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret")
	}
	return tokenStr
}
