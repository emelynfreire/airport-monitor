# airport-monitor

Prática da arquitetura **Publish–Subscribe** usando **Apache Kafka** — exercício para disciplina de **Sistemas Distribuídos**.

Este repositório demonstra um cenário simples: um **AirportDatabase** (publisher) publica eventos de voos e três tótens (subscribers) consomem esses eventos:
- `FlightTotemArrival` — só chegadas (`arrival-flights`)
- `FlightTotemDeparture` — só partidas (`departure-flights`)
- `FlightTotemAll` — recebe ambos os tópicos

---

## Pré-requisitos

- Docker & Docker Compose instalados e funcionando.
- Go (Golang) instalado (versão compatível — ex.: `go version` deve responder).
- `git` (opcional).
- Portas 2181 (ZooKeeper) e 9092 (Kafka) livres no host.

---

## Estrutura do repositório

```
airport-monitor/
├── docker-compose.yaml        # Kafka + Zookeeper (Confluent)
└── airport/
    ├── go.mod
    ├── go.sum
    ├── publisher.go
    ├── subscriber.go
    ├── totem_all.go
    ├── totem_arrival.go
    ├── totem_departure.go
    ├── README.md               # este arquivo
    └── ... (docs, LICENSE, etc.)
```

---

## Passo a passo — subir o ambiente Kafka

1. Na raiz do projeto (onde está o `docker-compose.yaml`), rode:
```bash
docker-compose down -v   # (opcional) limpa contêineres antigos/volumes
docker-compose pull      # baixa imagens mais recentes (recomendado)
docker-compose up -d
```

2. Verifique que os containers estão rodando:
```bash
docker ps
# deve listar "zookeeper" e "kafka" com portas 2181 e 9092 mapeadas
```

3. (Opcional) Logs:
```bash
docker logs -f kafka
```

---

## Criar tópicos Kafka

Depois que os containers estiverem `Up` crie os tópicos usados pelo exemplo:

```bash
docker exec kafka kafka-topics --create   --topic arrival-flights   --bootstrap-server localhost:9092   --partitions 1 --replication-factor 1

docker exec kafka kafka-topics --create   --topic departure-flights   --bootstrap-server localhost:9092   --partitions 1 --replication-factor 1
```

Verifique:
```bash
docker exec kafka kafka-topics --list --bootstrap-server localhost:9092
# saída esperada: __consumer_offsets, arrival-flights, departure-flights
```

---

## Configurar o projeto Go (apenas se ainda não estiver)

No diretório `airport/`:

```bash
cd airport
go mod init airport              # cria go.mod (se ainda não existir)
go get github.com/segmentio/kafka-go
go mod tidy                       # organiza dependências
```

---

## Arquivos principais (resumo)

- `publisher.go`  
  Cria algumas mensagens de exemplo (voos) e publica em `arrival-flights` ou `departure-flights`. Um voo possui: `id, companhia, origem, destino, dataHora` (ou string equivalente).

- `subscriber.go`  
  Contém a função `consume(topic string, group string)` que cria um `kafka.Reader` e faz loop de leitura. Arquivo compartilhado pelos tótens.

- `totem_arrival.go` / `totem_departure.go` / `totem_all.go`  
  Chamam `consume` para o(s) tópico(s) apropriado(s). Ex.: `consume("arrival-flights", "arrival-group")`.

---

## Executar o exemplo (ordem recomendada)

Abra **4** terminais (ou abas):

1. Terminal A — Totem que lê ambos os tópicos:
```bash
cd airport
go run totem_all.go subscriber.go
```
<img width="1144" height="256" alt="image" src="https://github.com/user-attachments/assets/4f9a4994-1769-4f3f-b330-52a28e73670c" />


2. Terminal B — Totem de chegada:
```bash
cd airport
go run totem_arrival.go subscriber.go
```
<img width="1241" height="198" alt="image" src="https://github.com/user-attachments/assets/8521663e-2bbb-4d28-abfd-24dcd8908bd1" />


3. Terminal C — Totem de partida:
```bash
cd airport
go run totem_departure.go subscriber.go
```
<img width="1166" height="206" alt="image" src="https://github.com/user-attachments/assets/5cc16a23-8b53-4523-9ff1-ff47ff226e68" />


4. Terminal D — Publisher (publica alguns voos):
```bash
cd airport
go run publisher.go
```
<img width="961" height="304" alt="image" src="https://github.com/user-attachments/assets/0d58efd4-1f08-499b-b4a6-eac0124ebf8f" />


> Obs.: execute sempre os arquivos `totem_x.go` **junto** com `subscriber.go`, por exemplo:  
> `go run totem_arrival.go subscriber.go` — isso permite que a função `consume` (definida em `subscriber.go`) seja visível para o executável.

---

## Troubleshooting rápido

- `failed to dial: failed to open connection to localhost:9092`  
  → Verifique se o container Kafka está `Up` e escutando (use `docker ps` e `docker logs kafka`).

- `go: command not found`  
  → Instale o Go e adicione `C:\Program Files\Go\bin` ao `PATH`. Feche e reabra o terminal.

- `undefined: consume` ao executar apenas `go run totem_departure.go`  
  → Execute `go run totem_departure.go subscriber.go` (o Go precisa de ambos os arquivos).

- Se os tótens não recebem mensagens, verifique se os tópicos existem (`docker exec kafka kafka-topics --list ...`) e se o publisher está enviando para os tópicos corretos.

---

## Questões de estudo / exercícios propostos

1. Explique a diferença entre **broker**, **topic**, **partition** e **consumer group** no contexto do Kafka.  
2. Por que usamos **grupos de consumidores**? O que acontece se dois consumidores do mesmo grupo se inscrevem no mesmo tópico/partição?  
3. Modifique `publisher.go` para enviar mensagens no formato **JSON** (inclua campos adicionais, p.ex.: `status` , `gate`). Atualize os consumidores para decodificarem JSON.  
4. Implemente **persistência local** no `totem_all` que grava todas as mensagens recebidas em um arquivo CSV. Teste se o arquivo contém todas as entradas após várias execuções.  
5. Explore **partições**: altere o número de partições do tópico e observe como as mensagens são balanceadas entre consumidores (crie mais de um consumidor no mesmo grupo).  
6. Implemente um pequeno **web UI** (página estática + websocket) que mostra em tempo real as mensagens consumidas por `totem_all`.  
7. Substitua `localhost:9092` por `kafka:9092` nos readers/writers e empacote publisher + subscribers como serviços no `docker-compose` (cada um com seu `Dockerfile`) para rodar tudo via Docker.  
8. Analise garantias de entrega: qual a diferença entre `at-most-once`, `at-least-once` e `exactly-once` e como configurar escolhas no produtor/consumidor?

