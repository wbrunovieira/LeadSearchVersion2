package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Driver PostgreSQL oficial
)

// ConnectDB recebe dados de conexão e retorna *sql.DB
func ConnectDB(host, port, user, password, dbName string) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no PostgreSQL: %v", err)
	}

	// Verifica se a conexão está OK
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro no ping do PostgreSQL: %v", err)
	}

	log.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
	return db, nil
}
