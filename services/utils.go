package services

import (
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"strings"
)

func IsMp4(c *fiber.Ctx, file *multipart.FileHeader) (error, bool) {
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".mp4") {
		return c.Status(400).SendString("File must be mp4 format"), true
	}
	return nil, false
}
