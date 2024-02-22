package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func criarTransacao(c *fiber.Ctx) error {
	var req CriarTransacaoRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	var (
		id        = c.Params("id")
		tipo      = req.Tipo
		valor     = req.Valor
		descricao = req.Descricao
	)

	if (valor <= 0) || (tipo != "c" && tipo != "d") || (len(descricao) < 1 || len(descricao) > 10) {
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	tx, err := db.Begin(c.Context())
	if err != nil {
		log.Println("erro ao iniciar transação:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	defer tx.Rollback(c.Context())

	var (
		limite int64
		saldo  int64
	)

	err = tx.QueryRow(c.Context(), queryClienteLimiteAndSaldoForUpdate, id).Scan(&limite, &saldo)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		log.Println("erro ao recuperar limite e saldo do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	novoSaldo := valor
	if tipo == "d" {
		if (saldo - valor) < -limite {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		novoSaldo = -valor
	}

	_, err = tx.Exec(c.Context(), queryCriarTransacao, id, tipo, valor, descricao)
	if err != nil {
		log.Println("erro ao registrar transação:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	_, err = tx.Exec(c.Context(), queryAtualizarSaldo, novoSaldo, id)
	if err != nil {
		log.Println("erro ao atualizar saldo do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = tx.Commit(c.Context())
	if err != nil {
		log.Println("erro ao commitar transação:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"saldo":  saldo,
		"limite": limite,
	})
}

func extrato(c *fiber.Ctx) error {
	var (
		id     = c.Params("id")
		limite int64
		saldo  int64
	)
	err := db.QueryRow(c.Context(), queryClienteLimiteAndSaldo, id).Scan(&limite, &saldo)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		log.Println("erro ao recuperar limite e saldo do cliente:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	rows, err := db.Query(c.Context(), queryTransacoesExtrato, id)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"saldo": Saldo{
			Total:       saldo,
			DataExtrato: time.Now(),
			Limite:      limite,
		},
		"ultimas_transacoes": transacoes,
	})
}
