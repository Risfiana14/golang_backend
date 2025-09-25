package routes

import (
	"tugas5/app/repository"
	"tugas5/app/services"
	"tugas5/database"
	"tugas5/middleware"

	"github.com/gofiber/fiber/v2"
)

	// UserRoutes -> definisi route untuk user
func UserRoutes(app *fiber.App) {
	alumniRepo := repository.NewAlumniRepository(database.DB)
	pekerjaanRepo := repository.NewPekerjaanRepository(database.DB)
	
	api := app.Group("/api")
	app.Get("/users", services.GetUsersService)

	// Init service
	alumniSvc := services.NewAlumniService(alumniRepo)
	pekerjaanSvc := services.NewPekerjaanService(pekerjaanRepo)

	// ---------- AUTH ----------
	api.Post("/login", func(c *fiber.Ctx) error {
		return services.LoginService(c, database.DB)
	})

	protected := api.Group("", middleware.AuthRequired())
	protected.Get("/profile", services.GetProfileHandler)

	// ---------- ALUMNI ----------
	protected.Get("/alumni", alumniSvc.GetAllService)
	protected.Get("/alumni/:id", alumniSvc.GetByIDService)
	protected.Post("/alumni", middleware.AuthRequired(), alumniSvc.CreateService)
	protected.Put("/alumni/:id", middleware.AuthRequired(), alumniSvc.UpdateService)
	protected.Delete("/alumni/:id", middleware.AuthRequired(), alumniSvc.DeleteService)

	// ---------- PEKERJAAN ----------
	protected.Get("/pekerjaan", pekerjaanSvc.GetAllHandler)
	protected.Get("/pekerjaan/:id", pekerjaanSvc.GetByIDHandler)
	protected.Get("/pekerjaan/alumni/:alumni_id", pekerjaanSvc.GetByAlumniIDHandler)
	protected.Post("/pekerjaan", middleware.AuthRequired(), pekerjaanSvc.CreateHandler)
	protected.Put("/pekerjaan/:id", middleware.AuthRequired(), pekerjaanSvc.UpdateHandler)
	protected.Delete("/pekerjaan/:id", middleware.AuthRequired(), pekerjaanSvc.DeleteHandler)
}