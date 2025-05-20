package mw

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func CheckAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(401, map[string]string{"error": "Missing Authorization header"})
		}

		tokenString := authHeader[len("Bearer "):]
		userToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_ACCESS_KEY")), nil
		})

		if err != nil || !userToken.Valid {
			return c.JSON(401, map[string]string{"error": "Unauthorized"})
		}

		claims, ok := userToken.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(401, map[string]string{"error": "Invalid token claims"})
		}

		// fmt.Println(claims)

		user_id := claims["user_id"].(float64)
		email := claims["email"].(string)
		exp := claims["exp"].(float64)

		// Set the user ID in the context for further use
		c.Set("user_id", user_id) // Changed from "id" to "user_id"
		c.Set("email", email)
		c.Set("exp", exp)

		// fmt.Printf("User ID: %v, Email: %s, Expiration: %v\n", user_id, email, exp)

		return next(c)
	}
}
