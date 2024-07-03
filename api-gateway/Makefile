export POSTGRES_DB=postgres://postgres:1234@localhost:5432/examdb?sslmode=disable

.PHONY: services-gen umigfile uup udown uforce tmigfile tup tdown tforce

gateway-gen:
	@protoc \
        --go_out=. \
        --go-grpc_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        ./proto/gateway.proto

umigfile:
	cd internal && migrate create -ext sql -dir migrations/ user

uup:
	cd internal && migrate -path migrations -database $(POSTGRES_DB) up

udown:
	cd internal && migrate -path migrations -database $(POSTGRES_DB) down

uforce:
	cd internal && migrate -path migrations -database $(POSTGRES_DB) force $(version)	

