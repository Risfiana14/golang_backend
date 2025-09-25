package services

import (
	"tugas5/app/model"
	"tugas5/app/repository"
	"tugas5/utils"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func LoginService(c *fiber.Ctx, db *sql.DB) error {
	var loginData model.LoginRequest
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"success": false,
		})
	}

	if loginData.Username == "" || loginData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":   "Harap masukkan username dan password",
			"success": false,
		})
	}

	user, err := repository.Login(db, loginData.Username, loginData.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Username atau password salah",
				"error":   err.Error(),
				"success": false,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal terhubung ke database",
			"error":   err.Error(),
			"success": false,
		})
	}

	// Generate JWT token
    token, err := utils.GenerateToken(user)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Gagal membuat token",
            "error":   err.Error(),
            "success": false,
        })
    }

	response := model.LoginResponse{
		Token: token,
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   response.Token,
		"data":    user,
		"message": "Login berhasil",
		"success": true,
	})
}