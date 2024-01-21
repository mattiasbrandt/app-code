package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dapr "github.com/dapr/go-sdk/client"
)

// Data structures
type Individual struct {
	Name           string `json:"name"`
	PersonalNumber string `json:"personalNumber"`
}

type Organization struct {
	Name      string `json:"name"`
	OrgNumber string `json:"orgNumber"`
}

// Validation functions
func (i *Individual) IsValid() bool {
	// Check for non-empty name and 11-digit personal number
	return i.Name != "" && len(i.PersonalNumber) == 11
}

func (o *Organization) IsValid() bool {
	// Check for non-empty name and non-empty org number
	return o.Name != "" && o.OrgNumber != ""
}

func main() {
	// Create a Dapr client
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer client.Close()

	// Define the state store name
	const stateStore = "statestore"

	// Sample data
	ind := Individual{Name: "John Doe", PersonalNumber: "12345678901"}
	org := Organization{Name: "Acme Corp", OrgNumber: "ACME123"}

	// Validate and save Individual
	if ind.IsValid() {
		saveData(client, stateStore, "individual", ind)
	} else {
		fmt.Println("Invalid individual data")
	}

	// Validate and save Organization
	if org.IsValid() {
		saveData(client, stateStore, "organization", org)
	} else {
		fmt.Println("Invalid organization data")
	}

	// Set up a simple HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, this is my Go/Dapr app!")
	})

	// Start listening on a port
	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Save data to state store
func saveData(client dapr.Client, storeName, key string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	// Adding an empty map for metadata
	metadata := map[string]string{}

	if err := client.SaveState(context.Background(), storeName, key, dataBytes, metadata); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	} else {
		fmt.Printf("Data saved for key %s\n", key)
	}
}
