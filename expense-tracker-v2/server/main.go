package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"expense-tracker-v2/src/repository"
	"expense-tracker-v2/src/service"
	"expense-tracker-v2/src/transport"

	_ "github.com/lib/pq"
)

const (
	// Update these with your actual Postgres details
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "password" // Change this!
	dbName     = "postgres"
)

func main() {
	// 1. Database Connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify connection
	if err = db.Ping(); err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	// 2. Dependency Injection
	// Init Repository
	expenseRepo := repository.NewPostgresRepository(db)

	// Init Service (injecting repo)
	expenseService := service.NewExpenseService(expenseRepo)

	// Init Transport (injecting service)
	handler := transport.NewHandler(expenseService)

	// 3. Setup Routes
	// Route A: Exact match for "/expenses" (No trailing slash)
	// Handles: Listing (GET) and Creating (POST)
	http.HandleFunc("/expenses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateExpenseHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.ListExpensesHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route B: Prefix match for "/expenses/" (Has trailing slash)
	// Handles: Actions on specific items like /expenses/1 (DELETE)
	http.HandleFunc("/expenses/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			handler.DeleteExpenseHandler(w, r)
		} else if r.Method == http.MethodGet {
			// New: Handle GET for single item
			handler.GetExpenseByIDHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 4. Start Server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
