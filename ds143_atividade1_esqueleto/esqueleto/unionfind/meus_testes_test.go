package unionfind

import (
	"math/rand"
	"testing"
)

// Detecta erros
func TestOraculo(t *testing.T) {
	const n = 100
	const operacoes = 500

	r := rand.New(rand.NewSource(67))

	// inicia todas as variantes de uf com 100 elementos
	instancias := map[string]UF{
		"QuickFind":               NewQuickFind(n),
		"QuickUnion":              NewQuickUnion(n),
		"WeightedQuickUnion":      NewWeightedQuickUnion(n),
		"WeightedPathCompression": NewWeightedPathCompression(n),
	}

	// faz 500 uniões entre os elem.
	for i := 0; i < operacoes; i++ {
		p := r.Intn(n)
		q := r.Intn(n)

		// Aplica a mesma união em todas as implementações
		for _, uf := range instancias {
			uf.Union(p, q)
		}

		// verifica se todas trazem o msm n de componentes
		var countEsperado int
		primeira := true
		for nome, uf := range instancias {
			// se for a primeira
			if primeira {
				countEsperado = uf.Count()
				primeira = false
			} else if uf.Count() != countEsperado {
				t.Fatalf("O count na operação %d na classe %s: obteve %d, esperado %d",
					i, nome, uf.Count(), countEsperado)
			}
		}

		// valida a conexão dos pares
		// sorteia 5 pares aleatórios para testar
		for j := 0; j < 5; j++ {
			// pegam objetos de 0 a N -1
			a := r.Intn(n)
			b := r.Intn(n)

			var connectedEsperado bool
			primeira = true
			// percorre cada implementação
			for nome, uf := range instancias {
				// se o par sorteado ta conectado, tem que voltar true se não false
				res := uf.Connected(a, b)
				if primeira {
					connectedEsperado = res
					primeira = false
				} else if res != connectedEsperado {
					t.Fatalf("Erro no Connected(%d, %d) na operação %d em %s: obteve %v, esperado %v",
						a, b, i, nome, res, connectedEsperado)
				}
			}
		}
	}
}
