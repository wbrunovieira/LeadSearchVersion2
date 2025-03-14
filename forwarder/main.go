// /forwarder/main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// CombinedLeadData representa os dados combinados recebidos e o prompt adicional.
type CombinedLeadData struct {
	Lead        interface{}            `json:"lead"`
	TavilyData  interface{}            `json:"tavily_data,omitempty"`
	TavilyExtra interface{}            `json:"tavily_extra,omitempty"`
	SerperData  map[string]interface{} `json:"serper_data,omitempty"`
	CNPJData    map[string]interface{} `json:"cnpj_data,omitempty"`
	Prompt      string                 `json:"prompt,omitempty"`
}

// OlhamaPayload define a estrutura que o Olhama espera.
type OlhamaPayload struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

var (
	amqpConn *amqp.Connection
	amqpCh   *amqp.Channel
)

// publishToOlhama envia os dados combinados (incluindo o prompt) para o endpoint do Olhama e retorna o corpo da resposta.
func publishToOlhama(data CombinedLeadData) ([]byte, error) {
	olhamaURL := os.Getenv("OLHAMA_URL")
	if olhamaURL == "" {
		// Endereço padrão; certifique-se de que o Olhama está acessível neste endereço.
		olhamaURL = "http://192.168.0.7:11434/api/chat"
	}

	// Converte os dados combinados para uma string JSON com identação.
	combinedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados combinados para JSON: %v", err)
	}

	// Monta o payload esperado pelo Olhama.
	var payload OlhamaPayload
	payload.Model = "deepseek-r1"
	payload.Stream = false
	messageContent := fmt.Sprintf("%s\n\nCombined Data:\n%s", data.Prompt, string(combinedJSON))
	payload.Messages = []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		{
			Role:    "user",
			Content: messageContent,
		},
	}

	// Converte o payload para JSON.
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	// Cria a requisição HTTP POST.
	req, err := http.NewRequest("POST", olhamaURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Usa um cliente HTTP com timeout.
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar dados para Olhama: %v", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta.
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta do Olhama: %v", err)
	}

	// Se o status HTTP for de erro, retorna erro.
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Olhama retornou status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Println("Dados enviados com sucesso para Olhama")
	return respBody, nil
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

	// Define prefetch count para 1 para processar uma mensagem de cada vez.
	if err := amqpCh.Qos(1, 0, false); err != nil {
		log.Fatalf("Erro ao definir QoS: %v", err)
	}

	// Declara a fila "combined_leads_queue".
	_, err = amqpCh.QueueDeclare(
		"combined_leads_queue", // nome da fila
		true,                   // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		log.Fatalf("Erro ao declarar a fila 'combined_leads_queue': %v", err)
	}
}

func consumeCombinedData() {
	msgs, err := amqpCh.Consume(
		"combined_leads_queue", // nome da fila
		"",
		false, // autoAck desativado: queremos confirmar só se o envio foi bem-sucedido.
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao consumir a fila: %v", err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Mensagem combinada recebida: %s", d.Body)

			var data CombinedLeadData
			if err := json.Unmarshal(d.Body, &data); err != nil {
				log.Printf("Erro ao decodificar mensagem combinada: %v", err)
				// Se o erro for irrecuperável, não requeue.
				d.Nack(false, false)
				continue
			}

			// Adiciona um prompt para a análise do Olhama.
			data.Prompt = "Please analyze the lead data below. Evaluate its completeness and quality, identify any missing critical details, and suggest additional data to enhance the lead profile."

			// Envia os dados para o Olhama e captura o retorno.
			respBody, err := publishToOlhama(data)
			if err != nil {
				log.Printf("Erro ao enviar dados para Olhama: %v", err)
				// Se houver erro, requeue a mensagem para tentar novamente.
				d.Nack(false, true)
				continue
			}

			// Loga o retorno do Olhama (você pode processá-lo ou salvar em um banco de dados, se desejar).
			log.Printf("Resposta do Olhama: %s", string(respBody))

			// Apenas após receber um retorno bem-sucedido, confirma o processamento da mensagem.
			d.Ack(false)
		}
	}()
	log.Println("Aguardando mensagens na fila 'combined_leads_queue' para enviar ao Olhama...")
	<-forever
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado")
	}

	initAMQP()
	defer amqpConn.Close()
	defer amqpCh.Close()

	consumeCombinedData()

	// Inicia um servidor HTTP simples para manter o container ativo (opcional).
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Forwarder service is running"))
	})
	log.Printf("Forwarder service rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
