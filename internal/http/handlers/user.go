package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

func CreateUser(c *fiber.Ctx) error {
	input := new(createUserInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Invalid request body",
		})
	}

	vldt := validator.New()
	if err := vldt.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Validation failed",
		})
	}

	db := database.DB

	var existing models.User
	if err := db.First(&existing, "telegram_id = ?", input.TelegramID).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(Res{
			Message: "User already exists",
		})
	}

	newUser := models.User{
		TelegramID: input.TelegramID,
		Role:       input.Role,
		Active:     false,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Res{
		Message: "User created successfully",
		Data: fiber.Map{
			"id":          newUser.ID,
			"telegram_id": newUser.TelegramID,
			"role":        newUser.Role,
		},
	})
}

func GetMe(c *fiber.Ctx) error {
	telegramID := c.Locals("telegram_id")

	db := database.DB

	var user models.User

	if err := db.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to retrieve user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "My user data",
		Data:    user,
	})
}

func ListUsers(c *fiber.Ctx) error {
	db := database.DB

	var users []models.User

	if err := db.Order("created_at DESC").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to retrieve users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "List of users",
		Data:    users,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("telegram_id")

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "Telegram ID is required",
		})
	}

	db := database.DB

	var user models.User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "User not found",
		})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "User deleted successfully",
	})
}
