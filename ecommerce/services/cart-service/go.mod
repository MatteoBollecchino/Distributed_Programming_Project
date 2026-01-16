module github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/cart-service

go 1.25.1

require github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto v0.0.0

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
)

require (
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11 // indirect
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

replace github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto => ../../proto
