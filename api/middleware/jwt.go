package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--JWT--Authing")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("Unautheticated")
	}
	if err := parseToken(token); err != nil {
		return err
	}
	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		secret := os.Getenv("JWT_SECRET")
		fmt.Println("DELETE ME AFTER DEBUG SECRET IS ", secret)
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT Token", err)
		return fmt.Errorf("Unauthorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unauthorized")

}
