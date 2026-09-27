package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nabinagrawal64/olx-api/internal/config"
) 

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down]")
		return
	}

	cfg := config.MustLoad()

	m, err := migrate.New( 
        "file://migrations",
        cfg.DatabaseUrl,
	)
	if err != nil {
		log.Fatalf("Migration error: %v", err)
	}

	cmd := os.Args[1]

	switch cmd {
	case "up":
		if err := m.Up(); err!=nil {
			log.Fatalf("Migration up error: %v", err)
		}
		fmt.Println("Migration up successfully")
	case "down":
		if err := m.Down(); err!=nil {
			log.Fatalf("Migration down error: %v", err)
		}
		fmt.Println("Migration down successfully")
	default:	
		fmt.Println("Usage: go run cmd/migrate/main.go [up|down]")
		return
	}
}