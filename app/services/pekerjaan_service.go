package services

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
	"tugas5/app/model"
	"tugas5/app/repository"
	"tugas5/utils"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	repo repository.PekerjaanRepository
}

func NewPekerjaanService(repo repository.PekerjaanRepository) *PekerjaanService {
	return &PekerjaanService{repo: repo}
}

// GET /pekerjaan?page=&limit=&sortBy=&order=&search=
func (s *PekerjaanService) GetAllService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit

	sortByWhitelist := map[string]bool{
		"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "created_at": true,
	}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	pekerjaan, err := s.repo.GetAll(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	response := fiber.Map{
		"data": pekerjaan,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	}
	return c.JSON(response)
}

// GET /pekerjaan/:id
func (s *PekerjaanService) GetByIDService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID tidak valid"})
	}
	data, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// GET /pekerjaan/alumni/:alumni_id
func (s *PekerjaanService) GetByAlumniIDService(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Alumni ID tidak valid"})
	}
	data, err := s.repo.GetByAlumniID(alumniID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// POST /pekerjaan
func (s *PekerjaanService) CreateService(c *fiber.Ctx) error {
	var req model.CreatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)

	// User biasa hanya bisa membuat pekerjaan untuk dirinya sendiri
	if role != "admin" {
		req.AlumniID = userID
	}

	// Simpan siapa yang membuat
	req.CreatedBy = utils.StringPtr(username)

	// Validasi tanggal
	if req.TanggalMulaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal mulai salah"})
		}
		req.TanggalMulaiKerja = t.Format("2006-01-02")
	}
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal selesai salah"})
		}
		formatted := t.Format("2006-01-02")
		req.TanggalSelesaiKerja = &formatted
	}

	data, err := s.repo.Create(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// PUT /pekerjaan/:id
func (s *PekerjaanService) UpdateService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var req model.UpdatePekerjaanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request body tidak valid"})
	}

	// Validasi tanggal
	if req.TanggalMulaiKerja != "" {
		t, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal mulai salah"})
		}
		req.TanggalMulaiKerja = t.Format("2006-01-02")
	}
	if req.TanggalSelesaiKerja != nil && *req.TanggalSelesaiKerja != "" {
		t, err := time.Parse("2006-01-02", *req.TanggalSelesaiKerja)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tanggal selesai salah"})
		}
		formatted := t.Format("2006-01-02")
		req.TanggalSelesaiKerja = &formatted
	}

	data, err := s.repo.Update(id, req)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"success": true, "data": data})
}

// DELETE /pekerjaan/:id
func (s *PekerjaanService) DeleteService(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	role := c.Locals("role").(string)
	username := c.Locals("username").(string)

	pekerjaan, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Validasi hak akses: hanya admin atau pembuat (CreatedBy) yang boleh hapus
	if role != "admin" && (pekerjaan.CreatedBy == nil || *pekerjaan.CreatedBy != username) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Anda tidak memiliki izin untuk menghapus pekerjaan ini",
		})
	}

	if err := s.repo.Delete(id); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{"error": "Pekerjaan tidak ditemukan"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil dihapus (soft delete)",
	})
}
