# Relatório de projeto: Antonio Amorim (GRR20255289) e Douglas Puchalski (GRR20250662)

## Remover um chamado recolhido sabendo só o id

Estrutura: tabela hash da atividade 3. Criada no preparo com 2000003 (numero primo bem grande) posições, para não precisar crescer durante a partida. id = posição do chamado no vetor. Quando chega um recolhido, acha pela hash e marca como Morto.

Custo por turno: O(1) para cada chamado recolhido.

Alternativa descartada: procurar o id varrendo o vetor de chamados (como o ingênuo fazia), porque é O(n) por recolhido e no estresse chegam 1000 chamados por turno, ai estoura o tempo.

## Descartar os chamados expirados

Estrutura: vetor expiraEm, com uma posição por turno. expiraEm[p] é a lista dos chamados com o prazo no turno p. No turno t, só percorre a lista de t − 1, os chamados viram Morto, saem da hash, e a lista é liberada.

Custo por turno: O(n de chamados que venceram naquele turno).

Alternativa descartada: varrer todos os chamados todo turno vendo o prazo, porque é O(n) toda rodada mesmo quando nenhum chamado venceu.

## Achar os chamados próximos de uma ambulância

Estrutura: matriz de seções, divide o mapa em quadrados de 25 × 25. Cada seção guarda os chamados que estão dentro dela. A ambulância só olha a própria seção e asvizinhas. Entre esses, escolhe o de maior nota, calculo da nota: (valor + 10) / (1,5 × distância até o chamado + distância até o hospital + 1). O +10 é pra diminuir a diferença entre chamado comum e grave, e o 1,5 faz preferir os mais perto, para pegar mais chamados em geral. Testamos varias formulas e essa deu uns +8% contra o bot da etapa 2, a antiga dava +3%.

Custo por busca: O(chamados nas 9 seções).

Alternativa descartada: Seção maior (tipo 50) ou busca em mais seções tambem foram descartadas, porque a ambulância ia buscar chamado longe demais e a pontuação caiu.

## Evitar que duas ambulâncias sigam para o mesmo chamado

Estrutura: campo Reservado no chamado e um vetor alvo com o chamado de cada ambulância. Quando a ambulância escolhe um chamado, marca como reservado, e as outras pulam ele. Ela só procura outro quando o dela morre ou não dá mais tempo. O vetor desistiu guarda o ultimo chamado que ela largou, pra não escolher ele de novo no turno seguinte.

Custo por turno: O(1) por ambulância. A busca só roda quando ela perde o alvo, e no maximo 30 buscas por turno (MAXBUSCAS).

Alternativa descartada: escolher o alvo de novo todo turno, porque o tempo ficou mais de 6 vezes o do bot de referencia e as ambulâncias ficavam trocando de alvo no meio do caminho.

## Achar o hospital mais próximo de uma ambulância ocupada

Estrutura: percurso em largura saindo de todos os hospitais ao mesmo tempo, feito uma vez no preparo. Pra cada célula guarda distHosp (distancia real até o hospital mais proximo) e passoHosp que é o movimento N, S, L, O ou P que leva até ele, em vetores de byte e int com uma posição por célula. A ambulância ocupada só lê o passoHosp da célula onde está. O DistHosp do chamado vem do mesmo vetor quando ele chega.

Custo por turno: O(1) por ambulância ocupada. No preparo, O(células do mapa).

Alternativa descartada: olhar os 20 hospitais com distancia de Manhattan, porque com bloqueios a ambulância andava em linha reta e batia em parede. 
## Calcular caminhos a cada turno ou no preparo

Estrutura: o caminho até o hospital é calculado no preparo, porque os hospitais não mudam. O caminho até o chamado é calculado na hora que a ambulância escolhe o alvo: um percurso em largura da ambulância até o chamado guardando o pai de cada célula, e depois volta pelos pais montando uma pilha de movimentos. A cada turno ela tira o movimento do topo da pilha. Se o tamanho da pilha passar do que o prazo deixa, ela desiste. O vetor visitado guarda o numero da rodada, então não precisa limpar ele a cada busca.

Custo por turno: no maximo 30 percursos, cada uma limitada a 100 passos. Andar é O(1), só tira da pilha.

Alternativa descartada: Recalcular o percurso de toda ambulância todo turno tambem foi descartado, porque é caro e o caminho não muda se o alvo é o mesmo.


# possiveis melhorias

AVL em cada seção que armazena a pontuação dos chamados daquela seção.