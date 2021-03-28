package main

import (
	"email-counter/connector"
	"email-counter/service"
	"flag"
	"log"
	"os"

	"github.com/gofiber/fiber"
)

const defaultPort = "8000"

var dbHostAddr string

func main() {

	log.Println("Running email counter...")

	app := fiber.New()
	app.Get("/health", service.HealthCheck)
	app.Post("/", service.CreateList)
	app.Post("/:iteration", service.UpdateEmailsSentCounter)
	app.Get("/:iteration", service.GetListReportByIteration)

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Cannot find $PORT variable, using default value %s instead", defaultPort)
		port = defaultPort
	}

	app.Listen(port)
}

func init() {
	flag.StringVar(&dbHostAddr, "db", "", "Database connection path, in format host:port")
	flag.Parse()
	connector.SetDatabaseAddress(dbHostAddr)
}
