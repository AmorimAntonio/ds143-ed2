# Atividade Avaliativa 3 - Respostas

Nome: Antonio Amorim e Douglas Puchalski

GRR: 20255289 e 20250662

Se a atividade foi feita em dupla, escreva aqui os dois nomes e os dois GRR. As
duas pessoas enviam o mesmo `.zip` na UFPR Virtual.

Ambiente (sistema operacional, processador e versão do Go):
- Windows 11
- Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz (2.40 GHz)
- go1.26.5 windows/amd64


Saída de `go test ./...` na pasta que você vai enviar:

```
ok      ds143/atividade3/tabela (cached)
```

## Medidas

Copie a coluna `ns/op` de cada linha da saída de
`go test -run '^$' -bench . ./tabela`, primeiro com `misturar = false` e depois
com `misturar = true`. Volte para `false` antes de enviar.

| Medida | Tabela, `false` | map, `false` | Tabela, `true` | map, `true` |
|---|---|---|---|---|
| Insere, aleatórias | 79.11 ns/op | 56.23 ns/op | 78.45 ns/op | 63.74 ns/op |
| Insere, sequenciais | 34.91 ns/op | 62.31 ns/op | 76.45 ns/op | 58.45 ns/op |
| Busca, aleatórias, presentes | 25.48 ns/op | 18.86 ns/op | 25.59 ns/op | 18.05 ns/op |
| Busca, aleatórias, ausentes | 38.67 ns/op | 30.70 ns/op | 37.48 ns/op | 27.89 ns/op |
| Busca, sequenciais, presentes | 12.90 ns/op | 19.50 ns/op | 25.35 ns/op | 18.70 ns/op |
| Janela | 28490 ns/op | 68.49 ns/op | 252.4 ns/op | 63.14 ns/op |


## Perguntas

1. Na sua `Remove`, depois que a chave removida deixa a posição vazia, as
   chaves seguintes podem descer para ela. Em que condição uma dessas chaves
   desce? Responda em função de três posições: a de origem da chave, a posição
   em que ela está e a posição vazia. No teste `TestRemoveVoltaCircular`, por
   que o 1 fica na posição 1 e o 10, que está depois dele, desce para a
   posição 0?
   
   
   **R:** A chave desce se a posição vazia está entre a posição de origem e a posição 
   que ela está, porque isso significa que a posição vazia ou é a posição de origem
   ou está mais perto da posição de origem do que a posição atual da chave.

   Por que o elemento removido está fora da distância entre o 1 e sua posição de origem,
   mas está dentro da distância entre o 10 e a sua posição de origem. 

2. Com `misturar = false`, a linha da janela da `Tabela` fica centenas de vezes
   mais lenta que a do `map`. Quantas posições a sua `Remove` examina, mais ou
   menos, a cada passo da janela, e por quê? Responda pela posição em que cada
   chave sequencial fica na tabela.

   **R:** A remove examina mais ou menos 1000 posições a cada passo, porque são 1000
   valores sequenciais, então é como se a tabela estivesse preenchida de forma contínua.
   Como ele remove a mais antiga, o espaço vazio fica lá e ele percorre toda a tabela até
   achar o primeiro espaço vazio, como ela está preenchida sem intervalos até a posição 1000
   ele faz 1000 examinações até achar o espaço vazio.

3. Com `misturar = true`, o que melhorou e o que piorou em relação a
   `misturar = false`? Num programa em que as chaves são identificadores
   sequenciais que entram e saem o tempo todo, qual das duas você usaria?

   **R:** O teste da janela melhorou muito após ativar o misturar = true, 
   ficando mais ou menos 100x mais rápido (porque as chaves agora ficam 
   espalhadas pela tabela em vez de formarem um bloco único contínuo, criando
   posições vazias entre os pequenos grupos). Porém, a inserção e 
   a busca de chaves sequenciais ficaram cerca de 2x mais lentas, devido ao 
   custo de processamento da função de mistura. Sobre qual eu usaria, depende, 
   mas no geral, em programa em que chaves mudam frequentemente, eu usaria misturar = true,
   pois o pequeno custo computacional do hash evita o custo de processamento alto
   de remoções em janela


## Uso de IA

Diga se usou ferramentas de IA generativa nesta atividade e em que partes. Não
influencia a nota.

Usei o gemini, me ajudou com a sintaxe do Go e à descobrir um jeito de percorrer 
o vetor circularmente (usando o mod)