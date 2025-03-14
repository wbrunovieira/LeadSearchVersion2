// /api/rabbitmq/producer.go
package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	ch   *amqp.Channel
)

func InitRabbitMQ(url string) error {
	var err error

	conn, err = amqp.Dial(url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"lead_exchange", // nome da exchange
		"fanout",        // tipo (fanout envia para todas as filas vinculadas)
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(
		"lead_queue", // nome da fila
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		"lead_queue",    // nome da fila
		"",              // routing key
		"lead_exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	log.Println("RabbitMQ inicializado com sucesso")
	return nil
}

func PublishLeadID(leadID string) error {
	err := ch.Publish(
		"lead_exchange", // exchange
		"",              // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(leadID),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	log.Printf("Publicado lead id %s na exchange", leadID)
	return nil
}

func PublishLead(lead interface{}) error {

	leadJSON, err := json.Marshal(lead)
	if err != nil {
		return fmt.Errorf("failed to marshal lead: %v", err)
	}

	err = ch.Publish(
		"lead_exchange", // exchange
		"",              // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        leadJSON,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish lead: %v", err)
	}
	log.Printf("Publicado lead com dados %s na exchange", leadJSON)
	return nil
}

func CloseRabbitMQ() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
