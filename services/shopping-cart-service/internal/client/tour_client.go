package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// TourDetails je struktura koja predstavlja odgovor od tour-service
// ID je 'uint' da bi se poklopilo sa odgovorom tour-servisa
type TourDetails struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// TourServiceClient je odgovoran za komunikaciju sa tour-service
type TourServiceClient struct {
	Client  *http.Client
	BaseURL string // Npr. "http://tour-service:8080"
}

//  kreira novu instancu klijenta
func NewTourServiceClient(baseURL string) *TourServiceClient {
	return &TourServiceClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

// GetTourDetails dobavlja detalje o turi na osnovu njenog ID-ja
func (c *TourServiceClient) GetTourDetails(tourID string, authorizationHeader string) (*TourDetails, error) {
	reqURL := fmt.Sprintf("%s/api/v1/tours/%s", c.BaseURL, tourID)

	// Kreiramo novi zahtev da bismo mogli da dodamo header
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	//Dodajemo Authorization header u interni poziv
	req.Header.Set("Authorization", authorizationHeader)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call tour service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tour service returned non-200 status: %d", resp.StatusCode)
	}

	var tourDetails TourDetails
	if err := json.NewDecoder(resp.Body).Decode(&tourDetails); err != nil {
		return nil, fmt.Errorf("failed to decode tour service response: %w", err)
	}

	return &tourDetails, nil
}