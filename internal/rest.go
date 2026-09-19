// Rest API Implementations

package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var (
	macRegex = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	ipRegex  = regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
)

// normalizeIP accepts "192.168.1.x" or "192.168.1.x:9" (any port),
// validates the IP part, and always returns "192.168.1.x:9".
func normalizeIP(raw string) (string, bool) {
	host := raw
	if idx := strings.LastIndex(raw, ":"); idx != -1 {
		host = raw[:idx] // strip existing port
	}
	if !ipRegex.MatchString(host) {
		return "", false
	}
	return host + ":9", true
}

// RestWakeUpWithComputerName - REST Handler for Processing URLS /api/wakeup/computer/<computerName>
func RestWakeUpWithComputerName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	computerName := vars["computerName"]

	var result WakeUpResponseObject
	result.Success = false

	if computerName == "" {
		result.Message = "Empty Computername is not allowed"
		result.ErrorObject = nil
		w.WriteHeader(http.StatusBadRequest)
	} else {
		for _, c := range ComputerList {
			if c.Name == computerName {
				if err := SendMagicPacket(c.Mac, c.BroadcastIPAddress, ""); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					result.Success = false
					result.Message = "Internal error on Sending the Magic Packet"
					result.ErrorObject = err
				} else {
					result.Success = true
					result.Message = fmt.Sprintf("Succesfully Wakeup Computer %s with Mac %s on Broadcast IP %s", c.Name, c.Mac, c.BroadcastIPAddress)
					result.ErrorObject = nil
				}
			}
		}

		if !result.Success && result.ErrorObject == nil {
			w.WriteHeader(http.StatusNotFound)
			result.Message = fmt.Sprintf("Computername %s could not be found", computerName)
		}
	}
	json.NewEncoder(w).Encode(result)
}

// RestAddComputer - REST Handler for adding a new computer
func RestAddComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newComputer Computer
	if err := json.NewDecoder(r.Body).Decode(&newComputer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$`)
	if !macRegex.MatchString(newComputer.Mac) {
		http.Error(w, "Invalid MAC address format", http.StatusBadRequest)
		return
	}
	normalizedIP, ok := normalizeIP(newComputer.BroadcastIPAddress)
	if !ok {
		http.Error(w, "Invalid IP address format", http.StatusBadRequest)
		return
	}
	newComputer.BroadcastIPAddress = normalizedIP


	for _, c := range ComputerList {
		if c.Name == newComputer.Name {
			http.Error(w, "Computer with this name already exists", http.StatusBadRequest)
			return
		}
		if c.Mac == newComputer.Mac {
			http.Error(w, "Computer with this MAC address already exists", http.StatusBadRequest)
			return
		}
		if c.BroadcastIPAddress == newComputer.BroadcastIPAddress {
			http.Error(w, "Computer with this IP address already exists", http.StatusBadRequest)
			return
		}
	}

	saved, err := AddComputer(newComputer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ComputerList = append(ComputerList, saved)
	log.Printf("Added new computer: %+v\n", saved)

	json.NewEncoder(w).Encode(WakeUpResponseObject{
		Success:     true,
		Message:     "Computer added successfully",
		ErrorObject: nil,
	})
}

// RestUpdateComputer - REST Handler for updating an existing computer
func RestUpdateComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	oldName := vars["computerName"]

	var updated Computer
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !macRegex.MatchString(updated.Mac) {
		http.Error(w, "Invalid MAC address format", http.StatusBadRequest)
		return
	}
	normalizedIP, ok := normalizeIP(updated.BroadcastIPAddress)
	if !ok {
		http.Error(w, "Invalid IP address format", http.StatusBadRequest)
		return
	}
	updated.BroadcastIPAddress = normalizedIP

	found := false
	for _, c := range ComputerList {
		if c.Name == oldName {
			found = true
			break
		}
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(WakeUpResponseObject{
			Success: false,
			Message: fmt.Sprintf("Computername %s could not be found", oldName),
		})
		return
	}

	for _, c := range ComputerList {
		if c.Name == oldName {
			continue
		}
		if c.Name == updated.Name {
			http.Error(w, "Computer with this name already exists", http.StatusBadRequest)
			return
		}
		if c.Mac == updated.Mac {
			http.Error(w, "Computer with this MAC address already exists", http.StatusBadRequest)
			return
		}
		if c.BroadcastIPAddress == updated.BroadcastIPAddress {
			http.Error(w, "Computer with this IP address already exists", http.StatusBadRequest)
			return
		}
	}

	if err := UpdateComputer(oldName, updated); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var err error
	ComputerList, err = LoadComputerList()
	if err != nil {
		log.Printf("Warning: failed to reload computer list after update: %v", err)
	}

	log.Printf("Updated computer: %s -> %+v\n", oldName, updated)

	json.NewEncoder(w).Encode(WakeUpResponseObject{
		Success:     true,
		Message:     fmt.Sprintf("Computer %s updated successfully", oldName),
		ErrorObject: nil,
	})
}

// RestDeleteComputer - REST Handler for deleting a computer
func RestDeleteComputer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	computerName := vars["computerName"]

	index := -1
	for i, c := range ComputerList {
		if c.Name == computerName {
			index = i
			break
		}
	}

	if index == -1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(WakeUpResponseObject{
			Success:     false,
			Message:     fmt.Sprintf("Computername %s could not be found", computerName),
			ErrorObject: nil,
		})
		return
	}

	if err := DeleteComputer(computerName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ComputerList = append(ComputerList[:index], ComputerList[index+1:]...)
	log.Printf("Deleted computer: %s\n", computerName)

	json.NewEncoder(w).Encode(WakeUpResponseObject{
		Success:     true,
		Message:     fmt.Sprintf("Successfully deleted computer %s", computerName),
		ErrorObject: nil,
	})
}

// RestCheckStatus - REST Handler for checking if a computer is online
func RestCheckStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	computerName := vars["computerName"]

	for _, c := range ComputerList {
		if c.Name == computerName {
			online := CheckOnline(c.BroadcastIPAddress)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"online": online,
				"name":   c.Name,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"online": false,
		"name":   computerName,
		"error":  "computer not found",
	})
}
