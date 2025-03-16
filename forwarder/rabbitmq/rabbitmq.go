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

			leadID := ""
			if leadMap, ok := data.Lead.(map[string]interface{}); ok {
				if id, exists := leadMap["id"]; exists {
					if idStr, ok := id.(string); ok {
						leadID = idStr
					} else {
						log.Printf("ID do lead não é uma string")
					}
				} else {
					log.Printf("ID do lead não encontrado")
				}
			} else {
				log.Printf("Erro ao converter 'lead' para map[string]interface{}")
			}

			log.Printf("ID do Lead: %s", leadID)

			data.Prompt = "Por favor, analise os dados do lead a seguir. Esses dados foram obtidos através de uma busca no Google, representando o cliente identificado. Adicionalmente, nas chaves 'TavilyData', 'SerperData' e 'CNPJData' foram realizadas buscas complementares sobre este lead. Com base nessas informações, extraia e identifique os seguintes dados:\n\n- Razão Social da empresa (RegisteredName)\n- CNPJ\n- Contatos e/ou dados do responsável pela empresa, incluindo:\n    - Nome\n    - Telefone\n    - Email\n- Data de Fundação\n- Website\n- Redes Sociais (Facebook, Instagram, TikTok, WhatsApp)\n\nInstruções adicionais:\n1. Retorne a resposta estritamente no seguinte formato JSON, sem nenhum comentário ou texto adicional, exceto um bloco interno de raciocínio (entre <think> e </think>) contendo uma breve explicação do processo de extração dos dados.\n2. Inclua no JSON uma chave \"AnaliseEmpresa\" contendo uma análise resumida das informações da empresa, destacando pontos fortes e oportunidades para abordagem comercial.\n3. No final, inclua uma chave \"MensagemWhatsApp\" com uma mensagem personalizada para ser enviada via WhatsApp. A mensagem deve:\n   - Variar para cada empresa, evitando repetições excessivas (pois mensagens muito padronizadas podem levar a bloqueio no WhatsApp);\n   - Substituir a placeholder [Nome] pelo nome do contato, se disponível; caso contrário, utilizar \"Prezado Cliente\";\n   - Ser clara e persuasiva, destacando os benefícios de uma presença digital profissional (como um site moderno e integração com redes sociais) e convidando o lead para uma conversa.\n4. Ao identificar as informações, compare-as com os dados enviados originalmente e, se houver divergências, adicione as informações complementares sem substituir os dados originais.\n\nExemplo de formato JSON esperado:\n\n```json\n{\n  \"<think>\": \"Explicação interna do raciocínio usado para extrair os dados.\",\n  \"RegisteredName\": \"<valor>\",\n  \"CNPJ\": \"<valor>\",\n  \"Contatos\": {\n    \"Nome\": \"<valor>\",\n    \"Telefone\": \"<valor>\",\n    \"Email\": \"<valor>\"\n  },\n  \"DataFundacao\": \"<valor>\",\n  \"Website\": \"<valor>\",\n  \"RedesSociais\": {\n    \"Facebook\": \"<valor>\",\n    \"Instagram\": \"<valor>\",\n    \"TikTok\": \"<valor>\",\n    \"WhatsApp\": \"<valor>\"\n  },\n  \"AnaliseEmpresa\": \"Breve análise das informações da empresa, destacando pontos fortes e oportunidades para abordagem comercial.\",\n  \"MensagemWhatsApp\": \"Olá, [Nome ou Prezado Cliente]! Aqui é o Bruno da WB Digital Solutions. Observamos que sua empresa possui [inserir diferencial ou oportunidade identificada]. Um website moderno pode ampliar significativamente sua presença digital e atrair mais clientes. Que tal conversarmos para explorar essa oportunidade? Estou à disposição para agendarmos uma reunião.\" \n}\n```\n\nResponda em português."
			log.Printf("Enviado para Olhama: %s", data)

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
			if err := helpers.UpdateLeadField(leadID, "RegisteredName", olhamaResp.RegisteredName); err != nil {
				log.Printf("Erro ao atualizar RegisteredName: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("CNPJ: %s", olhamaResp.CNPJ)
			if err := helpers.UpdateLeadField(leadID, "CompanyRegistrationID", olhamaResp.CNPJ); err != nil {
				log.Printf("Erro ao atualizar CompanyRegistrationID: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Contatos: %+v", olhamaResp.Contatos)
			if err := helpers.UpdateLeadField(leadID, "Owner", olhamaResp.Contatos); err != nil {
				log.Printf("Erro ao atualizar Owner: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("DataFundacao: %s", olhamaResp.DataFundacao)
			if err := helpers.UpdateLeadField(leadID, "FoundationDate", olhamaResp.DataFundacao); err != nil {
				log.Printf("Erro ao atualizar FoundationDate: %v", err)
				d.Nack(false, true)
				continue
			}

			log.Printf("Website: %s", olhamaResp.Website)
			if err := helpers.UpdateLeadField(leadID, "Website", olhamaResp.Website); err != nil {
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
