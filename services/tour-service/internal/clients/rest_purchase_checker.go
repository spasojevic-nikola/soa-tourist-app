// PRIVREMENO REST KLIJENT -> ZAMIJENI SA RPC KLIJENTOM
package clients

import (
    "encoding/json"
    "fmt"
    "net/http"
    "tour-service/internal/service"
)

type RESTPurchaseChecker struct {
    baseURL string
    client  *http.Client
}

func NewRESTPurchaseChecker(baseURL string) service.PurchaseChecker {
    return &RESTPurchaseChecker{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (r *RESTPurchaseChecker) HasPurchasedTour(touristID uint, tourID uint) (bool, error) {
    url := fmt.Sprintf("%s/api/purchase/verify/%d/%d", r.baseURL, touristID, tourID)
    
    resp, err := r.client.Get(url)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("purchase service returned status: %d", resp.StatusCode)
    }

    var result struct {
        HasPurchased bool `json:"hasPurchased"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return false, err
    }

    return result.HasPurchased, nil
}