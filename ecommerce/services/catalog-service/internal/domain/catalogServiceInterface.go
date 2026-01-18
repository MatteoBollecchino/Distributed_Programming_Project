package domain

type CatalogServiceInterface interface {

	// AddCatalogItem adds a new catalog item to the catalog.
	AddCatalogItem(item *CatalogItem) error

	// RemoveCatalogItem removes a catalog item from the catalog by its unique identifier.
	RemoveCatalogItem(itemID string) error

	// GetCatalogItem retrieves a catalog item by its unique identifier.
	GetCatalogItem(itemID string) (*CatalogItem, error) //DA METTERE pb.

	// UpdateQuantityAvailable updates the quantity available of a catalog item.
	UpdateQuantityAvailable(itemID string, quantity uint32) error

	// UpdatePrice updates the price of a catalog item.
	UpdatePrice(itemID string, price float64) error

	// ListCatalogItems retrieves all catalog items.
	ListCatalogItems() ([]*CatalogItem, error)
}
