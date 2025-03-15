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
	"time"

	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/types"
)

func Publish(data types.CombinedLeadData) ([]byte, error) {
	olhamaURL := os.Getenv("OLHAMA_URL")
	if olhamaURL == "" {

		olhamaURL = "http://192.168.0.7:11434/api/chat"
	}

	combinedJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("erro ao converter dados combinados para JSON: %v", err)
	}

	payload := types.OlhamaPayload{
		Model:  "deepseek-r1",
		Stream: false,
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: fmt.Sprintf("%s\n\nCombined Data:\n%s", data.Prompt, string(combinedJSON)),
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
