// /api/main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/wbrunovieira/LeadSearchVersion2/db"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/save-leads", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint /save-leads acessado de %s usando o método %s", r.RemoteAddr, r.Method)

		if r.Method != http.MethodPost {
			log.Printf("Método inválido %s. Apenas POST é permitido.", r.Method)
			http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
			return
		}

		var leadsData []map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&leadsData)
		if err != nil {
			log.Printf("Erro ao decodificar JSON: %v", err)
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		log.Printf("Recebidos %d leads para salvar", len(leadsData))

		for i, data := range leadsData {
			log.Printf("Processando lead #%d: %+v", i+1, data)
			err = saveLead(data)
			if err != nil {
				log.Printf("Erro ao salvar lead #%d: %v", i+1, err)
				http.Error(w, fmt.Sprintf("Falha ao salvar um lead: %v", err), http.StatusInternalServerError)
				return
			}
			log.Printf("Lead #%d salvo com sucesso", i+1)
		}

		log.Println("Todos os leads foram salvos com sucesso!")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Leads salvos com sucesso!"))
	})

	mux.HandleFunc("/list-leads", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Endpoint /list-leads acessado de %s usando o método %s", r.RemoteAddr, r.Method)

		if r.Method != http.MethodGet {
			log.Printf("Método inválido %s. Apenas GET é permitido.", r.Method)
			http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
			return
		}

		leads, err := db.GetLeads()
		if err != nil {
			log.Printf("Erro ao buscar leads: %v", err)
			http.Error(w, fmt.Sprintf("Falha ao buscar leads: %v", err), http.StatusInternalServerError)
			return
		}

		// Loga cada lead encontrado
		for i, lead := range leads {
			log.Printf("Lead #%d: %+v", i+1, lead)
		}

		jsonResponse, err := json.Marshal(leads)
		if err != nil {
			log.Printf("Erro ao converter leads para JSON: %v", err)
			http.Error(w, fmt.Sprintf("Falha ao converter dados: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
		log.Printf("Retornados %d leads com sucesso.", len(leads))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := withCORS(mux)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func saveLead(placeDetails map[string]interface{}) error {
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

	if v, ok := placeDetails["Radius"].(float64); ok {
		lead.Radius = int(v)
	} else if v, ok := placeDetails["Radius"].(int); ok {
		lead.Radius = v
	}

	if v, ok := placeDetails["Category"].(string); ok {
		lead.Category = v
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
	log.Printf("lead a enviar no banco de dados: %v", &lead)

	err := db.CreateLead(&lead)
	if err != nil {
		log.Printf("Erro ao salvar lead no banco de dados: %v", err)
		return fmt.Errorf("failed to save lead to database: %v", err)
	}
	log.Printf("Lead salvo no banco de dados: %+v", lead)
	return nil
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
