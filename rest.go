// Rest API Implementations

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

// restWakeUpWithComputerName - REST Handler for Processing URLS /api/computer/<computerName>
func restWakeUpWithComputerName(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	computerName := vars["computerName"]

	var result WakeUpResponseObject
	result.Success = false

	// Ensure computerName is not empty
	if computerName == "" {
		result.Message = "Empty Computername is not allowed"
		result.ErrorObject = nil
		w.WriteHeader(http.StatusBadRequest)
		// Computername is empty
	} else {

		// Get Computer from List
		for _, c := range ComputerList {
			if c.Name == computerName {

				// We found the Computername
				if err := SendMagicPacket(c.Mac, c.BroadcastIPAddress, ""); err != nil {
					// We got an internal Error on SendMagicPacket
					w.WriteHeader(http.StatusInternalServerError)
					result.Success = false
					result.Message = "Internal error on Sending the Magic Packet"
					result.ErrorObject = err
				} else {
					// Horray we send the WOL Packet succesfully
					result.Success = true
					result.Message = fmt.Sprintf("Succesfully Wakeup Computer %s with Mac %s on Broadcast IP %s", c.Name, c.Mac, c.BroadcastIPAddress)
					result.ErrorObject = nil
				}
			}
		}

		if result.Success == false && result.ErrorObject == nil {
			// We could not find the Computername
			w.WriteHeader(http.StatusNotFound)
			result.Message = fmt.Sprintf("Computername %s could not be found", computerName)
		}
	}
	json.NewEncoder(w).Encode(result)
}

// restAddComputer - REST Handler for adding a new computer
func restAddComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newComputer Computer
	if err := json.NewDecoder(r.Body).Decode(&newComputer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Append the new computer to the list
	ComputerList = append(ComputerList, newComputer)

	// Save the updated list to the CSV file
	if err := saveComputerList(DefaultComputerFilePath, ComputerList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success message
	response := WakeUpResponseObject{
		Success: true,
		Message: "Computer added successfully",
		ErrorObject: nil,
	}
	json.NewEncoder(w).Encode(response)
}
