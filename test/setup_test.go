package test

import (
	"log"
	"os"
	"project-0/db"
	"project-0/pkg/imagekit"
	"project-0/routes"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var App *fiber.App

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}
	log.Println("===> Initializing DB for tests")
	db.Init()
	imagekit.Init()
	App = fiber.New()
	App = fiber.New(fiber.Config{
		ReadTimeout: time.Duration(20)*time.Second,
	})

	routes.Setup(App)
	code := m.Run()
	log.Println("===> Tests finished")
	os.Exit(code)
}
