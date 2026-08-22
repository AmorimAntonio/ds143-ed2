package unionfind

// WeightedPathCompression é o Weighted Quick-Union com compressão de caminho:
// a cada subida até a raiz, os nós visitados são aproximados dela.
//
// É a única variante sem código pronto no material da disciplina. A seção 15 da
// aula de Análise de Algoritmos descreve a ideia e mostra o laço do Find que
// aponta cada nó para o avô. O resto parte do weightedquickunion.go, neste
// mesmo pacote.
type WeightedPathCompression struct {
	parent []int
	size   []int // size[i] = número de objetos na subárvore de raiz i
	count  int
}

var _ UF = (*WeightedPathCompression)(nil)

// NewWeightedPathCompression cria a estrutura com n objetos.
func NewWeightedPathCompression(n int) *WeightedPathCompression {
	uf := &WeightedPathCompression{
		parent: make([]int, n),
		size:   make([]int, n),
		count:  n,
	}
	for i := range uf.parent {
		uf.parent[i] = i
		uf.size[i] = 1
	}
	return uf
}

// Find sobe até a raiz e, no caminho, encurta a árvore.
func (uf *WeightedPathCompression) Find(p int) int {
	// lógica retirada da seção 15
	for p != uf.parent[p] {
		uf.parent[p] = uf.parent[uf.parent[p]] // aponta para o avô
		p = uf.parent[p]
	}
	return p
}

func (uf *WeightedPathCompression) Connected(p, q int) bool {
	return uf.Find(p) == uf.Find(q)
}

func (uf *WeightedPathCompression) Union(p, q int) {
	raizP, raizQ := uf.Find(p), uf.Find(q)
	if raizP == raizQ {
		return
	}
	if uf.size[raizP] < uf.size[raizQ] {
		uf.parent[raizP] = raizQ
		uf.size[raizQ] += uf.size[raizP]
	} else {
		uf.parent[raizQ] = raizP
		uf.size[raizP] += uf.size[raizQ]
	}
	uf.count--
}

func (uf *WeightedPathCompression) Count() int {
	return uf.count
}
