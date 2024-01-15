package middleware

import (
	"arxivhub/internal/delivery/http/jwt"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(jwt *jwt.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println(c.GetReqHeaders())
		authHeader := c.Get("Authorization")

		claims, err := AuthHeaderValidation(authHeader, jwt)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}

func AuthHeaderValidation(token string, jwt *jwt.JWTManager) (*jwt.Claims, error) {
	//fields := strings.Fields(authHeader)
	//if len(fields) < 2 {
	//	return nil, fmt.Errorf("invalid authorization header format")
	//}

	//authorizationType := strings.ToLower(fields[0])
	//if authorizationType != "bearer" {
	//	return nil, fmt.Errorf("unsupported authorization type %s", authorizationType)
	//}

	//token := fields[1]

	return jwt.ParseToken(token)
}
