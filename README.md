# rate-limiter-system-design

A tolerância a burst (ou burst tolerance) se refere à capacidade de um algoritmo de rate limiting permitir picos curtos e intensos de requisições em um curto espaço de tempo, sem violar o limite total.

📌 Definição rápida:
A tolerância a burst é a flexibilidade temporária que um sistema de rate limiting tem para aceitar múltiplas requisições de forma quase simultânea, desde que o total esteja dentro do permitido a médio prazo.

🧠 Exemplificando:
Suponha:
Limite: 10 requisições por minuto

Dois usuários enviam requisições:

Caso com tolerância a burst (ex: Token Bucket):
Um usuário envia 10 requisições de uma vez → 💚 Aceito (bucket cheio).

Depois, ele precisa esperar o refill dos tokens.

Caso sem tolerância a burst (ex: Leaking Bucket):
O sistema só permite, por exemplo, 1 requisição a cada 6 segundos.

Mesmo que ele tenha feito zero requisições antes, se tentar 10 de uma vez → ❌ Só a 1ª entra, as outras são rejeitadas.

📊 Comparação rápida:
Algoritmo Tolerância a Burst Comportamento
Token Bucket Alta Permite picos rápidos até encher o balde
Leaking Bucket Baixa Fluxo constante, limita picos abruptos
Fixed Window Média Depende do momento em que o "clock" reinicia
Sliding Window Média/Alta Mais equilibrado entre burst e suavização

✅ Quando usar alta tolerância a burst?
APIs públicas que devem ser rápidas com picos pequenos.

Sistemas tolerantes a carga variável.

Usuários legítimos que fazem ações em lote (ex: dashboards, scripts).

❌ Quando evitar?
Quando o sistema backend é sensível a carga repentina.

Quando você precisa garantir um ritmo constante e previsível (ex: filas de mensagens, robôs de consumo de dados).
