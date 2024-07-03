export POSTGRES_DB=postgres://postgres:1234@localhost:5432/gatewaydb?sslmode=disable

.PHONY: user-gen umigfile uup udown uforce tmigfile tup tdown tforce

user-gen:
	@protoc \
        --go_out=. \
        --go-grpc_out=. \
        --go_opt=paths=source_relative \
        --go-grpc_opt=paths=source_relative \
        ./proto/user.proto

umigfile:
	cd user && migrate create -ext sql -dir migrations/ user

uup:
	cd user && migrate -path migrations -database $(POSTGRES_DB) up

udown:
	cd user && migrate -path migrations -database $(POSTGRES_DB) down

uforce:
	cd user && migrate -path migrations -database $(POSTGRES_DB) force $(version)	

