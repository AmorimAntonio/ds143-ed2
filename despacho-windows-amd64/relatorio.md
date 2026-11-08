# Relatório de projeto: Antonio Amorim e Douglas Puchalski


## Remover um chamado recolhido sabendo só o id

Estrutura: tabela hash da atividade 3. Criada no preparo com 5000011 posições, para não precisar crescer durante a partida.

Custo por turno: O(1) para cada chamado recolhido.

Alternativa descartada: 

## Descartar os chamados expirados

Estrutura: vetor `expiraEm`, com uma posição por turno. `expiraEm[p]` é a lista dos chamados com o prazo  no turno p. No turno t, só percorre a lista de t − 1: as fichas viram `Morto`, saem da hash, e a lista é liberada.

Custo por turno: O(n de chamados que venceram naquele turno).

Alternativa descartada: 

## Achar os chamados próximos de uma ambulância

Estrutura: grade de seções, uma matriz secoes[linha][coluna] que divide o mapa em quadrados de TAM × TAM células (TAM = 50).
Cada seção guarda a lista das posições, no vetor de chamados, de cada chamado dentro dela. A seção de um chamado sai  da coordenada y/TAM, x/TAM. Para achar o chamado mais próximo, a ambulância olha só a própria seção e as vizinhas, e não todos os chamados abertos.

## Evitar que duas ambulâncias sigam para o mesmo chamado

Ainda nao fiz.

## Achar o hospital mais próximo de uma ambulância ocupada

Mantive como ta ainda.
