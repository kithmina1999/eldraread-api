package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kithmina1999/eldraread-api/config"
	"github.com/kithmina1999/eldraread-api/db"
	"github.com/kithmina1999/eldraread-api/models"
)

func AdminLogin(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	//firebase api input to sign in
	payload := map[string]string{
		"email":             input.Email,
		"password":          input.Password,
		"returnSecureToken": "true",
	}
	jsonPayload, _ := json.Marshal(payload)

	apiKey := os.Getenv("FIREBASE_WEB_API_KEY")
	if apiKey == "" {
		return c.Status(500).JSON(fiber.Map{
			"error": "API key is not confirmed",
		})
	}
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
			"debug": string(body),
		})
	}
	defer resp.Body.Close()

	//parse firebase response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Firebase response",
		})
	}

	//verify admin claim
	authClient, err := config.FirebaseAuth()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize with firebase auth",
		})
	}
	//verify id token and fetch user
	tokenStr := result["idToken"].(string)
	token, err := authClient.VerifyIDToken(c.Context(), tokenStr)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	uid := token.UID
	var user models.User
	if err := db.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if user.Role != models.RoleAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}
	// Set the session token for the admin login
	c.Cookie(&fiber.Cookie{
		Name:     "admin_session_token",
		Value:    tokenStr,                       // Use the Firebase ID token or any custom session ID
		Expires:  time.Now().Add(24 * time.Hour), // 24 hours expiration
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Admin login successful",
		"idToken": result["idToken"],
		"user": fiber.Map{
			"email":       result["email"],
			"localId":     result["localId"],
			"displayName": result["displayName"],
		},
	})
}

