package cnpjsearch

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func FetchDataCNPJBIZ(companyName string, cityName string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	searchQuery := strings.ReplaceAll(companyName, " ", "%20")
	url := fmt.Sprintf("https://cnpj.biz/procura/%s", searchQuery)
	log.Printf("url do FetchDataCNPJBIZ : %v", url)

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

	found := false

	doc.Find(".flex.items-center.text-sm.text-gray-500").EachWithBreak(func(index int, element *goquery.Selection) bool {
		cnpjData := strings.TrimSpace(element.Text())
		if strings.Contains(cnpjData, "/") {
			cnpjParts := strings.Split(cnpjData, "/")
			if len(cnpjParts) >= 2 {
				cnpjCity := strings.TrimSpace(cnpjParts[1])
				if strings.Contains(strings.ToLower(cnpjCity), strings.ToLower(cityName)) {
					cnpj := strings.TrimSpace(cnpjParts[0])

					// Captura os detalhes da empresa utilizando a função refatorada
					details, err := FetchCompanyDetailsCNPJBIZ(cnpj)
					if err == nil {
						result["details"] = details
					} else {
						result["details_error"] = err.Error()
					}

					result["cnpj"] = cnpj
					result["city"] = cnpjCity
					found = true
					return false // interrompe o loop
				}
			}
		}
		return true // continua o loop
	})

	SleepRandom()

	if !found {
		return nil, fmt.Errorf("no CNPJ data found for company %s in city %s", companyName, cityName)
	}
	return result, nil
}

func SleepRandom() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(2+rng.Intn(2)) * time.Second)
}
