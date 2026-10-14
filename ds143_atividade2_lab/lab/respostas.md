# Atividade Avaliativa 2 - Respostas

Nome: Antonio Amorim e Douglas Puchalski

GRR: 2025289 e 20250662

Se a atividade foi feita em dupla, escreva aqui os dois nomes e os dois GRR. As
duas pessoas enviam o mesmo `.zip` na UFPR Virtual.

Ambiente (sistema operacional e versão do Go):

Saída de `go test ./...` na pasta que você vai enviar:

```
cole aqui
```

## Perguntas

1. Em uma frase, como a sua `EhABB` rejeita a árvore que tem o 60 na subárvore
   esquerda do 50, sendo que o 60 é maior que o pai 30?

   R: Porque atribuimos o valor maximo como referencia na chamada de função e - nessa situação - ele é 50 pois o atribuimos ao fazer a chamada pela esquerda, nossa função verifica se o valor do nó é igual ou maior que o máximo para alertar caso não seja uma ABB, o que é o caso.

2. Em uma frase, o que a sua `Sucessor` guarda durante a descida, e por que isso
   basta quando `v` não está na árvore.

   R: A sucessor guarda o valor do NoAtual antes de ir para a esquerda, porque ela vai indo para a direita até achar um maior, quanda acha vai indo para a esquerda e, caso ache um menor que o próprio V, o valor do ultimo NoAtual maior está salvo (e é o próximo maior).


3. A `VerificaAVL` confere duas condições independentes. Dê um exemplo de árvore
   que satisfaz a primeira e falha na segunda, desenhando-a com o campo `Alt` de
   cada nó.

   ```
                50 (2)
               /  
              20(1)
            /   
           2(0)  
   ```
   Esse é um caso válido, porque a altura do primeiro nó está certa, mas o fator de balanceamento ta maior que 1

## Uso de IA

Diga se usou ferramentas de IA generativa nesta atividade e em que partes. Não
influencia a nota.

   R: Utilizamos IA (Gemini) como auxílio para tirar dúvidas sobre a sintaxe do Go (como manipulação de ponteiros e retornos múltiplos)