package main

import (
	"fmt"
	"go-rest-wol/internal"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Start Processing Shell Arguments or use Default Values defined in const.go
	httpPort, dbPath := internal.ProcessShellArgs()

	// Process Environment Variables
	httpPort, dbPath = internal.ProcessEnvVars(httpPort, dbPath)

	// Initialize SQLite database
	if err := internal.InitDB(dbPath); err != nil {
		log.Fatalf("Error initializing SQLite database \"%s\": %v", dbPath, err)
	}

	// Load Computer List from DB into memory
	var loadErr error
	if internal.ComputerList, loadErr = internal.LoadComputerList(); loadErr != nil {
		log.Fatalf("Error loading computer list from database: %v", loadErr)
	}

	// Init HTTP Router - mux
	router := mux.NewRouter()

	// Define Home Route
	router.HandleFunc("/", internal.RenderHomePage).Methods("GET")

	// Wakeup API
	router.HandleFunc("/api/wakeup/computer/{computerName}", internal.RestWakeUpWithComputerName).Methods("GET")
	router.HandleFunc("/api/wakeup/computer/{computerName}/", internal.RestWakeUpWithComputerName).Methods("GET")

	// Add computer
	router.HandleFunc("/api/add/computer", internal.RestAddComputer).Methods("POST")

	// Update computer
	router.HandleFunc("/api/update/computer/{computerName}", internal.RestUpdateComputer).Methods("PUT")

	// Delete computer
	router.HandleFunc("/api/delete/computer/{computerName}", internal.RestDeleteComputer).Methods("DELETE")

	// Setup Webserver
	httpListen := fmt.Sprint(":", httpPort)
	log.Printf("Startup Webserver on \"%s\"", httpListen)

	log.Fatal(http.ListenAndServe(httpListen, handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router)))
}
