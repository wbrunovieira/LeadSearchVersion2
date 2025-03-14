// /data-collector/tavily/tavily.go
package tavily

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"net/http"
	"os"
)

type TavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type TavilyResponse struct {
	Query             string         `json:"query"`
	FollowUpQuestions interface{}    `json:"follow_up_questions"`
	Answer            interface{}    `json:"answer"`
	Images            []interface{}  `json:"images"`
	Results           []TavilyResult `json:"results"`
	ResponseTime      float64        `json:"response_time"`
}

func EnrichLead(query string) (*TavilyResponse, error) {
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY não definida no ambiente")
	}

	payload := map[string]string{
		"query": query,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("erro ao converter payload para JSON: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/search", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição à Tavily: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta: %v", err)
	}
	log.Printf("Resposta bruta da API Tavily: %s", string(bodyBytes))

	var tavilyResp TavilyResponse
	if err := json.Unmarshal(bodyBytes, &tavilyResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta da API Tavily: %v", err)
	}

	return &tavilyResp, nil
}

func ExtractLeadInfo(resp *TavilyResponse) (cnpj, phone, owner, email, website string) {

	cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}\/\d{4}\-\d{2}`)
	phoneRegex := regexp.MustCompile(`\+\d{2}\s?\d{2,3}\s?\d{4,5}\-\d{4}`)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

	for _, result := range resp.Results {

		if cnpj == "" {
			if matches := cnpjRegex.FindStringSubmatch(result.Content); len(matches) > 0 {
				cnpj = matches[0]
			}
		}

		if phone == "" {
			if matches := phoneRegex.FindStringSubmatch(result.Content); len(matches) > 0 {
				phone = matches[0]
			}
		}

		if email == "" {
			if matches := emailRegex.FindStringSubmatch(result.Content); len(matches) > 0 {
				email = matches[0]
			}
		}

		if website == "" && result.URL != "" {

			if len(result.URL) >= 4 && (result.URL[:4] == "http" || result.URL[:5] == "https") {
				website = result.URL
			}
		}

		if owner == "" && result.Title != "" {

		}
	}

	return cnpj, phone, owner, email, website
}
