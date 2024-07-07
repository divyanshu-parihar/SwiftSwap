run :
	go run cmd/main.go


build :
	go build -o bin/SwiftCoin cmd/main.go


migrate :
	go run cmd/main.go migrate
