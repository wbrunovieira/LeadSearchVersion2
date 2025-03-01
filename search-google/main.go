package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	// Importe seu pacote googleplaces abaixo (ajuste o path conforme o seu módulo).
	"github.com/wbrunovieira/LeadSearchVersion2/googleplaces"
)

func main() {

	log.Println("Starting the service...")

	// Carrega variáveis de ambiente do .env (opcional, mas útil)
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or couldn't be loaded. Continuing...")
	}

	// Verifica se a API key está presente
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		log.Fatal("API key is required. Set the GOOGLE_PLACES_API_KEY environment variable.")
	}

	// Registra as rotas HTTP
	http.HandleFunc("/start-search", startSearchHandler)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("Starting server on port 8082...")
	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startSearchHandler(w http.ResponseWriter, r *http.Request) {

	// Lê parâmetros da URL: category_id, zipcode_id, radius, max_results
	categoryID := r.URL.Query().Get("category_id")
	zipcodeIDString := r.URL.Query().Get("zipcode_id")
	radiusStr := r.URL.Query().Get("radius")
	maxResultsStr := r.URL.Query().Get("max_results")

	if categoryID == "" || zipcodeIDString == "" || radiusStr == "" {
		http.Error(w, "Missing required parameters (category_id, zipcode_id, radius)", http.StatusBadRequest)
		return
	}

	radiusInt, err := strconv.Atoi(radiusStr)
	if err != nil {
		http.Error(w, "Invalid radius value", http.StatusBadRequest)
		return
	}

	zipcodeID, err := strconv.Atoi(zipcodeIDString)
	if err != nil {
		http.Error(w, "Invalid zipcode_id value", http.StatusBadRequest)
		return
	}

	maxResults := 1 // valor padrão
	if maxResultsStr != "" {
		maxResults, err = strconv.Atoi(maxResultsStr)
		if err != nil {
			http.Error(w, "Invalid max_results value", http.StatusBadRequest)
			return
		}
	}

	// Inicia a busca
	err = startSearch(categoryID, zipcodeID, radiusInt, maxResults)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start search: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Search started for categoryID: %s, zipcodeID: %d, radius: %d",
		categoryID, zipcodeID, radiusInt)
}

func startSearch(categoryID string, zipcodeID, radius, maxResults int) error {

	// Log inicial
	log.Printf("Iniciando pesquisa: categoryID=%s, zipcodeID=%d, radius=%d, maxResults=%d",
		categoryID, zipcodeID, radius, maxResults)

	// Pega API key
	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("API key is required. Set the GOOGLE_PLACES_API_KEY environment variable.")
	}

	// Instancia o serviço GooglePlaces
	service := googleplaces.NewService(apiKey)

	// Converte zipcode (ex: "12345") em coordenadas "lat,lng"
	// Se zipcodeID for o próprio CEP, então:
	zipcodeString := strconv.Itoa(zipcodeID)
	locationStr, err := service.GeocodeZip(zipcodeString)
	if err != nil {
		return fmt.Errorf("erro ao geocodificar o CEP %d: %v", zipcodeID, err)
	}

	log.Printf("Localização obtida para o CEP %s: %s", zipcodeString, locationStr)

	// Define quantas páginas você quer buscar. Aqui, 1 é só exemplo.
	// O Google Places Text Search libera até 20 resultados por página, e máx de 60 no total
	maxPages := 3

	// Realiza a pesquisa. Aqui passamos "categoryID" como o “query”
	// e a coordenada "lat,lng" como “location”
	places, err := service.SearchPlaces(categoryID, locationStr, radius, maxPages)
	if err != nil {
		return fmt.Errorf("erro ao buscar lugares: %v", err)
	}

	// Itera nos resultados, obtendo detalhes de cada PlaceID
	totalLeadsExtracted := 0
	for _, place := range places {
		placeID, ok := place["PlaceID"].(string)
		if !ok {
			log.Println("PlaceID não encontrado ou não é string")
			continue
		}

		details, err := service.GetPlaceDetails(placeID)
		if err != nil {
			log.Printf("Erro ao obter detalhes do place: %v", err)
			continue
		}

		totalLeadsExtracted++
		log.Printf("Lead #%d extraído: %+v", totalLeadsExtracted, details)

		// Se atingiu o limite de resultados
		if totalLeadsExtracted >= maxResults {
			log.Printf("Limite de %d resultados atingido.", maxResults)
			break
		}
	}

	// Log final
	log.Printf("Busca concluída com sucesso! Total de leads: %d", totalLeadsExtracted)
	return nil
}
