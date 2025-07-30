package main

import (
	"log"
	"os"
	"project-0/db"
	"project-0/pkg/dotenv"
	"project-0/pkg/imagekit"
	"project-0/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load .env variables
	dotenv.Init()

	// Init external services
	imagekit.Init()
	db.Init()

	// Create Fiber app with default config
	app := fiber.New(fiber.Config{
		AppName: "Project 0 API",
		Prefork:       true,
	})

	// ---------- Middleware for security ----------
	// Recover from panics
	app.Use(recover.New())

	// Secure HTTP headers
	app.Use(helmet.New())

	// Compression (GZIP)
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault, // LevelBestSpeed, LevelBestCompression, etc.
	}))

	// Enable CORS (restrict on production)
	app.Use(cors.New())

	// Log all requests (custom format optional)
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Makassar",
	}))

	// Rate limiting: max 100 requests per minute per IP
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
	}))

	// ---------- App routes ----------
	routes.Setup(app)

	// ---------- Start Server ----------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
	
}
