BINARY_NAME=simple_bank

VERSION?=0.0.0

build: 	
	go build -o bin/app 
 
run:  
	build ./bin/app  
 
test: 	
	go test -v -cover ./... -count=1 
 
cpuprof: 
	go test -run=XXX -cpuprofile cpu.prof -bench . 
 
blockprof:   

createdb: 
	docker exec -it postgres16 createdb --username root --owner root simple_bank

dropdb: 
	docker exec -it postgres16 dropdb simple_bank

postgres:
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -d postgres:16-alpine

migrateup:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go

mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/manther/simplebank/db/sqlc Store 
.PHONY: createdb migrateup migratedown sqlc test server mock cpuprof blockprof