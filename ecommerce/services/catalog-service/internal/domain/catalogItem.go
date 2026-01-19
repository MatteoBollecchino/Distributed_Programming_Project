package domain

import (
	"fmt"

	pb "github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto/catalog"
)

type CatalogItem struct {

	// ItemID is the unique identifier for the catalog item.
	ItemID string `gorm:"primaryKey;not null; check:item_id <> ''"`

	// Description provides details about the catalog item.
	Description string `gorm:"not null; check:description <> ''"`

	// QuantityAvailable indicates how many units of the item are available in stock.
	QuantityAvailable uint32 `gorm:"not null; check:quantity_available > 0"`

	// Price indicates the price of the catalog item.
	Price float64 `gorm:"not null; check:price >= 0"`
}

// DomainCatalogItemToProtoCatalogItem converts a model.CatalogItem into a pb.CatalogItem
func DomainCatalogItemToProtoCatalogItem(item *CatalogItem) (*pb.CatalogItem, error) {
	if item == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	return &pb.CatalogItem{
		ItemId:            item.ItemID,
		Description:       item.Description,
		QuantityAvailable: item.QuantityAvailable,
		Price:             item.Price,
	}, nil
}
