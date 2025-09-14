package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/streadway/amqp"
)

type CombinedLeadData struct {
	Lead        interface{}            `json:"lead"`
	TavilyData  interface{}            `json:"tavily_data,omitempty"`
	TavilyExtra interface{}            `json:"tavily_extra,omitempty"`
	SerperData  map[string]interface{} `json:"serper_data,omitempty"`
	CNPJData    map[string]interface{} `json:"cnpj_data,omitempty"`
}

var (
	esClient *elasticsearch.Client
	amqpConn *amqp.Connection
	amqpCh   *amqp.Channel
)

func initElasticsearch() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	var err error
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})
	if err != nil {
		log.Fatalf("Erro ao criar o cliente Elasticsearch: %v", err)
	}

	res, err := esClient.Info()
	if err != nil {
		log.Fatalf("Erro ao obter informações do Elasticsearch: %v", err)
	}
	defer res.Body.Close()
	log.Println("Cliente Elasticsearch criado com sucesso.")
}

func initAMQP() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL não definida no ambiente")
	}

	var err error
	amqpConn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	amqpCh, err = amqpConn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal do RabbitMQ: %v", err)
	}

	// Declara o exchange fanout
	err = amqpCh.ExchangeDeclare(
		"leads.fanout", // nome do exchange
		"fanout",       // tipo
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar exchange fanout: %v", err)
	}

	// Cria fila exclusiva para o datalake
	_, err = amqpCh.QueueDeclare(
		"datalake_queue", // nome único da fila
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila 'datalake_queue': %v", err)
	}

	// Faz bind da fila ao exchange fanout
	err = amqpCh.QueueBind(
		"datalake_queue", // nome da fila
		"",               // routing key (não usado em fanout)
		"leads.fanout",   // exchange
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao fazer bind da fila ao exchange: %v", err)
	}
}

func consumeCombinedData() {
	msgs, err := amqpCh.Consume(
		"datalake_queue", // consome da fila específica do datalake
		"",               // consumer tag
		true,             // auto-acknowledge
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao consumir fila: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Mensagem recebida: %s", d.Body)
			indexCombinedData(d.Body)
		}
	}()

	log.Println("Aguardando mensagens na fila combined_leads_queue...")
	<-forever
}

func indexCombinedData(body []byte) {

	var doc CombinedLeadData
	if err := json.Unmarshal(body, &doc); err != nil {
		log.Printf("Erro ao decodificar documento: %v", err)

	}

	res, err := esClient.Index(
		"combined_leads",
		bytes.NewReader(body),
		esClient.Index.WithContext(context.Background()),
		esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Printf("Erro ao indexar documento: %v", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("Erro ao indexar documento: %s", res.String())
		return
	}
	log.Println("Documento indexado com sucesso no índice 'combined_leads'.")
}

func main() {

	initElasticsearch()
	initAMQP()
	defer amqpConn.Close()
	defer amqpCh.Close()

	consumeCombinedData()
}
