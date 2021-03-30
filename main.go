package main

import (
	"email-counter/connector"
	"email-counter/service"
	"log"
	"os"

	"github.com/gofiber/fiber"
)

const defaultPort = "8000"

func main() {

	log.Println("Running email counter...")

	app := fiber.New()
	app.Get("/health", service.HealthCheck)
	app.Post("/", service.CreateList)
	app.Put("/:iteration", service.UpdateEmailsSentCounter)
	app.Get("/", service.GetAll)
	app.Get("/:iteration", service.GetListReportByIteration)
	app.Put("/archive/:id", service.ArchiveIteration)

	app.Options("*", service.Options)

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Cannot find $PORT variable, using default value %s instead", defaultPort)
		port = defaultPort
	}

	app.Listen(port)
}

func init() {
	dbHostAddr := os.Getenv("DATABASE_URI")
	log.Printf("hostAddr: %s", dbHostAddr)
	connector.SetDatabaseAddress(dbHostAddr)
}
