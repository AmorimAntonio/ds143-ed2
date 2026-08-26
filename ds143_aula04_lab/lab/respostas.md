# Aula 04 - Respostas

Nome:

Ambiente de medição (processador, memória, sistema operacional, versão do Go):

## Parte 1

Para cada função corrigida, uma linha dizendo qual era o defeito.

| Função | Defeito |
|---|---|
| `BuscaBinaria` | o loop não verificava se o lo e o hi fossem iguais, então não entrava no loop no caso em que encontrava o valor |
| `InsertionSort` | o loop interno não ia até o fim do vetor, bastou adicionar um j>0 |
| `Merge` | faltou o for para acrescentar os restantes do vetor b caso ele seja maior |
| `Particiona` | faltava mover o pivo no final do loop |

## Parte 2

Saída de `./benchmark` (entrada aleatória):

```
$ go build -o benchmark ./cmd/benchmark
./benchmark             
Tempos em segundos. Entrada aleatoria. ERRO significa vetor nao ordenado ao final.

Tabela 1: todos os algoritmos

algoritmo            5000 razao    10000 razao    20000 razao    40000 razao
SelectionSort      0.0113         0.0515   4.6   0.1993   3.9   0.8535   4.3
InsertionSort      0.0068         0.0281   4.2   0.1179   4.2   0.4741   4.0
MergeSort          0.0000         0.0010         0.0025   2.5   0.0067   2.7
QuickSort          0.0000         0.0000         0.0000         0.0010      
QuickSort emb.     0.0000         0.0000         0.0000         0.0015      
sort.Ints          0.0000         0.0000         0.0000         0.0015      

Tabela 2: apenas os algoritmos O(n log n)

algoritmo          125000 razao   250000 razao   500000 razao  1000000 razao
MergeSort          0.0226         0.0469   2.1   0.1128   2.4   0.2131   1.9
QuickSort          0.0068         0.0153   2.2   0.0345   2.3   0.0709   2.1
QuickSort emb.     0.0072         0.0181   2.5   0.0387   2.1   0.0814   2.1
sort.Ints          0.0072         0.0161   2.2   0.0344   2.1   0.0744   2.2

```

Saída de `./benchmark -ordenada`:

```
$ ./benchmark -ordenada 
Tempos em segundos. Entrada ja ordenada. ERRO significa vetor nao ordenado ao final.

Tabela 1: todos os algoritmos

algoritmo            5000 razao    10000 razao    20000 razao    40000 razao
SelectionSort      0.0038         0.0455  12.0   0.1906   4.2   0.7833   4.1
InsertionSort      0.0000         0.0000         0.0000         0.0000      
MergeSort          0.0000         0.0000         0.0000         0.0000      
QuickSort          0.0048         0.0317   6.6   0.1650   5.2   0.6475   3.9
QuickSort emb.     0.0000         0.0000         0.0000         0.0000      
sort.Ints          0.0000         0.0000         0.0000         0.0000      

A tabela 2 nao roda com -ordenada: o QuickSort sem embaralhamento
levaria tempo demais nesses tamanhos. Explique por que.
```
Porque o quicksort faria partições mesmo já estando ordenado, o que resultaria em várias operações com o pivô, independente do vetor estar ou não ordenado.


1. Pelas razões da tabela de entrada aleatória, quais algoritmos são O(n²) e
   quais são O(n log n)? Cite as razões que sustentam a resposta.
   - Selection sort (4.3) e insertion sort (4.0) são O(n²)
   - Merge sort (1.9),  QuickSort (2.1), QuickSort (2.1), sorts.Ints (2.2)
   Respostas para os algoritmos de 40k


2. Na entrada já ordenada, o `InsertionSort` e o `QuickSort` mudam de
   comportamento em direções opostas. Explique cada um.

## Parte 3

Saída de `./consultas`:

```
cole aqui
```

3. A partir de quantas consultas compensa ordenar o vetor antes de buscar?
   Compare o valor medido com o que a análise prevê.
