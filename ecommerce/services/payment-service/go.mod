module github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/services/payment-service

go 1.25.1

replace github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto => ../../proto

require github.com/MatteoBollecchino/Distributed_Programming_Project/ecommerce/proto v0.0.0-00010101000000-000000000000

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

require (
	github.com/glebarez/sqlite v1.11.0
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260114163908-3f89685c29c3 // indirect
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11 // indirect
	gorm.io/gorm v1.31.1
)
