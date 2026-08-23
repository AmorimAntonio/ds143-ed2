# Relatório de Medição Empírica — Union-Find

## Ambiente de Medição
* **Processador:** Intel(R) Core(TM) i5-9300H CPU @ 2.40GHz (2.40 GHz)
* **Memória RAM:** 16 GB
* **Sistema Operacional:** Windows 11 
* **Versão do Go:** go 1.26.5 windows/amd64

## Tabela de Resultados

| n | Quick-Find | Razão | Quick-Union | Razão | Weighted | Razão | Weighted+PC | Razão |
| :--- | :---: | :---: | :---: | :---: | :---: | :---: | :---: | :---: |
| **10.000** | 0.084 | - | 0.057 | - | 0.000 | - | 0.002 | - |
| **20.000** | 0.404 | 4.8 | 0.316 | 5.5 | 0.006 | - | 0.001 | 0.3 |
| **40.000** | 1.357 | 3.4 | 1.523 | 4.8 | 0.002 | 0.4 | 0.001 | 2.0 |
| **80.000** | 5.333 | 3.9 | 12.022 | 7.9 | 0.003 | 1.5 | 0.005 | 4.6 |
| **160.000** | 21.763 | 4.1 | 66.287 | 5.5 | 0.011 | 3.7 | 0.008 | 1.7 |

## Perguntas e Respostas

**1. As razões medidas confirmam a ordem de crescimento prevista na seção 14 da aula anterior para cada implementação? Cite os números que você obteve.**

Sim. O Quick-Find manteve razão próxima a 4 (variou de 3.4 a 4.8), confirmando a ordem O(n^2). O Quick-Union também exibiu comportamento quadrático com variações maiores (razões de 4.8 a 7.9). O Weighted e o Weighted+PC foram muito mais rápidos (menos de milissegundos), com razões oscilando perto de 2 (0.4 a 3.7 para o Weighted e 0.3 a 4.6 para o Weighted+PC), confirmando o crescimento sub-quadrático.

**2. O Weighted e O Weighted+PC marcam tempos praticamente iguais nesses tamanhos. Isso contraria a proposição da seção 15? Diga o que precisaria mudar na medição para a diferença entre os dois aparecer.**

Não contraria. Para n <= 160.000, a altura máxima da árvore no Weighted é de no máximo 18 passos ($\log_2$ n), tornando o ganho da compressão de caminho imperceptível em milissegundos. Para a diferença se destacar, seria necessário aumentar o valor de $n$ significativamente (para a casa dos milhões) ou realizar uma grande quantidade de consultas via `Find`.

## Uso de IA
Utilizei o Gemini em alguns momentos do trabalho, sendo estes:
* Esclarecer dúvidas conceituais sobre compressão de caminho e balanceamento.
* Explicar o funcionamento dos códigos do esqueleto disponibilizados, uma vez que não estou habitiuado com GoLang.
* Auxiliar na estruturação do teste do oráculo cruzado.
* Auxílio na formatação e revisão textual do relatório final.