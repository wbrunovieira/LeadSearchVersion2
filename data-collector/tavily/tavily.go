package tavily

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
)

type TavilyResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

type TavilyResponse struct {
	Query        string         `json:"query"`
	Results      []TavilyResult `json:"results"`
	ResponseTime float64        `json:"response_time"`
}

func EnrichLead(query string, maxResults int) (*TavilyResponse, error) {
	apiKey := os.Getenv("TAVILY_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("TAVILY_API_KEY não definida no ambiente")
	}

	payload := map[string]interface{}{
		"query":               query,
		"max_results":         maxResults,
		"search_depth":        "advanced",
		"include_answer":      false,
		"include_raw_content": true,
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta: %v", err)
	}
	log.Printf("Resposta bruta da API Tavily: %s", string(bodyBytes))

	var tavilyResp TavilyResponse
	if err := json.Unmarshal(bodyBytes, &tavilyResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta da API Tavily: %v", err)
	}

	sort.SliceStable(tavilyResp.Results, func(i, j int) bool {
		return tavilyResp.Results[i].Score > tavilyResp.Results[j].Score
	})

	return &tavilyResp, nil
}

func ExtractLeadInfo(resp *TavilyResponse, leadName string) (string, string, string, string, string) {
	cnpjRegex := regexp.MustCompile(`\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}`)
	phoneRegex := regexp.MustCompile(`\+\d{2}\s?\d{2,3}\s?\d{4,5}-\d{4}`)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	var cnpjList, phoneList, emailList []string
	var website, owner string

	for _, result := range resp.Results {
		if len(cnpjList) == 0 {
			cnpjList = cnpjRegex.FindAllString(result.Content, -1)
		}
		if len(phoneList) == 0 {
			phoneList = phoneRegex.FindAllString(result.Content, -1)
		}
		if len(emailList) == 0 {
			emailList = emailRegex.FindAllString(result.Content, -1)
		}
		if website == "" && result.URL != "" {
			website = result.URL
		}
		if owner == "" && result.Title != "" {
			owner = result.Title
		}
	}

	return firstOrEmpty(cnpjList), firstOrEmpty(phoneList), owner, firstOrEmpty(emailList), website
}

func firstOrEmpty(list []string) string {
	if len(list) > 0 {
		return list[0]
	}
	return ""
}
