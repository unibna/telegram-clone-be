package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"path/filepath"
	"time"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error getting file",
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	filepath := fmt.Sprintf("%s/%s", h.uploadDir, filename)

	if err := c.SaveFile(file, filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error saving file",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": fiber.StatusOK,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"url": fmt.Sprintf("/uploads/%s", filename),
		},
	})
} 