package service

type PurchaseChecker interface {
    HasPurchasedTour(touristID uint, tourID uint) (bool, error)
}

// PRIVREMENO DOK SE NE SREDI GATEWAY I RPC