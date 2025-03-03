package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/wbrunovieira/LeadSearchVersion2/search-google/googleplaces"
)

func main() {
	log.Println("Starting the API service...")

	// Carrega variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or not loaded. Continuing...")
	}
	log.Println(".env file loaded (if present)")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT não definido no ambiente")
	}
	fmt.Println("API rodando na porta", port)

	// Configura as rotas usando um mux
	mux := http.NewServeMux()

	// Rota para iniciar a busca (sem salvar no BD)
	mux.HandleFunc("/start-search", startSearchHandler)

	// Rota de healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Aplica middleware CORS a todas as rotas
	handlerComCORS := withCORS(mux)

	log.Println("Starting server on port", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerComCORS))
}

// startSearchHandler processa a busca via Google Places.
// Agora, essa rota não salva os leads em um banco, mas somente executa a busca
// e loga os resultados.
func startSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	// Lê os parâmetros da URL
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

	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not provided", http.StatusInternalServerError)
		return
	}

	err = startSearch(apiKey, categoryID, zipcodeID, radiusInt, maxResults)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start search: %v", err), http.StatusInternalServerError)
		return
	}

	// Resposta simples informando que a busca foi iniciada.
	fmt.Fprintf(w, "Search started for categoryID: %s, zipcodeID: %d, radius: %d", categoryID, zipcodeID, radiusInt)
}

// startSearch executa a lógica de busca usando o serviço do Google Places.
func startSearch(apiKey string, categoryID string, zipcodeID, radius, maxResults int) error {
	log.Printf("Iniciando pesquisa: categoryID=%s, zipcodeID=%d, radius=%d, maxResults=%d",
		categoryID, zipcodeID, radius, maxResults)

	service := googleplaces.NewService(apiKey)

	zipcodeString := strconv.Itoa(zipcodeID)
	locationStr, err := service.GeocodeZip(zipcodeString)
	if err != nil {
		return fmt.Errorf("erro ao geocodificar o CEP %d: %v", zipcodeID, err)
	}
	log.Printf("Localização obtida para o CEP %s: %s", zipcodeString, locationStr)

	maxPages := 3
	places, err := service.SearchPlaces(categoryID, locationStr, radius, maxPages)
	if err != nil {
		return fmt.Errorf("erro ao buscar lugares: %v", err)
	}

	totalLeadsExtracted := 0
	// Em vez de salvar no banco, apenas logamos os detalhes dos leads.
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
		log.Printf("Lead #%d obtido: %+v", totalLeadsExtracted, details)

		if totalLeadsExtracted >= maxResults {
			log.Printf("Limite de %d resultados atingido.", maxResults)
			break
		}
	}

	log.Printf("Busca concluída com sucesso! Total de leads: %d", totalLeadsExtracted)
	return nil
}

// withCORS é um middleware que adiciona os cabeçalhos CORS a todas as requisições.
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
