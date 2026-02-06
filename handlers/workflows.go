package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/umbe77/yasb/models"
	"github.com/umbe77/yasb/store"
)

func GetWorflows(db *store.Storage) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		wfs, err := db.GetWorflows(c.Context(), "{}")
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(models.GetErrorMessage("Cannot get workflows", err))
		}

		return c.JSON(wfs)
	}
}

func AddWorkflow(db *store.Storage) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) <= 0 {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(models.GetErrorMessage("Cannot get workflows", fmt.Errorf("no data")))
		}
		var wf models.Workflow
		if err := json.Unmarshal(c.Body(), &wf); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(models.GetErrorMessage("Cannot get workflows", err))
		}

		if err := db.AddWorkflow(c.Context(), wf); err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(models.GetErrorMessage("Cannot get workflows", err))
		}

		return c.JSON(fiber.Map{
			"message": "Workflow added",
		})
	}
}
