package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
	"github.com/kithmina1999/eldraread-api/db"
	"github.com/kithmina1999/eldraread-api/models"
)

func AdminOnly(c *fiber.Ctx) error {

	
		tokenStr := c.Cookies("admin_session_token")

		if tokenStr == ""{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":"Missing session token",
			})
		}

		authClient, _ := config.FirebaseAuth()
		token, err := authClient.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Retrieve the user from the database
		var user models.User
		if err := db.DB.Where("uid = ?", token.UID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Check if the user has admin role
		if user.Role != models.RoleAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}

		c.Locals("user_id", token.UID)
		return c.Next()

}

func AdminSessionCheck(c *fiber.Ctx)error{
	tokenStr := c.Cookies("admin_session_token")

	if tokenStr == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Missing session token",
		})
	}

	authClient,_:=config.FirebaseAuth()
	token,err:=authClient.VerifyIDToken(context.Background(),tokenStr)
	if err!=nil{
		c.Cookie(&fiber.Cookie{
			Name:     "admin_session_token",
			Value:    "",
			Expires:  time.Now().Add(-1 * time.Hour), // Set expiration to the past
			HTTPOnly: true,
			Secure:   false,
			SameSite: "Lax",
		})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Invalid or expired token",
		})
	}

	var user models.User
	if err := db.DB.Where("uid = ?",token.UID).First(&user).Error;err!=nil{
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":"User not found",
		})
	}

	  // Check if the user has admin role
	  if user.Role != models.RoleAdmin {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": "Admin access required",
        })
    }
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Session is valid",
    })
}
