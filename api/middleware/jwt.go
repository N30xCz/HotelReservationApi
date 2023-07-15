package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--JWT--Authing")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unautheticated")
	}
	claims, err := ValidateToken(token)
	if err != nil {
		fmt.Errorf("unauthorized")
	}
	expiresFloat := claims["expires"].(float64)
	expires := int64(expiresFloat)
	// check token expiration
	if time.Now().Unix() > expires {
		return fmt.Errorf("token is expired")
	}
	return c.Next()
}

func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT Token", err)
		return nil, fmt.Errorf("Unauthorized")
	}
	if !token.Valid {
		fmt.Println("Invalid Token")
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	return claims, nil

}
