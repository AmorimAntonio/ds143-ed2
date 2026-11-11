// Bot ingênuo: ponto de partida do trabalho prático.
//
// Este bot implementa o protocolo inteiro e joga de forma simples:
//
//   - guarda os chamados abertos num slice;
//   - a cada turno, cada ambulância livre procura, por varredura linear,
//     o chamado mais próximo e anda na direção dele;
//   - cada ambulância ocupada anda na direção do hospital mais próximo.
//
// Ele não verifica se o chamado ainda pode ser alcançado dentro do prazo,
// não impede que duas ambulâncias sigam para o mesmo chamado e ignora os
// bloqueios da etapa 2. No cenário de estresse, a varredura linear estoura
// o limite de tempo por turno.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"bot/tabela" // reutilizado da atividade 3
)


const TAM = 25 // tamanho das seções

type Chamado struct {
	ID    int
	X, Y  int
	Valor int
	Prazo int
	Morto bool // se ta ativo ou nao
	Reservado bool // se tem alguem indo atras ou n
	DistHosp int // distancia ate o hospital
}

type Ambulancia struct {
	ID    int
	X, Y  int
	Livre bool
}

type Hospital struct {
	X, Y int
}

var (
	leitor   *bufio.Scanner
	escritor *bufio.Writer

	indice *tabela.Tabela // indice hash
	hospitais []Hospital
	chamados  []Chamado

	totalTurnos int // para saber o numero de turnos total
	expiraEm [][]int // guarda os chamados prestes a expirar

	// variaveis do seccionamento
	largura int
	altura int
	colunasS int
	linhasS int

	// x = colunas da secao, y = linhas da secao, z = posicao do chamado especifico dentro da secao
	secoes [][][]int 
	
	// guarda a posição do chamado que a ambulancia do id do alvo ta indo (-1 se n tem alvo)
	alvo []int
)

// lerCampos lê a próxima linha da entrada e a divide em campos.
// Devolve nil quando a entrada termina.
func lerCampos() []string {
	if !leitor.Scan() {
		return nil
	}
	return strings.Fields(leitor.Text())
}

func inteiro(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic("número inválido na entrada: " + s)
	}
	return v
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func distancia(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

// direcao devolve o movimento que aproxima (x, y) de (alvoX, alvoY),
// andando primeiro na horizontal.
func direcao(x, y, alvoX, alvoY int) string {
	switch {
	case x < alvoX:
		return "L"
	case x > alvoX:
		return "O"
	case y < alvoY:
		return "S"
	case y > alvoY:
		return "N"
	}
	return "P"
}

// lerInicio lê a mensagem de início. Este bot só usa os hospitais; os
// bloqueios são lidos e descartados.
func lerInicio() {
	c := lerCampos() // INICIO largura altura turnos limite_turno limite_preparo equipe num_equipes
	totalTurnos = inteiro(c[3]) // guarda o total de turnos

	// guarda altura e largura em variaveis globais
	largura = inteiro(c[1]) 
	altura = inteiro(c[2])

	// calcula as secoes
	colunasS = largura / TAM
	linhasS  = altura/ TAM


	h := inteiro(lerCampos()[1])
	for i := 0; i < h; i++ {
		c := lerCampos()
		hospitais = append(hospitais, Hospital{inteiro(c[0]), inteiro(c[1])})
	}

	
	k := inteiro(lerCampos()[1])
	alvo = make([]int, k) // uma posicao no vetor p cada ambulancia
	for i := 0; i < k; i++ {
		lerCampos() // id x y
		alvo[i] = -1 // inicia as ambulancias sem alvo
	}
	b := inteiro(lerCampos()[1])
	for i := 0; i < b; i++ {
		lerCampos() // x1 y1 x2 y2
	}
	lerCampos() // FIM
}

// removerPorID procura o chamado no slice e o remove, trocando-o pelo
// último elemento. Custo: proporcional ao número de chamados abertos.
/*
func removerPorID(id int) {
	for i := range chamados {
		if chamados[i].ID == id {
			chamados[i] = chamados[len(chamados)-1]
			chamados = chamados[:len(chamados)-1]
			return
		}
	}
}
*/

// removerExpirados mantém só os chamados com prazo maior ou igual ao turno.
func removerExpirados(turno int) {
	t := turno - 1
	if t >= 1 {
		for _, pos := range expiraEm[t]  {
			chamados[pos].Morto = true
			indice.Remove(chamados[pos].ID)
		}
		expiraEm[t] = nil
	}
}

// retorna qual a secao que o chamado ta
func qualSecao(x, y int) (int, int) {
	sl := y/TAM
	sc := x/TAM
	return sl, sc
}

func nota(pos, d int) int{
	// explicando o calculo da nota
	// basicamente uma relação entre o valor do chamado e a distancia 
	// +1 no divisor é pra não ter div por zero e *10000 é pra n ter q lidar com casa decimal
	return chamados[pos].Valor * 10000 / (d + chamados[pos].DistHosp + 1)
}

// retorna se da tempo ou nao de chegar ao chamado alvo
func daTempo(x, y, pos, turno int) bool {
	d := distancia(x, y, chamados[pos].X, chamados[pos].Y)
	return d <= chamados[pos].Prazo-turno+1
}


// procura em uma secao e limpa os mortos, devolve o melhor resultado atualizado
func olharSecao(sl, sc, x, y, melhor, melhorNota, turno int) (int, int){
	secaoAtual := secoes[sl][sc]
	i := 0
	for i < len(secaoAtual) {
		pos := secaoAtual[i] // posicao do chamado
		if chamados[pos].Morto {
			secaoAtual[i] = secaoAtual[len(secaoAtual)-1] // troca o chamado atual com o ultimo
			secaoAtual = secaoAtual[:len(secaoAtual)-1] // remove o ultimo da lista
			continue
		} else {

			// ve se o chamado ja ta reservado por outra ambulancia e se da tempo
			if chamados[pos].Reservado || !daTempo(x, y, pos, turno) {
				i++
				continue 
			}

			d := distancia(x, y, chamados[pos].X, chamados[pos].Y)
			
			if nota(pos, d) > melhorNota { // se achou um chamado com melhor nota
				melhor = pos // atualiza o melhor alvo
				melhorNota = nota(pos, d) // atualzia a melhor nota	
			}
			i++
		}
	}
	secoes[sl][sc] = secaoAtual // atualiza a lista
	return melhor, melhorNota
}


func melhorChamado(x, y, turno int) int {
	sl0, sc0 := qualSecao(x, y)
	melhor := -1
	melhorNota := -1 // -1 pra iniciar

	for dl := -1; dl <= 1; dl++ {
		for dc := -1; dc <= 1; dc++ {
			sl := sl0 + dl
			sc := sc0 + dc
			if sl < 0 || sl >= linhasS || sc < 0 || sc >= colunasS {
				continue // adjacente alvo nao existe
			}
			melhor, melhorNota = olharSecao(sl, sc, x, y, melhor, melhorNota, turno)
		}
	}
	return melhor
}

// hospitalMaisProximo devolve o índice, em hospitais, do hospital mais
// próximo de (x, y). Olha todos: são no máximo 20.
func hospitalMaisProximo(x, y int) int {
	melhor := 0
	for i, h := range hospitais {
		if distancia(x, y, h.X, h.Y) < distancia(x, y, hospitais[melhor].X, hospitais[melhor].Y) {
			melhor = i
		}
	}
	return melhor
}

func jogarTurno(turno int, ambulancias []Ambulancia) {
	for _, a := range ambulancias {
		A := alvo[a.ID] // chamado alvo da ambulancia
		mov := "P"
		if a.Livre {
			// se o alvo ainda vale
			if A != -1 {
				d := distancia(a.X, a.Y, chamados[A].X, chamados[A].Y)
				// se o alvo nao faz sentido 
				if chamados[A].Morto || d > chamados[A].Prazo-turno+1 {
					chamados[A].Reservado = false // libera para outra ambulância
					A = -1 // fica sem alvo
				}
			}

			// se nao tem alvo
			if A == -1 {
				A = melhorChamado(a.X, a.Y, turno)
				if A != -1 {
					chamados[A].Reservado = true
				}
			}

			// guarda o alvo
			alvo[a.ID] = A

			// anda
			if A != -1 {
				mov = direcao(a.X, a.Y, chamados[A].X, chamados[A].Y)
			} else {
				h := hospitalMaisProximo(a.X, a.Y)
				mov = direcao(a.X, a.Y, hospitais[h].X, hospitais[h].Y)
			}
		
		} else {
			if A != -1 { // libera se tinha alvo
				chamados[A].Reservado = false
				alvo[a.ID] = -1
			}
			h := hospitalMaisProximo(a.X, a.Y)
			mov = direcao(a.X, a.Y, hospitais[h].X, hospitais[h].Y)
		}
		fmt.Fprintln(escritor, a.ID, mov)
	}
	fmt.Fprintln(escritor, "FIM")
	escritor.Flush()
}

func main() {
	leitor = bufio.NewScanner(os.Stdin)
	leitor.Buffer(make([]byte, 1<<16), 1<<20)
	escritor = bufio.NewWriter(os.Stdout)

	
	lerInicio()



	// criando a tabela hash pra armazenar os chamados
	indice = tabela.Nova(5000011) // primo grande pra nao precisar expandir durante a execucao
	expiraEm = make([][]int, totalTurnos+1) // inicia o vetor de chamados p vencer, de 0 ao total de turnos
	// inicia vetor de secoes
	secoes = make([][][]int, linhasS)
	for bl := 0; bl < linhasS; bl++ { // precisa do for por ter 3 dimensoes
		secoes[bl] = make([][]int, colunasS)
	}

	fmt.Fprintln(escritor, "PRONTO")
	escritor.Flush()



	for {
		c := lerCampos() // TURNO t
		if c == nil {
			return // fim da partida
		}
		turno := inteiro(c[1])

		n := inteiro(lerCampos()[1]) // CHAMADOS n
		for i := 0; i < n; i++ {
			c := lerCampos()
			posicao := len(chamados)
			chamados = append(chamados, Chamado{inteiro(c[0]), inteiro(c[1]), inteiro(c[2]), 
				inteiro(c[3]), inteiro(c[4]), false, false, 0})

			// calcula a distancia ate o hospital mais proximo
			h := hospitalMaisProximo(chamados[posicao].X, chamados[posicao].Y)
			// passos do chamado ate o hospital
			chamados[posicao].DistHosp = distancia(chamados[posicao].X, chamados[posicao].Y, 
				hospitais[h].X, hospitais[h].Y)

			indice.Insere(inteiro(c[0]), posicao) // adiicona na tabela hash

			// guarda na lista de prestes a expirar
			prazo := inteiro(c[4])
			if prazo <= totalTurnos {
				expiraEm[prazo] = append(expiraEm[prazo], posicao)
			}

			// descobre qual secao ta o chamado e insere na matriz
			sl, sc := qualSecao(chamados[posicao].X, chamados[posicao].Y)
			secoes[sl][sc] = append(secoes[sl][sc], posicao)
		}

		r := inteiro(lerCampos()[1]) // RECOLHIDOS r
		for i := 0; i < r; i++ {
			id := inteiro(lerCampos()[0])
			posicao, achou := indice.Busca(id)
			if achou {
				chamados[posicao].Morto = true
				indice.Remove(id)
			}
		}
		removerExpirados(turno)

		k := inteiro(lerCampos()[1]) // AMBULANCIAS k
		ambulancias := make([]Ambulancia, k)
		for i := 0; i < k; i++ {
			c := lerCampos()
			ambulancias[i] = Ambulancia{inteiro(c[0]), inteiro(c[1]), inteiro(c[2]), c[3] == "LIVRE"}
		}

		e := inteiro(lerCampos()[1]) // PLACAR e
		for i := 0; i < e; i++ {
			lerCampos()
		}
		lerCampos() // FIM

		jogarTurno(turno, ambulancias)
	}
}
