package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/keyvault/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
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
	return i.Name != "" && len(i.PersonalNumber) == 11
}

func (o *Organization) IsValid() bool {
	return o.Name != "" && o.OrgNumber != ""
}

// Azure Key Vault interaction
func getKeyVaultToken(vaultName, secretName string) (string, error) {
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return "", err
	}

	basicClient := keyvault.New()
	basicClient.Authorizer = authorizer

	vaultURL := fmt.Sprintf("https://%s.vault.azure.net/", vaultName)
	secretBundle, err := basicClient.GetSecret(context.Background(), vaultURL, secretName, "")
	if err != nil {
		return "", err
	}

	return *secretBundle.Value, nil
}

// Cosmos DB Client
func getCosmosDBClient(accountKey string) *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

// Store data in Cosmos DB
func storeDataInCosmosDB(client *http.Client, dbURL string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}

	req, err := http.NewRequest("POST", dbURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// Add necessary headers here
	// For example, req.Header.Add("Authorization", "type=master&ver=1.0&sig="+accountKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Handle the response
	// ...
}

func main() {
	// Create a Dapr client
	daprClient, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}
	defer daprClient.Close()

	// Define the state store name
	const stateStore = "statestore"

	// Sample data
	ind := Individual{Name: "John Doe", PersonalNumber: "12345678901"}
	org := Organization{Name: "Acme Corp", OrgNumber: "ACME123"}

	// Validate and save Individual
	if ind.IsValid() {
		saveData(daprClient, stateStore, "individual", ind)
	} else {
		fmt.Println("Invalid individual data")
	}

	// Validate and save Organization
	if org.IsValid() {
		saveData(daprClient, stateStore, "organization", org)
	} else {
		fmt.Println("Invalid organization data")
	}

	// Retrieve tokens and Cosmos DB credentials from Azure Key Vault
	token, err := getKeyVaultToken("YOUR_VAULT_NAME", "YOUR_SECRET_NAME")
	if err != nil {
		log.Fatalf("Error getting secret from Key Vault: %v", err)
	}

	accountName, err := getKeyVaultToken("YOUR_VAULT_NAME", "COSMOS_DB_ACCOUNT_NAME_SECRET")
	if err != nil {
		log.Fatalf("Error getting Cosmos DB account name: %v", err)
	}

	accountKey, err := getKeyVaultToken("YOUR_VAULT_NAME", "COSMOS_DB_ACCOUNT_KEY_SECRET")
	if err != nil {
		log.Fatalf("Error getting Cosmos DB account key: %v", err)
	}

	cosmosClient := getCosmosDBClient(accountKey)

	// Store validated data in Cosmos DB
	if ind.IsValid() {
		storeDataInCosmosDB(cosmosClient, "COSMOS_DB_URL", ind)
	}
	if org.IsValid() {
		storeDataInCosmosDB(cosmosClient, "COSMOS_DB_URL", org)
	}
}

// Save data to Dapr state store
func saveData(client dapr.Client, storeName, key string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal data: %v", err)
	}

	metadata := map[string]string{} // Adding an empty map for metadata

	if err := client.SaveState(context.Background(), storeName, key, dataBytes, metadata); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	} else {
		fmt.Printf("Data saved for key %s in Dapr state store\n", key)
	}
}
