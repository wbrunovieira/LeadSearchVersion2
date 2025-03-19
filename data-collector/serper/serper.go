// /data-collector/serper/serper.go
package serper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type OrganicResult struct {
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Link    string `json:"link"`
}

type SerperResponse struct {
	Organic []OrganicResult `json:"organic"`
}

// FetchSerperDataForCNPJ faz uma busca no Serper e retorna os CNPJs encontrados
func FetchSerperDataForCNPJ(name, city string, numResults int) (map[string]interface{}, error) {

	query := fmt.Sprintf("%s, %s CNPJ", name, city)
	payload := map[string]interface{}{
		"q":   query,
		"gl":  "br",
		"hl":  "pt-br",
		"num": numResults, // Agora podemos controlar a quantidade de resultados retornados
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

	// Se houver erro, capturar resposta para debugging
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta da API Serper (%d): %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta da API Serper: %v", err)
	}

	var serperResp SerperResponse
	if err := json.Unmarshal(bodyBytes, &serperResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar a resposta da API Serper: %v", err)
	}

	// Extrair CNPJs
	cnpjList := extractCNPJs(serperResp.Organic, name)

	return map[string]interface{}{
		"serper_info":    serperResp.Organic,
		"captured_cnpjs": cnpjList,
	}, nil
}

// extractCNPJs melhora a extração dos CNPJs e prioriza os mais relevantes
func extractCNPJs(results []OrganicResult, leadName string) []string {
	cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)
	digitRegex := regexp.MustCompile(`\b\d{14}\b`)

	cnpjScores := make(map[string]int) // Armazena pontuações dos CNPJs para priorizar o mais relevante

	for _, result := range results {
		foundCNPJs := append(cnpjRegex.FindAllString(result.Title, -1), cnpjRegex.FindAllString(result.Snippet, -1)...)
		foundCNPJs = append(foundCNPJs, digitRegex.FindAllString(result.Link, -1)...)

		for _, match := range foundCNPJs {
			normalizedCNPJ := NormalizeCNPJ(match)
			if normalizedCNPJ != "" {
				// Aumenta a pontuação se o nome do lead for mencionado no resultado
				score := 1
				if strings.Contains(strings.ToLower(result.Title), strings.ToLower(leadName)) ||
					strings.Contains(strings.ToLower(result.Snippet), strings.ToLower(leadName)) {
					score += 5
				}
				cnpjScores[normalizedCNPJ] += score
			}
		}
	}

	// Ordenar os CNPJs por relevância (maior pontuação primeiro)
	var sortedCNPJs []string
	for cnpj := range cnpjScores {
		sortedCNPJs = append(sortedCNPJs, cnpj)
	}
	sort.Slice(sortedCNPJs, func(i, j int) bool {
		return cnpjScores[sortedCNPJs[i]] > cnpjScores[sortedCNPJs[j]]
	})

	return sortedCNPJs
}

// NormalizeCNPJ formata um CNPJ para o padrão XX.XXX.XXX/XXXX-XX
func NormalizeCNPJ(cnpj string) string {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")
	if len(digits) == 14 {
		return fmt.Sprintf("%s.%s.%s/%s-%s", digits[:2], digits[2:5], digits[5:8], digits[8:12], digits[12:14])
	}
	return ""
}
