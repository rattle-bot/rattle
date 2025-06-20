package handlers

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilyxenc/rattle/internal/database"
	"github.com/ilyxenc/rattle/internal/models"
)

func CreateContainer(c *fiber.Ctx) error {
	input := new(saveContainerInput)

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

	container := models.Container{
		Type:  input.Type,
		Value: input.Value,
		Mode:  input.Mode,
	}

	if err := db.Create(&container).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to create container",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(Res{
		Message: "Container created",
		Data:    container,
	})
}

func ListContainers(c *fiber.Ctx) error {
	db := database.DB
	var containers []models.Container

	if err := db.Order("created_at DESC").Find(&containers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to retrieve containers",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "List of containers",
		Data:    containers,
	})
}

func ListRunningContainers(c *fiber.Ctx) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to initialize Docker client",
		})
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: false, // Only running containers
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to list containers",
		})
	}

	result := make([]getRunningContainer, 0, len(containers))
	for _, c := range containers {
		result = append(result, getRunningContainer{
			ID:      c.ID,
			Name:    c.Names[0],
			Image:   c.Image,
			Labels:  c.Labels,
			ShortID: c.ID[:12],
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Containers received",
		Data:    result,
	})
}

func UpdateContainer(c *fiber.Ctx) error {
	id := c.Params("id")

	input := new(saveContainerInput)
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

	result := db.Model(&models.Container{}).Where("id = ?", id).Updates(map[string]any{
		"type":  input.Type,
		"value": input.Value,
		"mode":  input.Mode,
	})

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to update container",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Container not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Container updated",
	})
}

func DeleteContainer(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.DB

	result := db.Delete(&models.Container{}, "id = ?", id)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Res{
			Message: "Failed to delete container",
		})
	}
	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(Res{
			Message: "Container not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(Res{
		Message: "Container deleted",
	})
}
