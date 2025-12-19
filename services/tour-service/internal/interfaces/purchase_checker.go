package interfaces

// PurchaseChecker defines the interface for checking purchases.
type PurchaseChecker interface {
	HasUserPurchasedTour(userID, tourID uint, authorizationHeader string) (bool, error)
}