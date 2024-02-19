package main

const queryCliente = `SELECT id, limite, saldo FROM clientes WHERE id = $1`

const queryAtualizarSaldoCliente = `UPDATE 
    							 		  clientes
								    SET saldo = $1
									WHERE id = $2`

const queryCriarTransacao = `INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)`

const querySaldoExtrato = `SELECT 
    							 saldo,
    						     limite,
    						     now() AS data_extrato
                           FROM clientes
                           WHERE id = $1`

const queryTransacoesExtrato = `SELECT 
    							      tipo,
    						          valor,
    						          descricao,
    						          realizada_em
                           		FROM transacoes
                           		WHERE cliente_id = $1
                           		ORDER BY realizada_em DESC
                           		LIMIT 10`
