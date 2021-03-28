package main

import (
	"email-counter/connector"
	"email-counter/service"
	"flag"
	"log"

	"github.com/gofiber/fiber"
)

const port = 8000

var dbHostAddr string

func main() {

	log.Println("Running email counter...")

	app := fiber.New()
	app.Get("/health", service.HealthCheck)
	app.Post("/", service.CreateList)
	app.Post("/:iteration", service.UpdateEmailsSentCounter)
	app.Get("/:iteration", service.GetListReportByIteration)

	app.Listen(port)
}

func init() {
	flag.StringVar(&dbHostAddr, "db", "", "Database connection path, in format host:port")
	flag.Parse()
	connector.SetDatabaseAddress(dbHostAddr)
}
