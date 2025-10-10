package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tour-service/internal/interfaces" 
)


//  je klijent specificno za  shopping-cart-service
type RESTShoppingCartClient struct {
	Client  *http.Client
	BaseURL string //"http://shopping-cart-service:8081"
}

func NewShoppingCartClient(baseURL string) interfaces.PurchaseChecker { 
	return &RESTShoppingCartClient{
		Client:  &http.Client{},
		BaseURL: baseURL,
	}
}

//  proverava da li korisnik ima kupljenu turu
func (c *RESTShoppingCartClient) HasUserPurchasedTour(userID, tourID uint, authorizationHeader string) (bool, error) {
	// Kreiramo URL koji gađa tvoj shopping-cart-service endpoint
	reqURL := fmt.Sprintf("%s/api/v1/cart/purchase-status/%d", c.BaseURL, tourID)

	// Kreiramo novi zahtev kako bismo mogli da dodamo Authorization header
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	// Prosledjujemo autorizaciju originalnog korisnika
	req.Header.Set("X-User-ID", strconv.FormatUint(uint64(userID), 10))

	resp, err := c.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to call shopping cart service: %w", err)
	}
	defer resp.Body.Close()
	fmt.Printf(">>> SHOPPING-CART-SERVICE ODGOVORIO SA STATUSOM: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var result struct {
		IsPurchased bool `json:"isPurchased"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("failed to decode shopping cart response: %w", err)
	}

	return result.IsPurchased, nil
}