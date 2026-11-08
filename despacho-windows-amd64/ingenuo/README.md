# Relatório de projeto: Antonio Amorim e Douglas Puchalski

## Remover um chamado recolhido sabendo só o id

Estrutura: tabela hash da atividade 3. Criada no preparo com 5000011 (numero primo bem grande)posições, para não precisar crescer durante a partida. id = posição do chamado no vetor. Quando chega um recolhido, acha pela hash e marca como Morto.

Custo por turno: O(1) para cada chamado recolhido.

## Descartar os chamados expirados

Estrutura: vetor expiraEm, com uma posição por turno. expiraEm[p] é a lista dos chamados com o prazo no turno p. No turno t, só percorre a lista de t − 1, os chamados viram Morto, saem da hash, e a lista é liberada.

Custo por turno: O(n de chamados que venceram naquele turno).

## Achar os chamados próximos de uma ambulância

Estrutura: matriz de seções, divide o mapa em quadrados de 50 × 50. Cada seção guarda os chamados que estão dentro dela. A ambulância só olha a própria seção e as 8 vizinhas. Entre esses, escolhe o que dá mais pontos por turno gasto: valor / (distância até o chamado + distância até o hospital). Assim os graves têm prioridade.

Custo por busca: O(chamados nas 9 seções).

## Evitar que duas ambulâncias sigam para o mesmo chamado

Estrutura: campo Reservado no chamado e um vetor alvo com o chamado de cada ambulância. Quando a ambulância escolhe um chamado, marca como reservado, e as outras pulam ele. Ela só procura outro quando o dela morre ou não dá mais tempo.

Custo por turno: O(1) por ambulância. A busca só roda quando ela perde o alvo.

## Achar o hospital mais próximo de uma ambulância ocupada

Estrutura: olha todos os hospitais, que são no máximo 20. A distância do chamado até o hospital é calculada uma vez, quando ele chega, e fica guardada no DistHosp.

Custo por turno: O(20) por ambulância ocupada.


# possiveis melhorias

AVL em cada seção que armazena a pontuação dos chamados daquela seção. 