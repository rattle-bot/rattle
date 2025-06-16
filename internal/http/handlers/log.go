package handlers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

func CreateLog(c *fiber.Ctx) error {
	input := new(createLogInput)

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

	log := models.LogExclusion{
		Pattern:   input.Pattern,
		MatchType: input.MatchType,
		EventType: input.EventType,
	}

	if err := db.Create(&log).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to create log",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Res{
		Message: "Log created",
		Data:    log,
	})
}

func ListLog(c *fiber.Ctx) error {
	db := database.DB
	var logs []models.LogExclusion

	if err := db.Order("created_at DESC").Find(&logs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to retrieve logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "List of logs",
		Data:    logs,
	})
}

func UpdateLog(c *fiber.Ctx) error {
	id := c.Params("id")

	input := new(updateLogInput)
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

	if input.Pattern != nil {
		if _, err := regexp.Compile(*input.Pattern); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(Res{
				Message: "Invalid regex pattern",
			})
		}
	}

	updates := map[string]interface{}{}
	if input.Pattern != nil {
		updates["pattern"] = *input.Pattern
	}
	if input.MatchType != nil {
		updates["match_type"] = *input.MatchType
	}
	if input.EventType != nil {
		updates["event_type"] = *input.EventType
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(Res{
			Message: "No valid fields provided for update",
		})
	}

	db := database.DB

	result := db.Model(&models.LogExclusion{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to update log",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Log not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Log updated",
	})
}

func DeleteLog(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	result := db.Delete(&models.LogExclusion{}, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to delete log",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Log not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Log deleted",
	})
}
