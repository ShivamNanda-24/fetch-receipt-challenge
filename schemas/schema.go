package schemas

// Receipt structure
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// Item structure
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// Response structure for /receipts/process endpoint
type ProcessResponse struct {
	ID string `json:"id"`
}

// Response structure for /receipts/{id}/points endpoint
type PointsResponse struct {
	Points int `json:"points"`
}
