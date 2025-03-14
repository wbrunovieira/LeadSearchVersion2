// /data-collector/serper/serper.go
package serper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// Estruturas que representam a resposta da API Serper
type OrganicResult struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Link    string `json:"link"`
}

type SerperResponse struct {
	Organic []OrganicResult `json:"organic"`
}

// FetchSerperDataForCNPJ faz uma requisição à API Serper com o termo "CNPJ"
// e extrai os CNPJs dos resultados orgânicos.
func FetchSerperDataForCNPJ(name, city string) (map[string]interface{}, error) {
	// Monta a query, adicionando o termo "CNPJ"
	query := fmt.Sprintf("%s, %s CNPJ", name, city)
	payload := map[string]interface{}{
		"q":   query,
		"gl":  "br",
		"hl":  "pt-br",
		"num": 30,
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
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta da API Serper: %s", string(body))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta: %v", err)
	}

	var serperResp SerperResponse
	if err := json.Unmarshal(bodyBytes, &serperResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar a resposta da API Serper: %v", err)
	}

	cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)

	digitRegex := regexp.MustCompile(`\b\d{14}\b`)
	capturedCNPJs := make(map[string]bool)

	for _, result := range serperResp.Organic {
		for _, match := range cnpjRegex.FindAllString(result.Title, -1) {
			if norm := NormalizeCNPJ(match); norm != "" {
				capturedCNPJs[norm] = true
			}
		}
		for _, match := range cnpjRegex.FindAllString(result.Snippet, -1) {
			if norm := NormalizeCNPJ(match); norm != "" {
				capturedCNPJs[norm] = true
			}
		}
		for _, match := range digitRegex.FindAllString(result.Link, -1) {
			if norm := NormalizeCNPJ(match); norm != "" {
				capturedCNPJs[norm] = true
			}
		}
	}

	capturedList := []string{}
	for cnpj := range capturedCNPJs {
		capturedList = append(capturedList, cnpj)
	}

	return map[string]interface{}{
		"serper_info":    serperResp.Organic,
		"captured_cnpjs": capturedList,
	}, nil
}

func NormalizeCNPJ(cnpj string) string {

	digits := regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")
	if len(digits) == 14 {
		return fmt.Sprintf("%s.%s.%s/%s-%s", digits[:2], digits[2:5], digits[5:8], digits[8:12], digits[12:14])
	}
	return ""
}
