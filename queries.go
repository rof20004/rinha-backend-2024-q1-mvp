package main

const queryClienteLimiteAndSaldo = `SELECT 
    									  c.limite, 
    									  s.valor 
									FROM clientes c 
									INNER JOIN saldos s ON s.cliente_id = c.id 
									WHERE c.id = $1`

const queryClienteLimiteAndSaldoForUpdate = `SELECT 
    											   c.limite, 
    											   s.valor 
											 FROM clientes c 
											 INNER JOIN saldos s ON s.cliente_id = c.id 
											 WHERE c.id = $1 FOR UPDATE`

const queryAtualizarSaldo = `UPDATE 
    							   saldos
						     SET valor = valor + $1
							 WHERE cliente_id = $2`

const queryCriarTransacao = `INSERT INTO transacoes(cliente_id, tipo, valor, descricao) VALUES ($1, $2, $3, $4)`

const queryTransacoesExtrato = `SELECT 
    							      tipo,
    						          valor,
    						          descricao,
    						          realizada_em
                           		FROM transacoes
                           		WHERE cliente_id = $1
                           		ORDER BY realizada_em DESC
                           		LIMIT 10`
