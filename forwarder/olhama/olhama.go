// forwarder/olhama/olhama.go
package olhama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/types"
)

func Publish(data types.CombinedLeadData) ([]byte, error) {
	olhamaURL := os.Getenv("OLHAMA_URL")
	if olhamaURL == "" {

		olhamaURL = "http://localhost:11434/api/chat"
	}

	combinedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados combinados para JSON: %v", err)
	}

	payload := types.OlhamaPayload{
		Model:       "qwen2.5:14b",
		Stream:      false,
		Temperature: 0.4,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "system",
				Content: "Você é um especialista em geração de listas de leads. Seu papel é identificar e extrair informações relevantes sobre a empresa, incluindo CNPJ, socios, telefone e dados de contato.",
			},
			{
				Role:    "user",
				Content: fmt.Sprintf("%s\n\nDados combinados:\n%s", data.Prompt, string(combinedJSON)),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	req, err := http.NewRequest("POST", olhamaURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição HTTP: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar dados para Olhama: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta do Olhama: %v", err)
	}

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("olhama retornou status %d: %s", resp.StatusCode, string(respBody))
	}

	log.Println("Dados enviados com sucesso para Olhama")
	return respBody, nil
}

func CallOlhama(data types.CombinedLeadData) (types.OlhamaResponse, error) {
	var respData types.OlhamaResponse
	var outerResp types.OlhamaOuterResponse

	rawResp, err := Publish(data)
	if err != nil {
		return respData, fmt.Errorf("erro ao enviar dados para Olhama: %v", err)
	}

	if err := json.Unmarshal(rawResp, &outerResp); err != nil {
		return respData, fmt.Errorf("erro ao decodificar resposta externa do Olhama: %v", err)
	}
	log.Printf("CallOlhama - Outer Response: %+v", outerResp)

	cleanedStr := cleanResponse(outerResp.Message.Content)
	cleanedBytes := []byte(cleanedStr)

	if err := json.Unmarshal(cleanedBytes, &respData); err != nil {
		return respData, fmt.Errorf("erro ao decodificar resposta interna do Olhama: %v", err)
	}
	log.Printf("CallOlhama - Inner Response (respData): %+v", respData)

	return respData, nil
}
func CallOlhama2(data types.CombinedLeadData) (types.OlhamaResponse, error) {
	var respData types.OlhamaResponse
	var outerResp types.OlhamaOuterResponse

	rawResp, err := Publish(data)
	if err != nil {
		return respData, fmt.Errorf("erro ao enviar dados para Olhama 2: %v", err)
	}

	if err := json.Unmarshal(rawResp, &outerResp); err != nil {
		return respData, fmt.Errorf("erro ao decodificar resposta externa do Olhama2: %v", err)
	}
	log.Printf("CallOlhama2 - Outer Response: %+v", outerResp)

	cleanedStr := cleanResponse(outerResp.Message.Content)
	cleanedBytes := []byte(cleanedStr)

	if err := json.Unmarshal(cleanedBytes, &respData); err != nil {
		return respData, fmt.Errorf("erro ao decodificar resposta interna do Olhama2: %v", err)
	}
	log.Printf("CallOlhama 2- Inner Response (respData): %+v", respData)

	return respData, nil
}

func cleanResponse(response string) string {
	log.Printf("cleanResponse - Raw response received: %s", response)
	cleaned := response
	if strings.HasPrefix(cleaned, "```json") {
		log.Printf("cleanResponse - Detected Markdown delimiters. Removing them...")
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
		log.Printf("cleanResponse - After removing markdown: %s", cleaned)
	} else {
		log.Printf("cleanResponse - No Markdown delimiters detected.")
	}

	if idx := strings.Index(cleaned, "<think>"); idx != -1 {
		log.Printf("cleanResponse - Detected '<think>' block. Removing content starting from index %d", idx)
		cleaned = strings.TrimSpace(cleaned[:idx])
		log.Printf("cleanResponse - After removing <think> block: %s", cleaned)
	}

	cleaned = strings.Trim(cleaned, "`")
	log.Printf("cleanResponse - After trimming extra backticks: %s", cleaned)

	return cleaned
}
