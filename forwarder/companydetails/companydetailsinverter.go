// /forwarder/companydetails/companydetailsinverter.go
package companydetails

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func FetchCNPJDataInverter(cnpj string) (map[string]interface{}, error) {

	cleanCNPJ := regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")
	apiToken := os.Getenv("INVERTEXTO_API_TOKEN")
	if apiToken == "" {
		return nil, fmt.Errorf("INVERTEXTO_API_TOKEN não definida")
	}
	apiURL := fmt.Sprintf("https://api.invertexto.com/v1/cnpj/%s?token=%s", cleanCNPJ, apiToken)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para Invertexto: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro na resposta da API Invertexto: %s", string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta da Invertexto: %v", err)
	}

	var cnpjData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &cnpjData); err != nil {
		return nil, fmt.Errorf("erro ao decodificar a resposta da Invertexto: %v", err)
	}

	return cnpjData, nil
}
