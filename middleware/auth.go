package middleware

import (
	// "fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	

)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") { //checking token starts with bearer //
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Malformed token",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ") //to remove bearer//

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			return []byte(os.Getenv("JWT_SECRET")), nil     //to verify the token//
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims) //data inside the JWT//
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		c.Locals("user", claims) //Fiber feature requesr context//
		return c.Next()
	}
}
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims, ok := c.Locals("user").(jwt.MapClaims) //converting it back into its real type we can acces user info//
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		role, ok := userClaims["role"].(string) //checking rile in user info//
		if !ok || role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access requried",
			})
		}
		return c.Next()
	}
}
