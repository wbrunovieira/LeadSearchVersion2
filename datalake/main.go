package main

import (
	"log"
	"net/http"
	"os"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func main() {

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL},
	})
	if err != nil {
		log.Fatalf("Erro ao criar o cliente Elasticsearch: %v", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Erro ao obter informações do Elasticsearch: %v", err)
	}
	defer res.Body.Close()
	log.Println("Cliente Elasticsearch criado com sucesso.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Datalake service is running"))
	})
	log.Printf("Datalake service rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
