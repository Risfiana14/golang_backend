package services
import (
"strconv"
"strings"
"github.com/gofiber/fiber/v2"
"tugas5/app/model"
"tugas5/app/repository"
)

	// GetUsersService -> service untuk ambil data user dengan pagination, search, sorting
func GetUsersService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// Validasi input
	sortByWhitelist := map[string]bool{"id": true, "name": true, "email": true, "created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

// Ambil data dari repository
	users, err := repository.GetUsersRepo(search, sortBy, order,
limit, offset)
if err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Failed to fetchusers"})
}
total, err := repository.CountUsersRepo(search)
if err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Failed to countusers"})
}

// Buat response pakai model
response := model.UserResponse{
Data: users,
Meta: model.MetaInfo{
Page: page,
Limit: limit,
Total: total,
Pages: (total + limit - 1) / limit,
SortBy: sortBy,
Order: order,
Search: search,
},
}
return c.JSON(response)
}

func GetProfileHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(model.User)
	return c.JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}