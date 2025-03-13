package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
)

var (
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

// initRabbitMQ conecta ao RabbitMQ, declara a fila e inicia o consumidor em uma goroutine.
func initRabbitMQ() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL não definida no ambiente")
	}

	var err error
	rabbitConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	rabbitCh, err = rabbitConn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal do RabbitMQ: %v", err)
	}

	// Declara a fila "lead_queue"
	_, err = rabbitCh.QueueDeclare(
		"lead_queue", // nome da fila
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila: %v", err)
	}

	// Inicia o consumidor em uma goroutine
	go consumeRabbitMQ()
}

// consumeRabbitMQ consome as mensagens da fila e as imprime no log.
func consumeRabbitMQ() {
	msgs, err := rabbitCh.Consume(
		"lead_queue", // nome da fila
		"",           // consumer tag
		true,         // auto-acknowledge
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao registrar o consumidor: %v", err)
	}

	for msg := range msgs {
		log.Printf("Mensagem recebida do RabbitMQ: %s", string(msg.Body))
	}
}

func main() {
	// Inicializa o RabbitMQ e o consumidor de mensagens
	initRabbitMQ()
	defer rabbitCh.Close()
	defer rabbitConn.Close()

	// Configura o endpoint HTTP
	http.HandleFunc("/", handler)
	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
