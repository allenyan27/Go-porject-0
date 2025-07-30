package routes

import (
	"project-0/handlers"
	"project-0/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/routes", func(c *fiber.Ctx) error {
		routes := app.GetRoutes()
		return c.JSON(routes)
	})
	
	// Villas group
	villas := api.Group("/villas")
	villas.Get("/", handlers.GetVillas).Name("Getting all Villas")
	villas.Get("/:id", handlers.GetSingleVilla).Name("Getting single villa")
	villas.Post("/", handlers.PostVilla).Name("Add new villa")
	villas.Put("/:id", 
		middlewares.CheckUpdateImage, 
		handlers.UpdateVilla,
	).Name("Update villa")
	villas.Delete("/:id", 
		middlewares.DeleteGalleryByVillaId, 
		middlewares.DeleteImage, 
		handlers.DeleteVilla,
	).Name("Delete villa")

	gallery := api.Group("/gallery")
	gallery.Get("/:id", handlers.GetGalleryVilla).Name("Get galleries villa")
	gallery.Post("/", handlers.AddGallery).Name("Add gallery")
	gallery.Post("/delete", 
		middlewares.DeleteGalleryByID, 
		handlers.DeleteGalleryByID,
	).Name("Delete multiple gallery using Id")

	
	// api.Post("/bookings", handlers.CreateBooking)
	// api.Get("/prices/:villaId", handlers.GetPrices)
}
