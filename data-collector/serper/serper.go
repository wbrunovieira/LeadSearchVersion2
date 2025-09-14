package serper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchSerperDataForCNPJ(name, city string, numResults int) (map[string]interface{}, error) {
	query := fmt.Sprintf("%s, %s CNPJ", name, city)
	payload := map[string]interface{}{
		"q":   query,
		"gl":  "br",
		"hl":  "pt-br",
		"num": numResults,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	apiURL := "https://google.serper.dev/search"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}
	req.Header.Set("X-API-KEY", os.Getenv("SERPER_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição à API Serper: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta da API Serper (%d): %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta da API Serper: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("erro ao decodificar a resposta da API Serper: %v", err)
	}
	return result, nil
}
