package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

func CreateChat(c *fiber.Ctx) error {
	input := new(createChatInput)

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

	chat := models.Chat{
		ChatID: input.ChatID,
		Send:   true,
	}
	if input.Send != nil {
		chat.Send = *input.Send
	}

	if err := db.Create(&chat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to create chat",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Res{
		Message: "Chat created",
		Data:    chat,
	})
}

func ListChats(c *fiber.Ctx) error {
	db := database.DB
	var chats []models.Chat

	if err := db.Order("created_at DESC").Find(&chats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to retrieve chats",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "List of chats",
		Data:    chats,
	})
}

func UpdateChat(c *fiber.Ctx) error {
	id := c.Params("id")

	input := new(updateChatInput)
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

	result := db.Model(&models.Chat{}).Where("id = ?", id).Update("send", input.Send)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to update chat",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Chat not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Chat updated",
	})
}

func DeleteChat(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	result := db.Delete(&models.Chat{}, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to delete chat",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Chat not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Chat deleted",
	})
}
