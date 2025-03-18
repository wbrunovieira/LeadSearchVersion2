package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func UpdateLeadField(leadID string, field string, value interface{}) error {

	if field == "FoundationDate" {
		if dateStr, ok := value.(string); ok && dateStr == "" {
			log.Printf("Valor vazio para FoundationDate, atualização ignorada")
			return nil
		}
	}

	updateURL := "http://api:8085/update-lead-field"
	payload := map[string]interface{}{
		"id":    leadID,
		"field": field,
		"value": value,
	}
	log.Printf("entrou no updateLeadField com payload: %v", payload)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	req, err := http.NewRequest("PUT", updateURL, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar a requisição: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	log.Printf("updateLeadField req: %v", req)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro na atualização (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
