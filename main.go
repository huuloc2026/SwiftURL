package main

import (
	"log"

	"github.com/gofiber/fiber"
)

func main() {
	db := database.NewPostgres()
	repo := repository.NewURLRepository(db)
	service := service.NewURLService(repo)
	h := handler.NewURLHandler(service)

	app := fiber.New()

	app.Post("/shorten", h.ShortenURL)
	app.Get("/:code", h.ResolveURL)

	log.Fatal(app.Listen(":8080"))
}