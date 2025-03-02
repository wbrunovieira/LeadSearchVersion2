package database

import (
	"database/sql"
	"fmt"
	"log"
)

// Migrate executa a criação/atualização das tabelas necessárias
func Migrate(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS leads (
	  id SERIAL PRIMARY KEY,
	  name VARCHAR(255),
	  formatted_address VARCHAR(255),
	  city VARCHAR(100),
	  state VARCHAR(100),
	  country VARCHAR(100),
	  phone VARCHAR(50),
	  rating DECIMAL(3,1),
	  place_id VARCHAR(100),
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  source VARCHAR(50)  -- Indica origem, ex: 'google_places'
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("erro ao criar tabela leads: %v", err)
	}

	log.Println("Tabela 'leads' criada/verificada com sucesso.")
	return nil
}
