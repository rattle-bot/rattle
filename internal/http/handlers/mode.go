package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

func GetFilteringMode(c *fiber.Ctx) error {
	db := database.DB

	var mode models.Mode
	if err := db.First(&mode).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Filtering mode not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "",
		Data:    mode,
	})
}

func UpdateFilteringMode(c *fiber.Ctx) error {
	input := new(updateModeInput)
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

	// Update first (and only) record
	var mode models.Mode
	err := db.First(&mode).Error
	if err != nil {
		// No record - create
		newMode := models.Mode{
			Value: input.Value,
		}
		if err := db.Create(&newMode).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(Res{
				Message: "Failed to create filtering mode",
			})
		}

		return c.Status(fiber.StatusOK).JSON(Res{
			Message: "Filtering mode created",
			Data:    newMode,
		})
	}

	if err := db.Model(&mode).Update("value", input.Value).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to update filtering mode",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Filtering mode updated",
		Data:    mode,
	})
}
