package utils

import(
	"strings"
)

func ValidateEmail(email string) bool{
	return strings.Contains(email,"@")&&strings.Contains(email,".")
}

func ValidateUsername(username string)bool{
	return len(username) >=3
}

func ValidatePassword(password string)bool{
	return len(password) >= 6
}