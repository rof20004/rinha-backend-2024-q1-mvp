package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("/clientes/:id/transacoes", criarTransacao)
	app.Get("/clientes/:id/extrato", extrato)
	log.Fatalln(app.Listen(":8080"))
}
