package domain

type Product struct {
	ID           string  `json:"id"`
	ProductName  string  `json:"productName"`
	ProductPrice float64 `json:"productPrice"`
	StoreID      string  `json:"storeId"`
}

type Store struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type StoreOwner struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	StoreID string `json:"storeId"`
}
