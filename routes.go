package main

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type Cliente struct {
	Id     int64
	Saldo  int64
	Limite int64
}

func criarTransacao(c *fiber.Ctx) error {
	id := c.Params("id")
	clienteId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var req CriarTransacaoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	if !req.isValid() {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	tx, err := db.Begin(c.Context())
	if err != nil {
		log.Println("erro ao iniciar transação:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer tx.Rollback(c.Context())

	_, err = tx.Exec(c.Context(), "SELECT pg_advisory_xact_lock($1)", clienteId)
	if err != nil {
		log.Println("erro ao adquirir lock:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var cliente Cliente
	err = tx.QueryRow(c.Context(), queryCliente, clienteId).Scan(&cliente.Id, &cliente.Limite, &cliente.Saldo)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		log.Println("erro ao recuperar dados do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if req.Tipo == "c" {
		cliente.Saldo += req.Valor
	} else {
		cliente.Saldo -= req.Valor
	}

	if (cliente.Saldo + cliente.Limite) < 0 {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	_, err = tx.Exec(c.Context(), queryCriarTransacao, cliente.Id, req.Tipo, req.Valor, req.Descricao)
	if err != nil {
		log.Println("erro ao registrar transação:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	_, err = tx.Exec(c.Context(), queryAtualizarSaldoCliente, cliente.Saldo, cliente.Id)
	if err != nil {
		log.Println("erro ao atualizar saldo do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = tx.Commit(c.Context())
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(CriarTransacaoResponse{
		Saldo:  cliente.Saldo,
		Limite: cliente.Limite,
	})
}

func extrato(c *fiber.Ctx) error {
	id := c.Params("id")
	clienteId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var cliente Cliente
	err = db.QueryRow(c.Context(), queryCliente, clienteId).Scan(&cliente.Id, &cliente.Limite, &cliente.Saldo)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		log.Println("erro ao recuperar dados do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var saldo Saldo
	err = db.QueryRow(c.Context(), querySaldoExtrato, cliente.Id).Scan(&saldo.Total, &saldo.Limite, &saldo.DataExtrato)
	if err != nil {
		log.Println("erro ao recuperar informações do saldo:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	rows, err := db.Query(c.Context(), queryTransacoesExtrato, cliente.Id)
	if err != nil {
		log.Println("erro ao recuperar informações das transações:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer rows.Close()

	transacoes := make([]Transacao, 0)
	for rows.Next() {
		var (
			tipo        sql.NullString
			valor       sql.NullInt64
			descricao   sql.NullString
			realizadaEm sql.NullTime
		)
		err = rows.Scan(&tipo, &valor, &descricao, &realizadaEm)
		if err != nil {
			log.Println("erro ao recuperar informações das transações:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if valor.Int64 == 0 {
			continue
		}

		transacoes = append(transacoes, Transacao{
			Tipo:        tipo.String,
			Valor:       valor.Int64,
			Descricao:   descricao.String,
			RealizadaEm: realizadaEm.Time,
		})
	}

	return c.Status(fiber.StatusOK).JSON(ExtratoResponse{
		Saldo:             saldo,
		UltimasTransacoes: transacoes,
	})
}
