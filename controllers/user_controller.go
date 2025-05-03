package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
	"github.com/kithmina1999/eldraread-api/db"
	"github.com/kithmina1999/eldraread-api/models"
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

	userRecord, err := authClient.CreateUser(c.Context(), params)
	if err != nil {
		if auth.IsEmailAlreadyExists(err) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Email is already registered",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	//save data in postgres
	newUser := models.User{
		UID:userRecord.UID,
		Email: userRecord.Email,
		Username: userRecord.DisplayName,
		Status: models.StatusActive,
		Role:models.RoleUser,
	}

	if err:= db.DB.Create(&newUser).Error; err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"Failed to save user in database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
		"uid":     userRecord.UID,
	})
}

func RegisterWithGoogle(c *fiber.Ctx) error {
	type GoogleInput struct {
		IdToken string `json:"idToken"`
	}

	var input GoogleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	authClient, err := config.FirebaseAuth()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Auth error",
		})
	}

	// verfiy google token
	token, err := authClient.VerifyIDToken(c.Context(), input.IdToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	//fetch user info
	user, err := authClient.GetUser(c.Context(), token.UID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not fetch user",
		})
	}

	if !user.EmailVerified{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Google account email is not verified",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Google login success",
		"user": fiber.Map{
			"uid":      user.UID,
			"email":    user.Email,
			"username": user.DisplayName,
		},
	})
}

func LoginUser(c *fiber.Ctx)error {
	type LoginInput struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err:=c.BodyParser(&input); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

	//firebase api req
	payload := map[string]string{
		"email":             input.Email,
		"password":          input.Password,
		"returnSecureToken": "true",
	}
	jsonPayload,_:=json.Marshal(payload)
	
	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	if apiKey == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "API key not configured",
		})
	}
	url:= fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	resp,err:=http.Post(url,"application/json",bytes.NewBuffer(jsonPayload))
	if err != nil || resp.StatusCode != 200 {	
		body,_:=io.ReadAll(resp.Body)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":"Invalid email or password",
			"debug":string(body),
		})
	}
	defer resp.Body.Close()

	//parse fireabse response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	c.Cookie(&fiber.Cookie{
		Name:     "session_token",
		Value:    result["idToken"].(string),
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // set true in prod w/ HTTPS
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message":"Login successful",
		"idToken":result["idToken"],
		"user":fiber.Map{
			"email":    result["email"],
			"localId":  result["localId"],
			"displayName": result["displayName"],
		},
	})
}