package handlers

import (
	"kimsha/internal/services"
	"kimsha/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct{ svc *services.StorageService }

func NewUploadHandler(svc *services.StorageService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "image file required")
	}

	if file.Size > 5*1024*1024 {
		return response.BadRequest(c, "image must be under 5MB")
	}

	ct := file.Header.Get("Content-Type")
	allowed := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	if !allowed[ct] {
		return response.BadRequest(c, "only jpeg, png, webp images are allowed")
	}

	f, err := file.Open()
	if err != nil {
		return response.Internal(c)
	}
	defer f.Close()

	url, err := h.svc.UploadImage(f, file)
	if err != nil {
		return response.Internal(c)
	}

	return response.OK(c, fiber.Map{"url": url})
}
