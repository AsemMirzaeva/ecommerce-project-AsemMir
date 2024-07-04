export POSTGRES_DB=postgres://postgres:1234@localhost:5432/gatewaydb?sslmode=disable

run:
	cd cmd && go run main.go


gen-product:
	protoc \
        --go_out=. \
        --go-grpc_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        ./protos/product-service/product.proto

pmigfile:
	migrate create -ext sql -dir migrations/ product-service

pup:
	migrate -path migrations -database $(POSTGRES_DB) up

pdown:
	migrate -path migrations -database $(POSTGRES_DB) down

pforce:
	migrate -path migrations -database $(POSTGRES_DB) force $(version)	

