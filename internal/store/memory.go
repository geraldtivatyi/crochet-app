package store

import (
	"Crochet/internal/domain"
	"slices"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	products []domain.Product
	stores   []domain.Store
	owners   []domain.StoreOwner
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		owners: []domain.StoreOwner{
			{ID: "1", Name: "Alice", StoreID: "1"},
			{ID: "2", Name: "Bob", StoreID: "2"},
		},
		stores: []domain.Store{
			{ID: "1", Name: "Alice's Crochet Nook"},
			{ID: "2", Name: "Knit & Knot"},
		},
		products: []domain.Product{
			{ID: "1", ProductName: "Summer Hat", ProductPrice: 100.0, StoreID: "1"},
			{ID: "2", ProductName: "Blanket", ProductPrice: 200.0, StoreID: "1"},
		},
	}
}

func (s *MemoryStore) GetProducts() []domain.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return slices.Clone(s.products)
}

func (s *MemoryStore) AddProduct(p domain.Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products = append(s.products, p)
}

func (s *MemoryStore) DeleteProduct(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	initialLen := len(s.products)
	s.products = slices.DeleteFunc(s.products, func(p domain.Product) bool {
		return p.ID == id
	})

	return len(s.products) < initialLen
}
