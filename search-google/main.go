package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wbrunovieira/LeadSearchVersion2/search-google/database"
	"github.com/wbrunovieira/LeadSearchVersion2/search-google/googleplaces"
)

func main() {

	log.Println("Starting the service...")

	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found or not loaded. Continuing...")
	}
	log.Println(".env file loaded (if present)")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("Faltam variáveis de ambiente para conexão com PostgreSQL (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME).")
	}

	db, err := database.ConnectDB(dbHost, dbPort, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatalf("Falha na conexão com o DB: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Falha ao rodar migrações: %v", err)
	}

	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		log.Fatal("API key is required. Set the GOOGLE_PLACES_API_KEY environment variable.")
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/start-search", func(w http.ResponseWriter, r *http.Request) {
		startSearchHandler(w, r, db, apiKey)
	})

	mux.HandleFunc("/list-leads", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método não suportado", http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT name, formatted_address, place_id FROM leads ORDER BY id DESC")
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao buscar leads: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var results []map[string]interface{}
		for rows.Next() {
			var name, address, placeID string
			if err := rows.Scan(&name, &address, &placeID); err != nil {
				log.Println("Erro no Scan de leads:", err)
				continue
			}
			results = append(results, map[string]interface{}{
				"name":              name,
				"formatted_address": address,
				"place_id":          placeID,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handlerComCORS := withCORS(mux)

	log.Println("Starting server on port 8082...")
	if err := http.ListenAndServe(":8082", handlerComCORS); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func startSearchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, apiKey string) {
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

	err = startSearch(db, apiKey, categoryID, zipcodeID, radiusInt, maxResults)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start search: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Search started for categoryID: %s, zipcodeID: %d, radius: %d", categoryID, zipcodeID, radiusInt)
}

func startSearch(db *sql.DB, apiKey string, categoryID string, zipcodeID, radius, maxResults int) error {
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

		err = saveLead(db, details)
		if err != nil {
			log.Printf("Falha ao salvar lead no DB: %v", err)
		} else {
			log.Printf("Lead #%d salvo no DB com sucesso!", totalLeadsExtracted)
		}

		if totalLeadsExtracted >= maxResults {
			log.Printf("Limite de %d resultados atingido.", maxResults)
			break
		}
	}

	log.Printf("Busca concluída com sucesso! Total de leads: %d", totalLeadsExtracted)
	return nil
}

func saveLead(db *sql.DB, placeDetails map[string]interface{}) error {
	query := `
	INSERT INTO leads (
		name,
		formatted_address,
		city,
		state,
		country,
		phone,
		rating,
		place_id,
		source
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'google_places')
	`

	_, err := db.Exec(query,
		placeDetails["Name"],
		placeDetails["FormattedAddress"],
		placeDetails["City"],
		placeDetails["State"],
		placeDetails["Country"],
		placeDetails["InternationalPhoneNumber"],
		placeDetails["Rating"],
		placeDetails["PlaceID"],
	)
	if err != nil {
		return fmt.Errorf("erro no INSERT: %v", err)
	}
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
