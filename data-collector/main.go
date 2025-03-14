// /data-collector/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/cnpj"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/common"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/serper"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/tavily"
)

var (
	rabbitConn *amqp.Connection
	rabbitCh   *amqp.Channel
)

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

	// Se o BusinessName estiver vazio, não realiza o enriquecimento.
	if lead.BusinessName == "" {
		log.Printf("BusinessName vazio para o lead com ID: %s. Pulando enriquecimento.", lead.ID)
		return
	}

	// Monta a query para o Tavily usando BusinessName, City e State.
	query := lead.BusinessName
	if lead.City != "" {
		query += " " + lead.City
	}
	if lead.State != "" {
		query += " " + lead.State
	}
	// Se o lead for brasileiro, adiciona "CNPJ" à query.
	if lead.Country == "Brazil" || lead.Country == "Brasil" {
		query += " CNPJ"
	}

	log.Printf("Query enviada para o Tavily: %s", query)

	// Chamada à API Tavily para enriquecimento.
	enrichedData, err := tavily.EnrichLead(query)
	if err != nil {
		log.Printf("Erro ao enriquecer o lead com Tavily: %v", err)
		// Você pode optar por continuar mesmo sem os dados do Tavily.
	} else {
		log.Printf("Resposta bruta da API Tavily: %+v", enrichedData)
		// Extrai os dados do Tavily.
		cnpjTavily, phone, owner, email, website := tavily.ExtractLeadInfo(enrichedData)
		log.Printf("Dados do Tavily - CNPJ: %s, Phone: %s, Owner: %s, Email: %s, Website: %s",
			cnpjTavily, phone, owner, email, website)
	}

	serperResult, err := serper.FetchSerperDataForCNPJ(lead.BusinessName, lead.City)
	if err != nil {
		log.Printf("Erro na chamada à API Serper: %v", err)
	} else {
		log.Printf("Dados da API Serper: %+v", serperResult)

		if capturedIface, ok := serperResult["captured_cnpjs"]; ok {

			if cnpjs, ok := capturedIface.([]string); ok && len(cnpjs) > 0 {

				cnpjData, err := cnpj.FetchCNPJData(cnpjs[0])
				if err != nil {
					log.Printf("Erro ao consultar dados do CNPJ %s: %v", cnpjs[0], err)
				} else {
					log.Printf("Dados detalhados do CNPJ %s: %+v", cnpjs[0], cnpjData)
				}
			}
		}
	}

	// Aqui você pode combinar os dados de Tavily, Serper e Invertexto e enviá-los para o RabbitMQ,
	// se necessário, em um objeto estruturado para processamento posterior.
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
