package tabela

// Insere associa o valor v à chave k. Se a chave já está na tabela, só troca
// o valor, e o número de pares não muda.
//
// Se a chave é nova, ela vai para a primeira posição vazia a partir de h(k),
// andando uma posição por vez e dando a volta no fim do vetor, como na Busca.
// Depois de guardar uma chave nova, se o fator de carga passou de 0,7, a
// tabela cresce para proximoPrimo(2*m) posições. Em inteiros, a condição é
// 10*n > 7*m.
func (t *Tabela) Insere(k, v int) {
	p := t.hash(k)
	achou := false
	for t.ocupada[p] {
		if t.chaves[p] == k {
			t.valores[p] = v
			achou = true
		}
		p = (p + 1) % t.m()
	}
	if !achou {
		t.chaves[p] = k
		t.valores[p] = v
		t.ocupada[p] = true
		t.n++
	}

	if 10*(t.Len()) > 7*(t.m()) {
		t.redimensiona(proximoPrimo(2 * t.m()))
	}
}

// 
// redimensiona troca o vetor por um de m posições e reinsere todos os pares.
// A posição de cada chave muda, porque h(k) = k mod m depende de m.
//
// Sugestão: crie uma tabela com Nova(m), insira nela os pares das posições
// ocupadas e copie o resultado com *t = *nova, como no encadeamento.go da
// aula de tabelas de dispersão.
func (t *Tabela) redimensiona(m int) {
	nova := Nova(m)
	for p:= 0; p < t.m(); p++ {
		if t.ocupada[p] {
			k := t.chaves[p]
			v := t.valores[p]
			
			ip := nova.hash(k)
			for nova.ocupada[ip] {
				ip = (ip + 1) % nova.m()
			}
			nova.chaves[ip] = k
			nova.valores[ip] = v
			nova.ocupada[ip] = true
			nova.n++
		}
	}
	*t = *nova
}

// Remove tira a chave k da tabela e devolve true, ou devolve false se ela
// não estava lá.
//
// Esvaziar a posição da chave não basta: uma chave que passou por ela ao ser
// inserida ficaria inalcançável, porque a Busca para na primeira posição
// vazia. A remoção percorre as posições seguintes ao buraco até achar uma
// vazia. Cada chave encontrada no caminho que pode ocupar o buraco desce para
// ele, e o buraco passa para a posição que ela deixou.
//
// A chave na posição j, com h(k) = h, pode descer para o buraco i quando o
// buraco está no caminho que a sondagem fez de h até j. Use a distancia de
// tabela.go para escrever essa condição sem tratar a volta no fim do vetor
// como caso à parte.
func (t *Tabela) Remove(k int) bool {
	// procura o k na tabela
	i := t.hash(k)
	for t.ocupada[i] {
		if t.chaves[i] == k {
			break
		}
		i = (i + 1) % t.m()
	}

	// nao achou
	if !t.ocupada[i] {
		return false
	}

	// esvazia a pos
	t.ocupada[i] = false
	t.n--


	j := (i + 1) % t.m()
	for t.ocupada[j] {
		h := t.hash(t.chaves[j]) // posicao de origem da chave que o j ta

		// ve se a distancia de i (onde esvaziou) até j (onde está iterando agora)
		// é menor ou igual a distancia de origem até j
		if t.distancia(i, j) <= t.distancia(h, j) {
			t.chaves[i] = t.chaves[j]
			t.valores[i] = t.valores[j]
			t.ocupada[i] = true
			t.ocupada[j] = false
			i = j // o novo buraco vira j
		}

		j = (j + 1) % t.m()
	}

	return true
}
