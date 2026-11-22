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
const DISTMAX = 100 // distancia maxima do percurso em larg
const MAXBUSCAS = 30 // 30 buscas por turno no maximo p n estourar


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
	rodada int // qual rodada é

	// x = colunas da secao, y = linhas da secao, z = posicao do chamado especifico dentro da secao
	secoes [][][]int 
	
	// guarda a posição do chamado que a ambulancia do id do alvo ta indo (-1 se n tem alvo)
	alvo []int



	// adicionadas para a Etapa 2:
	bloqueado [] bool // ve se é parede ou n
	distHosp [] int // guarda a distancia da celula ate o hosp (-1 se n visitou ainda)
	passoHosp []byte // movimento que tem que fazer p ir ate o hosp
	// N, S, L e O ou P se ta no proprio hospital

	filaAmb []int // a fila do percurso da ambulancia
	visitado []int // rodada q a celula foi visitada
	pai []int // celula usada pra chegar nessa
	pilha [][]string // pilha de movimentos por ambulancia
	desistiu []int // o último chamado que cada ambulância abandonou (-1 = nenhum).
)

// ===================== leitura da entrada =====================

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

// lerInicio lê a mensagem de início. 
func lerInicio() {
	c := lerCampos() // INICIO largura altura turnos limite_turno limite_preparo equipe num_equipes
	totalTurnos = inteiro(c[3]) // guarda o total de turnos

	// guarda altura e largura em variaveis globais
	largura = inteiro(c[1]) 
	altura = inteiro(c[2])

	// vetor com as paredes
	bloqueado = make([]bool, largura*altura)

	// calcula as secoes
	colunasS = largura / TAM
	linhasS  = altura/ TAM


	h := inteiro(lerCampos()[1])
	for i := 0; i < h; i++ {
		c := lerCampos()
		hospitais = append(hospitais, Hospital{inteiro(c[0]), inteiro(c[1])})
	}

	
	k := inteiro(lerCampos()[1])
	alvo = make([]int, k)
	pilha = make([][]string, k)
	desistiu = make([]int, k)
	for i := 0; i < k; i++ {
		lerCampos()
		alvo[i] = -1
		desistiu[i] = -1
	}


	// lendo os bloqueios
	b := inteiro(lerCampos()[1])
	for i := 0; i < b; i++ {
		c := lerCampos()

		// pega os cantos dos bloqueios
		x1 := inteiro(c[0])
		y1 := inteiro(c[1])
		x2 := inteiro(c[2])
		y2 := inteiro(c[3])
		// marca a celula de dentro como bloqueadas
		for y := y1; y <= y2; y++ {
			for x := x1; x <= x2; x++ {
				bloqueado[cel(x, y)] = true
			}
		}
	}
	lerCampos() // FIM
}

func percursoLarguraHospitais() {
	n := largura * altura

	distHosp = make([]int, n)
	for i := 0; i < n; i++ {
		distHosp[i] = -1 // nao visitada
	}
	passoHosp = make([]byte, n)

	// fila com todas as celulas
	fila := make([]int, n)
	ini, fim := 0, 0

	// hospitais na fila
	for _, h := range hospitais {
		c := cel(h.X, h.Y)
		distHosp[c] = 0
		passoHosp[c] = 'P' // hospital = hospital logicamente
		fila[fim] = c
		fim++
	}


	for ini < fim {
		// pega a proxima celula da fila
		c := fila[ini]
		ini++
		x := c % largura
		y := c / largura

		// vizinho a leste = pra voltar ate c, anda pra oeste
		if x+1 < largura {
			nc := cel(x+1, y)
			if !bloqueado[nc] && distHosp[nc] == -1 {
				distHosp[nc] = distHosp[c] + 1
				passoHosp[nc] = 'O'
				fila[fim] = nc
				fim++
			}
		}

		// vizinho a oeste = pra voltar ate c, anda pra leste
		if x-1 >= 0 {
			nc := cel(x-1, y)
			if !bloqueado[nc] && distHosp[nc] == -1 {
				distHosp[nc] = distHosp[c] + 1
				passoHosp[nc] = 'L'
				fila[fim] = nc
				fim++
			}
		}

		// vizinho ao sul = pra voltar ate c, anda pro norte
		if y+1 < altura {
			nc := cel(x, y+1)
			if !bloqueado[nc] && distHosp[nc] == -1 {
				distHosp[nc] = distHosp[c] + 1
				passoHosp[nc] = 'N'
				fila[fim] = nc
				fim++
			}
		}

		// vizinho ao norte = pra voltar ate c, anda pro sul
		if y-1 >= 0 {
			nc := cel(x, y-1)
			if !bloqueado[nc] && distHosp[nc] == -1 {
				distHosp[nc] = distHosp[c] + 1
				passoHosp[nc] = 'S'
				fila[fim] = nc
				fim++
			}
		}
	}
}


func caminho(id, x, y, pos, turno int) bool {

	// estabelece limite, limitado ao maximo da var global constante
	limite := chamados[pos].Prazo - turno + 1
	if limite >  DISTMAX {
		limite = DISTMAX
	}

	// preparando a busca
	rodada++ // muda p prox rodada 
	c0 := cel(x, y) // celula da ambulancia
	alvoC := cel(chamados[pos].X, chamados[pos].Y) // celulla chamado
	visitado[c0] = rodada 
	filaAmb[0] = c0
	ini := 0 // quem sai da fila
	fim := 1; // quem entra na fila
	d := 0 // distancia ate a ambulancia (0 pq a amabulancia ta aqui)
	fimNivel := 1; 
	achou := c0 == alvoC // se ja ta em cima do chamado, é true
	
	// percorre em largura ate achar o chamado ou passar do limite
	for ini < fim && !achou {
		// terminou a camada atual, aumenta a distancia
		if ini == fimNivel {
			d++
			fimNivel = fim
			// se o vizinho ta fora do limite
			if d >= limite {
				break 
			}
		}

		// pega as coord da prox cel na fila
		c := filaAmb[ini]
		ini++
		x := c % largura
		y := c / largura

		// leste
		if x+1 < largura {
			nc := cel(x+1, y)
			if !bloqueado[nc] && visitado[nc] != rodada {
				visitado[nc] = rodada
				pai[nc] = c
				filaAmb[fim] = nc
				fim++
				if nc == alvoC {
					achou = true
				}
			}
		}

		// oeste
		if x-1 >= 0 {
			nc := cel(x-1, y)
			if !bloqueado[nc] && visitado[nc] != rodada {
				visitado[nc] = rodada
				pai[nc] = c
				filaAmb[fim] = nc
				fim++
				if nc == alvoC {
					achou = true
				}
			}
		}

		// sul
		if y+1 < altura {
			nc := cel(x, y+1)
			if !bloqueado[nc] && visitado[nc] != rodada {
				visitado[nc] = rodada
				pai[nc] = c
				filaAmb[fim] = nc
				fim++
				if nc == alvoC {
					achou = true
				}
			}
		}

		// norte
		if y-1 >= 0 {
			nc := cel(x, y-1)
			if !bloqueado[nc] && visitado[nc] != rodada {
				visitado[nc] = rodada
				pai[nc] = c
				filaAmb[fim] = nc
				fim++
				if nc == alvoC {
					achou = true
				}
			}
		}
	}

	if !achou {
		return false
	}

	// monta a pilha voltando do chamado ate a ambulancia
	// o ultimo empilhado é o primeiro movimento
	pilha[id] = pilha[id][:0]
	c := alvoC
	for c != c0 {
		p := pai[c]
		pilha[id] = append(pilha[id], direcao(p%largura, p/largura, c%largura, c/largura))
		c = p
	}
	return true
}


// ===================== funcoes auxiliares =====================

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func distancia(x1, y1, x2, y2 int) int {
	return abs(x1-x2) + abs(y1-y2)
}

// função pra deixar uma coordenada em uma dimensão
func cel(x, y int) int { 
	return y*largura + x 
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

// ===================== controle dos chamados =====================

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

func nota(pos, d int) int{
	// explicando o calculo da nota
	// basicamente uma relação entre o valor do chamado e a distancia 
	// +1 no divisor é pra não ter div por zero e *10000 é pra n ter q lidar com casa decimal
	// formula original chamados[pos].Valor * 10000 / (d + chamados[pos].DistHosp + 1)
	// alterada depois de mts testes para ir achando o ponto ideal
	return (chamados[pos].Valor + 10) * 20000 / (3*d + 2*chamados[pos].DistHosp + 2)
}

// retorna se da tempo ou nao de chegar ao chamado alvo
func daTempo(x, y, pos, turno int) bool {
	d := distancia(x, y, chamados[pos].X, chamados[pos].Y)
	return d <= chamados[pos].Prazo-turno+1
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

// ===================== busca nas secoes =====================

// retorna qual a secao que o chamado ta
func qualSecao(x, y int) (int, int) {
	sl := y/TAM
	sc := x/TAM
	return sl, sc
}

// procura em uma secao e limpa os mortos, devolve o melhor resultado atualizado
func olharSecao(id, sl, sc, x, y, melhor, melhorNota, turno int) (int, int) {
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
			if chamados[pos].Reservado || pos == desistiu[id] || !daTempo(x, y, pos, turno) {
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


func melhorChamado(id, x, y, turno int) int {
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
			melhor, melhorNota = olharSecao(id, sl, sc, x, y, melhor, melhorNota, turno)
		}
	}
	return melhor
}


// ===================== turno e main =====================

func jogarTurno(turno int, ambulancias []Ambulancia) {

	buscas := 0

	// p cada ambulancia
	for _, a := range ambulancias {
		A := alvo[a.ID] // chamado alvo da ambulancia
		mov := "P" // parado por padrão

		// para as livres
		if a.Livre {

			// se o alvo ainda vale
			if A != -1 {
				d := len(pilha[a.ID])
				// ou seja, se morreu ou se n tem mais como chegar a tempo
				if chamados[A].Morto || d > chamados[A].Prazo-turno+1 {
					chamados[A].Reservado = false // libera para outra ambulância
					pilha[a.ID] = pilha[a.ID][:0] // libera a pilha tb
					A = -1 // fica sem alvo
				}
			}

			// se nao tem alvo
			if A == -1 && buscas < MAXBUSCAS{
				buscas++
				// procura outro
				A = melhorChamado(a.ID, a.X, a.Y, turno)
				if A != -1 {
					if caminho(a.ID, a.X, a.Y, A, turno) {
						chamados[A].Reservado = true
					} else {
						desistiu[a.ID] = A // nao da tempo pelas ruas, nao escolhe de novo
						A = -1
					}
				}
			}

			// guarda o alvo
			alvo[a.ID] = A

			// anda
			if A != -1 {
				if len(pilha[a.ID]) > 0 {
					topo := len(pilha[a.ID]) - 1
					mov = pilha[a.ID][topo]
					pilha[a.ID] = pilha[a.ID][:topo]
				}
			} else {
				// hospital mais proximo
				mov = string(passoHosp[cel(a.X, a.Y)])
			}
		

		// para as ocupadas:
		} else {
			if A != -1 { // libera se tinha alvo
				chamados[A].Reservado = false
				alvo[a.ID] = -1
				pilha[a.ID] = pilha[a.ID][:0]
			}
			mov = string(passoHosp[cel(a.X, a.Y)])
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
	indice = tabela.Nova(2000003) // primo grande pra nao precisar expandir durante a execucao
	expiraEm = make([][]int, totalTurnos+1) // inicia o vetor de chamados p vencer, de 0 ao total de turnos
	// inicia vetor de secoes
	secoes = make([][][]int, linhasS)
	for bl := 0; bl < linhasS; bl++ { // precisa do for por ter 3 dimensoes
		secoes[bl] = make([][]int, colunasS)
	}

	// função que marca o caminho das celulas
	percursoLarguraHospitais()

	filaAmb = make([]int, largura*altura)
	visitado = make([]int, largura*altura)
	pai = make([]int, largura*altura)



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

			
			// ve qnt falta pra chegar no hospital mais proximo
			chamados[posicao].DistHosp = distHosp[cel(chamados[posicao].X, chamados[posicao].Y)]

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
