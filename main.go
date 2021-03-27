package main

import (
	"email-counter/service"
	"log"

	"github.com/gofiber/fiber"
)

const port = 8000

func main() {

	log.Println("Running email counter...")

	app := fiber.New()
	app.Get("/health", service.HealthCheck)
	app.Post("/", service.CreateList)
	app.Post("/:iteration", service.UpdateEmailsSentCounter)
	app.Get("/:iteration", service.GetListReportByIteration)

	app.Listen(port)
}
