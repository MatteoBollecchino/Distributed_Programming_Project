package domain

type CatalogItem struct {

	// ItemID is the unique identifier for the catalog item.
	ItemID string

	// Description provides details about the catalog item.
	Description string

	// QuantityAvailable indicates how many units of the item are available in stock.
	QuantityAvailable uint32

	// Price indicates the price of the catalog item.
	Price float64
}

/*
// DomainUserToProtoUser converts a model.User into a pb.User
func DomainUserToProtoUser(user *User) (*pb.User, error) {
	if user == nil {
		return nil, fmt.Errorf("Input argument is nil")
	}

	var r pb.Role
	if user.Role == AdminRole {
		r = pb.Role_ADMIN
	} else {
		r = pb.Role_USER
	}
	return &pb.User{
		Username: user.Username,
		Password: user.Password,
		Role:     r}, nil
}*/
