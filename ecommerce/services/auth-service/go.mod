module github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/auth-service

go 1.25.1

// replace github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto => ../proto

require (
	github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto v0.0.0-20260114162401-55bee74f3b5d // indirect
	github.com/mattn/go-sqlite3 v1.14.33 // indirect
	golang.org/x/crypto v0.47.0
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/text v0.33.0 // indirect
	gorm.io/driver/sqlite v1.6.0
	gorm.io/gorm v1.31.1
)

// replace github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto => ./proto
