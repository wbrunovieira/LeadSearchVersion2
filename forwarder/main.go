// forwarder/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wbrunovieira/LeadSearchVersion2/forwarder/rabbitmq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado")
	}

	rabbitmq.Connect()
	defer rabbitmq.Conn.Close()
	defer rabbitmq.Ch.Close()

	go rabbitmq.ConsumeQueue()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Forwarder service is running"))
	})
	log.Printf("Forwarder service rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
