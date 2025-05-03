package middlewares

import(
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
)

func RequireAuth(c *fiber.Ctx) error{
	tokenStr:=c.Cookies("session_token")
	if tokenStr == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Unauthorized access",
		})
	}

	authClient,err:=config.FirebaseAuth()
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Firebase initialization error",
		})
	}

	token,err:=authClient.VerifyIDToken(c.Context(),tokenStr)
	if err != nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Token is invalid or expired",
		})
	}
	c.Locals("uid",token.UID)
	return c.Next()
}