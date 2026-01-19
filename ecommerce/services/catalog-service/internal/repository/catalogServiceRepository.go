package repository

import (
	"errors"

	"gorm.io/gorm"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
	"github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/catalog-service/internal/domain"
)

type CatalogServiceRepository struct {
	db *gorm.DB
}

func NewCatalogServiceRepository(db *gorm.DB) *CatalogServiceRepository {
	return &CatalogServiceRepository{db: db}
}

// AddCatalogItem adds a new item to the catalog, if the item already exists it returns an error.
func (r *CatalogServiceRepository) AddCatalogItem(item *pb.CatalogItem) error {

	// Check ItemID validity
	if err := checkItemIDValidity(item.ItemId); err != nil {
		return err
	}

	// Check ItemID uniqueness
	if err := checkItemIDUniqueness(item.ItemId, r.db); err != nil {
		return err
	}

	// Check Description validity
	if err := checkDescriptionValidity(item.Description); err != nil {
		return err
	}

	// Check QuantityAvailable validity
	if err := checkQuantityAvailableValidity(item.QuantityAvailable); err != nil {
		return err
	}

	// Check Price validity
	if err := checkPriceValidity(item.Price); err != nil {
		return err
	}

	// Create CatalogItem domain model
	catalogItem := &domain.CatalogItem{
		ItemID:            item.ItemId,
		Description:       item.Description,
		QuantityAvailable: item.QuantityAvailable,
		Price:             item.Price,
	}

	// Save to database
	if err := r.db.Create(catalogItem).Error; err != nil {
		return err
	}

	return nil
}

// RemoveCatalogItem removes a catalog item from the catalog by its unique identifier.
func (r *CatalogServiceRepository) RemoveCatalogItem(itemID string) error {

	// Retrieve item from catalog
	item, err := r.RetrieveCatalogItem(itemID)
	if err != nil {
		return err
	}

	// If the item exists, remove it
	if err := r.db.Delete(&item).Error; err != nil {
		return err
	}

	return nil
}

// GetCatalogItem retrieves a catalog item by its unique identifier.
func (r *CatalogServiceRepository) GetCatalogItem(itemID string) (*pb.CatalogItem, error) {

	// Retrieve item from catalog
	item, err := r.RetrieveCatalogItem(itemID)
	if err != nil {
		return nil, err
	}

	// if the item exists, return it
	protoItem, err := domain.DomainCatalogItemToProtoCatalogItem(item)
	return protoItem, err

}

// UpdateQuantityAvailable updates the quantity available of a catalog item.
func (r *CatalogServiceRepository) UpdateQuantityAvailable(itemID string, quantity uint32) error {

	// Check quantity validity

	// Retrieve item

	// If the item exists, update its quantity available

	// else return error

	return nil
}

// UpdatePrice updates the price of a catalog item.
func (r *CatalogServiceRepository) UpdatePrice(itemID string, price float64) error {

	// Check price validity

	// Retrieve item

	// If the item exists, update its price
	// else return error

	return nil
}

// ListCatalogItems retrieves all catalog items.
func (r *CatalogServiceRepository) ListCatalogItems() ([]*pb.CatalogItem, error) {
	var items []*pb.CatalogItem
	err := r.db.Find(&items).Error
	return items, err
}

// RetrieveCatalogItem retrieves a catalog item by its unique identifier.
func (r *CatalogServiceRepository) RetrieveCatalogItem(itemID string) (*domain.CatalogItem, error) {
	var item domain.CatalogItem
	err := r.db.First(&item, "item_id = ?", itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// PRIVATE FUNCTIONS TO VALIDATE INPUTS

func checkItemIDValidity(itemID string) error {
	if itemID == "" {
		return errors.New("Item ID cannot be empty")
	}
	return nil
}

func checkItemIDUniqueness(itemID string, db *gorm.DB) error {

	var count int64
	db.Model(&domain.CatalogItem{}).Where("item_id = ?", itemID).Count(&count)
	if count > 0 {
		return errors.New("Item ID must be unique")
	}
	return nil
}

func checkDescriptionValidity(description string) error {
	if description == "" {
		return errors.New("Description cannot be empty")
	}
	return nil
}

func checkQuantityAvailableValidity(quantity uint32) error {
	if quantity == 0 {
		return errors.New("Quantity available must be greater than zero")
	}
	return nil
}

func checkPriceValidity(price float64) error {
	if price < 0 {
		return errors.New("Price cannot be negative")
	}
	return nil
}
