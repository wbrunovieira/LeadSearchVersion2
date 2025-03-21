package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/cnpjsearch"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/common"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/serper"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/tavily"
)

var (
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
)

type CombinedLeadData struct {
	Lead       common.Lead            `json:"lead"`
	TavilyData *tavily.TavilyResponse `json:"tavily_data,omitempty"`

	TavilyExtra struct {
		CNPJ    string `json:"cnpj,omitempty"`
		Phone   string `json:"phone,omitempty"`
		Owner   string `json:"owner,omitempty"`
		Email   string `json:"email,omitempty"`
		Website string `json:"website,omitempty"`
	} `json:"tavily_extra,omitempty"`
	SerperData map[string]interface{} `json:"serper_data,omitempty"`

	CNPJData map[string]interface{} `json:"cnpj_data,omitempty"`
}

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

	go consumeRabbitMQ()
}

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
		processLeadMessage(msg.Body)
	}
}

func processLeadMessage(body []byte) {
	var lead common.Lead
	if err := json.Unmarshal(body, &lead); err != nil {
		log.Printf("Erro ao decodificar mensagem do RabbitMQ: %v", err)
		return
	}

	log.Printf("Lead recebido: %+v", lead)
	log.Printf("Data Collector - Lead recebido com ID: %s", lead.ID)

	if lead.BusinessName == "" {
		log.Printf("BusinessName vazio para o lead com ID: %s. Pulando enriquecimento.", lead.ID)
		return
	}

	query := lead.BusinessName
	if lead.City != "" {
		query += " " + lead.City
	}
	if lead.State != "" {
		query += " " + lead.State
	}

	log.Printf("Query enviada para o Tavily: %s", query)

	var combinedData CombinedLeadData
	combinedData.Lead = lead
	log.Printf("combinedData.Lead: %+v", combinedData.Lead)

	maxResults := 5

	tavilyData, err := tavily.FetchTavilyData(query, maxResults)
	if err != nil {
		log.Printf("Erro ao enriquecer o lead com Tavily: %v", err)
	} else {
		log.Printf("Resposta bruta da API Tavily: %+v", tavilyData)
		combinedData.TavilyData = tavilyData

		log.Printf("Dados do Tavily %v",
			tavilyData)
	}

	serperResult, err := serper.FetchSerperDataForCNPJ(lead.BusinessName, lead.City, 10)
	if err != nil {
		log.Printf("Erro na chamada à API Serper: %v", err)
	} else {
		log.Printf("Dados da API Serper: %+v", serperResult)
		combinedData.SerperData = serperResult
	}

	cnpjData, err := cnpjsearch.FetchDataCNPJBIZ(lead.BusinessName, lead.City)
	if err != nil {
		log.Printf("Erro ao buscar dados CNPJ: %v", err)
	} else {
		log.Printf("CNPJData BIZ: %+v", combinedData.CNPJData)
		combinedData.CNPJData = cnpjData
	}

	log.Printf("combinedData final - Lead: %+v", combinedData.Lead)
	log.Printf("combinedData final - Lead.ID: %s", combinedData.Lead.ID)
	log.Printf("combinedData final: %v", combinedData)

	if err := PublishCombinedLead(combinedData); err != nil {
		log.Printf("Erro ao publicar dados combinados no RabbitMQ: %v", err)
	} else {
		log.Printf("Dados combinados publicados com sucesso para o lead ID: %s", lead.ID)
	}
}

func PublishCombinedLead(data CombinedLeadData) error {
	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("erro ao converter dados combinados para JSON: %v", err)
	}

	err = rabbitCh.Publish(
		"",                     // Usando a default exchange
		"combined_leads_queue", // nome da fila
		false,                  // mandatory
		false,                  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("erro ao publicar dados combinados: %v", err)
	}
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado")
	}

	initRabbitMQ()
	defer rabbitCh.Close()
	defer rabbitConn.Close()

	log.Println("Servidor iniciado na porta 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
