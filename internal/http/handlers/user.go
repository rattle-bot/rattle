package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/config"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/http"
	"github.com/ilyxenc/rattle/internal/models"
	"gorm.io/gorm"
)

func AuthTelegram(c *fiber.Ctx) error {
	input := new(telegramDataInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Invalid request body",
		})
	}

	vldt := validator.New()

	if err := vldt.Struct(input); err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Validation failed",
		})
	}

	// Verify the integrity of the data received from the Telegram Web App
	token := config.Cfg.BotToken

	// Create a new HMAC using SHA-256 and the "WebAppData" key
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(token))  // Hash the bot token using HMAC-SHA256
	secretKey := h.Sum(nil) // Generate secret key based on token

	// Create a new HMAC with the derived secret key
	h = hmac.New(sha256.New, secretKey)

	// Decode the check string from URL-encoded format
	decodedString, err := url.QueryUnescape(input.CheckDataString)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to decode check string",
		})
	}

	// Hash the decoded check string using HMAC-SHA256
	h.Write([]byte(decodedString))
	// Generate computed hash as a hex string
	computedHash := hex.EncodeToString(h.Sum(nil))

	if computedHash != input.Hash {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Hash mismatch — data may have been tampered with",
		})
	}

	params := strings.Split(decodedString, "\n")

	var user map[string]any
	for _, param := range params {
		kv := strings.SplitN(param, "=", 2)
		if kv[0] == "user" {
			if err := json.Unmarshal([]byte(kv[1]), &user); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(Res{
					Message: "Failed to parse user data",
				})
			}
		}
	}

	tgID := fmt.Sprintf("%.0f", user["id"])

	db := database.DB

	var count int64
	db.Model(&models.User{}).Count(&count)

	var u models.User
	err = db.First(&u, "telegram_id = ?", tgID).Error

	// If no users exist, create the first one as admin
	if count == 0 && err == gorm.ErrRecordNotFound {
		newUser := models.User{
			TelegramID: tgID,
			Username:   user["username"].(string),
			FirstName:  user["first_name"].(string),
			Role:       models.RoleAdmin,
		}
		if err := db.Create(&newUser).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Res{
				Message: "Failed to create admin user",
			})
		}

		// Get access_token
		accessToken, err := http.GenerateAccessToken(u.TelegramID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Res{
				Message: "Failed to generate access token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(Res{
			Message: "Admin account created and access token issued",
			Data: fiber.Map{
				"access_token": accessToken,
			},
		})
	}

	// If user exists and has access
	if err == nil && u.ID > 0 {
		if err := db.Model(&models.User{}).
			Where("telegram_id = ?", tgID).
			Updates(models.User{Username: user["username"].(string), FirstName: user["first_name"].(string)}).
			Error; err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(Res{
				Message: "Failed to save user",
			})
		}

		accessToken, err := http.GenerateAccessToken(u.TelegramID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Res{
				Message: "Failed to generate access token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(Res{
			Message: "Access token issued",
			Data: fiber.Map{
				"access_token": accessToken,
			},
		})
	}

	// User not found and access denied
	return c.Status(fiber.StatusForbidden).JSON(Res{
		Message: "Access denied — user is not allowed",
	})
}
