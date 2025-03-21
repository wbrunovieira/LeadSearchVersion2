// rabbitmq/rabbitmq.go
package rabbitmq

import (
	"encoding/json"
	"log"
	"os"

	"github.com/streadway/amqp"
	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/helpers"
	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/olhama"
	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/types"
)

var (
	Conn *amqp.Connection
	Ch   *amqp.Channel
)

func Connect() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL não definida no ambiente")
	}

	var err error
	Conn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao RabbitMQ: %v", err)
	}

	Ch, err = Conn.Channel()
	if err != nil {
		log.Fatalf("Erro ao abrir canal do RabbitMQ: %v", err)
	}

	if err := Ch.Qos(1, 0, false); err != nil {
		log.Fatalf("Erro ao definir QoS: %v", err)
	}

	_, err = Ch.QueueDeclare(
		"combined_leads_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila 'combined_leads_queue': %v", err)
	}
}

func ConsumeQueue() {
	msgs, err := Ch.Consume(
		"combined_leads_queue",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao consumir a fila: %v", err)
	}

	go func() {
		for d := range msgs {
			log.Printf("Mensagem combinada recebida: %s", d.Body)

			var data types.CombinedLeadData
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Erro ao decodificar mensagem combinada: %v", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("ID do Lead: %s", data.Lead.ID)
			data.Prompt = "Por favor, analise os dados do lead a seguir, obtidos por buscas no Google e complementados pelas chaves \"TavilyData\", \"SerperData\" e \"CNPJData\". Extraia e identifique estritamente os seguintes dados:\n- Razão Social da empresa (RegisteredName)\n- CNPJ\n- Contatos: insira o(s) nome(s) do(s) responsável(is) pela empresa, concatenando-os em uma única string separada por vírgula\n- Data de Fundação\n- Website\n- Redes Sociais: retorne um objeto com as chaves \"Facebook\", \"Instagram\", \"TikTok\" e \"WhatsApp\"\n\nInstruções adicionais:\n1. Retorne a resposta estritamente no formato JSON, sem nenhum texto adicional fora do JSON.\n2. Inclua um bloco interno de raciocínio entre <think> e </think> que contenha uma breve explicação do processo de extração dos dados.\n3. Não crie um campo separado para \"Responsavel\"; utilize apenas o campo \"Contatos\".\n4. Se algum campo não estiver disponível, retorne null (para dados ausentes) ou uma string vazia, conforme apropriado.\n5. A estrutura do JSON deve ser exatamente conforme especificado, sem campos extras.\n6. Caso seja identificado mais de um CNPJ nos dados, verifique o endereço do lead para confirmar qual CNPJ corresponde à unidade correta.\n\nExemplo de saída:\n{\n  \"RegisteredName\": \"Nome da Razão Social\",\n  \"CNPJ\": \"XX.XXX.XXX/0001-XX\",\n  \"Contatos\": \"Nome1, Nome2\",\n  \"DataFundacao\": \"AAAA-MM-DD\",\n  \"Website\": \"https://exemplo.com\",\n  \"RedesSociais\": {\n    \"Facebook\": \"url ou null\",\n    \"Instagram\": \"url ou null\",\n    \"TikTok\": \"url ou null\",\n    \"WhatsApp\": \"url ou null\"\n  },\n  \"<think>\": \"Breve explicação de como os dados foram extraídos...\"\n}"
			log.Printf("Enviado para Olhama: %+v", data)

			response, err := olhama.Publish(data)
			if err != nil {
				log.Printf("Erro ao enviar dados para Olhama: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Resposta do Olhama: %s", response)

			olhamaResp, err := olhama.CallOlhama(data)
			if err != nil {
				log.Printf("Erro ao processar resposta do Olhama: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Bloco de Raciocínio: %s", olhamaResp.Think)
			log.Printf("RegisteredName: %s", olhamaResp.RegisteredName)
			if err := helpers.UpdateLeadField(data.Lead.ID.String(), "RegisteredName", olhamaResp.RegisteredName); err != nil {
				log.Printf("Erro ao atualizar RegisteredName: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("CNPJ: %s", olhamaResp.CNPJ)
			if err := helpers.UpdateLeadField(data.Lead.ID.String(), "CompanyRegistrationID", olhamaResp.CNPJ); err != nil {
				log.Printf("Erro ao atualizar CompanyRegistrationID: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Contatos: %+v", olhamaResp.Contatos)
			if err := helpers.UpdateLeadField(data.Lead.ID.String(), "Owner", olhamaResp.Contatos); err != nil {
				log.Printf("Erro ao atualizar Owner: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("DataDeFundacao: %s", olhamaResp.DataDeFundacao)
			if err := helpers.UpdateLeadField(data.Lead.ID.String(), "FoundationDate", olhamaResp.DataDeFundacao); err != nil {
				log.Printf("Erro ao atualizar FoundationDate: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Website: %s", olhamaResp.Website)
			if err := helpers.UpdateLeadField(data.Lead.ID.String(), "Website", olhamaResp.Website); err != nil {
				log.Printf("Erro ao atualizar Website: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("RedesSociais: %+v", olhamaResp.RedesSociais)
			log.Printf("AnaliseEmpresa: %s", olhamaResp.AnaliseEmpresa)
			log.Printf("MensagemWhatsApp: %s", olhamaResp.Message.Content)

			d.Ack(false)
		}
	}()
	log.Println("Aguardando mensagens na fila 'combined_leads_queue' para enviar ao Olhama...")
	select {}
}
