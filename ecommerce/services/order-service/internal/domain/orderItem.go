package domain

type OrderItem struct {

	// ItemID is the unique identifier for the item.
	ItemID string `gorm:"not null; check:item_id <> ''"`

	// Quantity indicates the number of units of the item in the order.
	Quantity uint32 `gorm:"not null; check:quantity > 0"`

	// Price represents the price of a single unit of the item at the time of the order.
	Price float64 `gorm:"not null; check:price >= 0"`
}
