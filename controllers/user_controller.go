package controllers

import (
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
	"github.com/kithmina1999/eldraread-api/utils"
)

func RegisterUser(c *fiber.Ctx) error {
	type RegisterInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}

	var data RegisterInput

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email,password or username",
		})
	}

	if !utils.ValidateEmail(data.Email) || !utils.ValidatePassword(data.Password) || !utils.ValidateUsername(data.Username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid email, password, or username",
		})
	}

	authClient, err := config.FirebaseAuth()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "firebase init error",
		})
	}

	params := (&auth.UserToCreate{}).Email(data.Email).Password(data.Password).DisplayName(data.Username)

	userRecord,err:= authClient.CreateUser(c.Context(),params)
	if err != nil{
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error":err.Error(),
			})
	}
	
	return c.JSON(fiber.Map{
		"message":"User registered successfully",
		"uid":userRecord.UID,
	})
}

func RegisterWithGoogle(c *fiber.Ctx)error{
	type GoogleInput struct{
		IdToken string `json:"idToken"`
	}

	var input GoogleInput
	if err:= c.BodyParser(&input); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error":"Invalid request",
		})
	}

	authClient,err:= config.FirebaseAuth()
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Auth error",
		})
	}

	// verfiy google token
	token,err := authClient.VerifyIDToken(c.Context(),input.IdToken)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error":"Invalid token",
		})
	}

	//fetch user info
	user, err := authClient.GetUser(c.Context(),token.UID)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"could not fetch user",
		})
	} 
	return c.JSON(fiber.Map{
		"message":"Google login success",
		"user":fiber.Map{
			"uid":user.UID,
			"email":user.Email,
			"username":user.DisplayName,
		},
	})
}
