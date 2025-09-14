// api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wbrunovieira/LeadSearchVersion2/db"
	"github.com/wbrunovieira/LeadSearchVersion2/handlers"
	"github.com/wbrunovieira/LeadSearchVersion2/middleware"
	"github.com/wbrunovieira/LeadSearchVersion2/rabbitmq"
)

func main() {
	log.Println("Starting API service...")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT não definido no ambiente")
	}
	fmt.Println("API rodando na porta", port)

	if err := db.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	if err := db.Migrate(); err != nil {
		log.Fatalf("Erro ao migrar o banco de dados: %v", err)
	}
	defer db.Close()

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL não definida no ambiente")
	}
	if err := rabbitmq.InitRabbitMQ(rabbitURL); err != nil {
		log.Fatalf("Erro ao inicializar RabbitMQ: %v", err)
	}
	defer rabbitmq.CloseRabbitMQ()

	mux := http.NewServeMux()
	mux.HandleFunc("/save-leads", handlers.SaveLeadsHandler)
	// mux.HandleFunc("/list-leads", handlers.ListLeadsHandler)
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/update-lead-field", handlers.UpdateLeadHandler)

	handler := middleware.CORS(mux)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
