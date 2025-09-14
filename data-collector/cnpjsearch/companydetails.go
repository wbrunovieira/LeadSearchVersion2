package cnpjsearch

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchCompanyDetailsCNPJBIZ(cnpj string) (map[string]string, error) {
	url := fmt.Sprintf("https://cnpj.biz/%s", cnpj)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the URL: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the page: %v", err)
	}

	details := make(map[string]string)
	details["cnpj"] = extractData(doc, "CNPJ", "p:contains('CNPJ')")
	details["razaoSocial"] = extractData(doc, "Razão Social", "p:contains('Razão Social')")
	details["nomeFantasia"] = extractData(doc, "Nome Fantasia", "p:contains('Nome Fantasia')")
	details["dataAbertura"] = extractData(doc, "Data da Abertura", "p:contains('Data de Abertura')")
	details["telefone"] = extractData(doc, "Telefone(s)", "p:contains('Telefone')")
	details["email"] = extractData(doc, "E-mail", "p:contains('E-mail')")

	return details, nil
}

func extractData(doc *goquery.Document, label string, selector string) string {
	var data string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {

		if s.Text() != "" && strings.Contains(strings.ToLower(s.Text()), strings.ToLower(label)) {
			data = s.Text()
		}
	})
	return data
}
