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
			Data:    err.Error(),
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
