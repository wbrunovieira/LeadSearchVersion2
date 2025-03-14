// api/handlers/handlers.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/wbrunovieira/LeadSearchVersion2/db"
	"github.com/wbrunovieira/LeadSearchVersion2/rabbitmq"
)

func SaveLeadsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Endpoint /save-leads acessado de %s usando o método %s", r.RemoteAddr, r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	var leadsData []map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&leadsData); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	log.Printf("Recebidos %d leads para salvar", len(leadsData))

	for i, data := range leadsData {
		log.Printf("Processando lead #%d: %+v", i+1, data)
		lead, err := saveLead(data)
		if err != nil {
			http.Error(w, fmt.Sprintf("Falha ao salvar um lead: %v", err), http.StatusInternalServerError)
			return
		}

		if err := rabbitmq.PublishLead(lead); err != nil {
			log.Printf("Falha ao publicar o lead %+v no RabbitMQ: %v", lead, err)

		}
		log.Printf("Lead #%d salvo e publicado com sucesso", i+1)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Leads salvos com sucesso!"))
}

func ListLeadsHandler(w http.ResponseWriter, r *http.Request) {
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
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func saveLead(placeDetails map[string]interface{}) (*db.Lead, error) {
	lead := db.Lead{}

	lead.ID = uuid.New()

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
		lead.Phone = v
	}
	if v, ok := placeDetails["Email"].(string); ok {
		lead.Email = v
	}
	if v, ok := placeDetails["Website"].(string); ok {
		lead.Website = v
		if strings.HasPrefix(lead.Website, "https://www.instagram.com") {
			lead.Instagram = lead.Website
			lead.Website = ""
		}
		if strings.HasPrefix(lead.Website, "https://www.facebook.com") {
			lead.Facebook = lead.Website
			lead.Website = ""
		}
	}
	if v, ok := placeDetails["Description"].(string); ok {
		lead.Description = v
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

	if err := db.CreateLead(&lead); err != nil {
		return nil, fmt.Errorf("failed to save lead to database: %v", err)
	}
	log.Printf("Lead salvo no banco de dados: %+v", lead)
	// Retorne o ID do lead salvo
	return &lead, nil
}
