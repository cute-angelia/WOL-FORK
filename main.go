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
	httpPort, computerFilePath := internal.ProcessShellArgs()

	// Process Environment Variables
	httpPort, computerFilePath = internal.ProcessEnvVars(httpPort, computerFilePath)

	// Loading Computer CSV File to Memory File in Memory
	var loadComputerCSVFileError error
	if internal.ComputerList, loadComputerCSVFileError = internal.LoadComputerList(computerFilePath); loadComputerCSVFileError != nil {
		log.Fatalf("Error on loading Computerlist File \"%s\" check File access and formating", computerFilePath)
	}

	// Init HTTP Router - mux
	router := mux.NewRouter()

	// Define Home Route
	router.HandleFunc("/", internal.RenderHomePage).Methods("GET")

	// Define Wakeup Api functions with a Computer Name
	router.HandleFunc("/api/wakeup/computer/{computerName}", internal.RestWakeUpWithComputerName).Methods("GET")
	router.HandleFunc("/api/wakeup/computer/{computerName}/", internal.RestWakeUpWithComputerName).Methods("GET")

	// Define route for adding a new computer
	router.HandleFunc("/api/add/computer", internal.RestAddComputer).Methods("POST")

	// Define route for deleting a computer
	router.HandleFunc("/api/delete/computer/{computerName}", internal.RestDeleteComputer).Methods("DELETE")

	// Setup Webserver
	httpListen := fmt.Sprint(":", httpPort)
	log.Printf("Startup Webserver on \"%s\"", httpListen)

	log.Fatal(http.ListenAndServe(httpListen, handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(router)))
}
