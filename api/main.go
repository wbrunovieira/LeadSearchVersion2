package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/wbrunovieira/LeadSearchVersion2/db" // importe o pacote db (aqui o módulo foi definido como "github.com/wbrunovieira/LeadSearchVersion2")
	// seu outro pacote
)

func main() {
	log.Println("Starting API service...")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT não definido no ambiente")
	}
	fmt.Println("API rodando na porta", port)

	// Conecta ao banco de dados usando o pacote db.
	if err := db.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/save-leads", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
			return
		}

		// Decodifica o corpo JSON em um array de maps
		var leadsData []map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&leadsData)
		if err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		// Itera sobre cada lead recebido e tenta salvar
		for _, data := range leadsData {
			err = saveLead(data)
			if err != nil {
				log.Printf("Erro ao salvar lead: %v", err)
				http.Error(w, fmt.Sprintf("Falha ao salvar um lead: %v", err), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Leads salvos com sucesso!"))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func saveLead(placeDetails map[string]interface{}) error {
	// Use o tipo Lead do pacote db, referenciando o pacote "db" e não a variável.
	lead := db.Lead{}

	if v, ok := placeDetails["Name"].(string); ok {
		lead.BusinessName = v
	}
	if v, ok := placeDetails["FormattedAddress"].(string); ok {
		lead.Address = v
	}
	if v, ok := placeDetails["City"].(string); ok {
		lead.City = v
	}
	if v, ok := placeDetails["State"].(string); ok {
		lead.State = v
	}
	if v, ok := placeDetails["ZIPCode"].(string); ok {
		lead.ZIPCode = v
	}
	if v, ok := placeDetails["Country"].(string); ok {
		lead.Country = v
	}
	if v, ok := placeDetails["InternationalPhoneNumber"].(string); ok {
		log.Printf("Verificando WhatsApp para o telefone: %s", v)
		lead.Phone = v
		log.Println("WhatsApp confirmado e salvo.")
	}
	if v, ok := placeDetails["Email"].(string); ok {
		log.Printf("Validando email: %s", v)
		lead.Email = v
		log.Println("Email salvo.")
	}
	if v, ok := placeDetails["Website"].(string); ok {
		lead.Website = v
		if strings.HasPrefix(lead.Website, "https://www.instagram.com") {
			lead.Instagram = lead.Website
			lead.Website = ""
			log.Println("Detectado Instagram, atualizado corretamente.")
		}
		if strings.HasPrefix(lead.Website, "https://www.facebook.com") {
			lead.Facebook = lead.Website
			lead.Website = ""
			log.Println("Detectado Facebook, atualizado corretamente.")
		}
	}
	if v, ok := placeDetails["Description"].(string); ok {
		lead.Description = v
		log.Println("Descrição atualizada.")
	}
	if v, ok := placeDetails["Rating"].(float64); ok {
		lead.Rating = v
	}
	if v, ok := placeDetails["UserRatingsTotal"].(float64); ok {
		lead.UserRatingsTotal = int(v)
	}
	if v, ok := placeDetails["PriceLevel"].(float64); ok {
		lead.PriceLevel = int(v)
	}
	if v, ok := placeDetails["BusinessStatus"].(string); ok {
		lead.BusinessStatus = v
	}
	if v, ok := placeDetails["Vicinity"].(string); ok {
		lead.Vicinity = v
	}
	if v, ok := placeDetails["PermanentlyClosed"].(bool); ok {
		lead.PermanentlyClosed = v
	}
	if v, ok := placeDetails["Types"].([]interface{}); ok {
		var types []string
		for _, t := range v {
			if typeStr, ok := t.(string); ok {
				types = append(types, typeStr)
			}
		}
		lead.Categories = strings.Join(types, ", ")
	}
	if v, ok := placeDetails["PlaceID"].(string); ok {
		lead.GoogleId = v
	}

	lead.Source = "GooglePlaces"

	// Usa a função CreateLead do pacote db para salvar o lead.
	err := db.CreateLead(&lead)
	if err != nil {
		log.Printf("Erro ao salvar lead no banco de dados: %v", err)
		return fmt.Errorf("failed to save lead to database: %v", err)
	}
	log.Printf("Lead salvo no banco de dados: %+v", lead)
	return nil
}
