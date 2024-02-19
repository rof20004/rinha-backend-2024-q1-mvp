package main

type CriarTransacaoRequest struct {
	Valor     int64  `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

func (r *CriarTransacaoRequest) isValid() bool {
	if r.Valor <= 0 {
		return false
	}

	if r.Tipo != "c" && r.Tipo != "d" {
		return false
	}

	if len(r.Descricao) < 1 || len(r.Descricao) > 10 {
		return false
	}

	return true
}
