package account

type OwnerInfo struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

type UsageRules struct {
	AllowedFromAccount   []string `json:"allowedFromAccount"`
	AllowedToAccount     []string `json:"allowedToAccount"`
	ForbiddenFromAccount []string `json:"forbiddenFromAccount"`
	ForbiddenToAccount   []string `json:"forbiddenToAccount"`
}

type Account struct {
	AdditionalData  map[string]interface{} `json:"additionalData"`
	Balance         int64                  `json:"balance"`
	BlockedBalance  int64                  `json:"blockedBalance"`
	Category        string                 `json:"category"`
	Channel         string                 `json:"channel"`
	CreatedAt       interface{}            `json:"createdAt"`
	Currency        string                 `json:"currency"`
	ID              string                 `json:"id"`
	LastMovementAt  interface{}            `json:"lastMovementAt"`
	MaxBalance      int64                  `json:"maxBalance"`
	OwnerID         string                 `json:"ownerId"`
	OwnerInfo       OwnerInfo              `json:"ownerInfo"`
	OwnerNationalID string                 `json:"ownerNationalId"`
	OwnerType       string                 `json:"ownerType"`
	Status          string                 `json:"status"`
	TransactionID   string                 `json:"transactionId"`
	Type            string                 `json:"type"`
	UpdatedAt       interface{}            `json:"updatedAt"`
	UsageRules      UsageRules             `json:"usageRules"`
}
