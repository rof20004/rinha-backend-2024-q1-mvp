package main

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:     true,
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Post("/clientes/:id/transacoes", criarTransacao)
	app.Get("/clientes/:id/extrato", extrato)
	log.Fatalln(app.Listen(":8080"))
}
