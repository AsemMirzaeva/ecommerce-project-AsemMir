create an e-commerce system using go micro services grpc


use grpc-gateway

connect every part to databases of each other using migrations use uuid
use loggers interceptors
config
gracefully shutting down
proper error handling
use .env for config

POSTGRES_HOST = localhost
POSTGRES_PORT = 5432
POSTGRES_DATABASE = gatewaydb
POSTGRES_USER = postgres
POSTGRES_PASSWORD = 1234

DATABASE_URL = postgres://postgres:1234@localhost:5432/gatewaydb

USERSERVICE_HOST = localhost
USERSERVICE_PORT = 


APIGATEWAY_HOST= localhost
APIGATEWAY_PORT= 


1.flow of creating users
client calls CreateUser through apigateway
apigateway sends request to user service
user service process the info and records the data in database
then user service sends the response to apigateway

first we use user service

2.flow of creating product
client calls AddProduct through api gateway
apigateway sends the request to Product service
product service processes the data and records the data of products to database
product service sends the response to apigateway and it sends that response to client

second the created user wants to get product 

3.flow of creating order
client calls CreateOrder through api gateway
apigateway sends the request to Order service
order service calls get user from user service  so it could get data from user
order service calls get product from product service  so it could get data from user
order service processes the data and records the data of products to database
order service sends the response to apigateway and it sends that response to client

user orders

orders are created through stream with being shown in protos stream key word

4.ListProducts(streaming)
client calls list products on api gateway
api gateway sends the request to product service
product service starts streaming so it could get all the list of products

use clientstream or servicestream

product service sends all the products to api gateway, and it send the data to a client



5.CreateOrders(streaming)
client calls create orders on api gateway
api gateway sends the request to order service
order service starts streaming so it could get all the list of products

use clientstream or servicestream

order service process the request and records product data to database


order service sends all the products to api gateway, and it sends the data to a client





user's order is written on the list 

lists are recorded through stream with being shown in protos stream key word


ecommerce-project/
├── api-gateway/
│   ├── cmd/
│   │   └── main.go                   // Main entry point for the API Gateway
│   ├── internal/
│   │   ├── config/                   // Configuration package
│   │   │   └── config.go             // Configuration handling
│   │   ├── handler/                  // Request handlers for API Gateway
│   │   │   └── handler.go            // Handlers for API endpoints
│   │   ├── interceptor/              // gRPC interceptors for API Gateway
│   │   │   └── interceptor.go        // Interceptors for API Gateway
│   │   ├── middleware/               // Middleware for API Gateway
│   │   │   └── middleware.go         // Middleware functions
│   │   ├── repository/               // Repository layer for API Gateway
│   │   │   └── repository.go         // Data access layer
│   │   └── proto/                    // gRPC proto definitions
│   │       ├── gateway.proto         // Proto file defining gRPC services
│   │       └── gateway.pb.go         // Generated Go code from proto
│   ├── migration/                    // Database migration scripts
│   │   └── 001_create_gateway_table.up.sql  // Example migration script
│   └── config/                       // Configuration files
│       └── config.yaml               // Configuration specific to API Gateway
├── product-service/
│   ├── cmd/
│   │   └── main.go                   // Main entry point for Product Service
│   ├── internal/
│   │   ├── config/                   // Configuration package
│   │   │   └── config.go             // Configuration handling
│   │   ├── handler/                  // Request handlers for Product Service
│   │   │   └── handler.go            // Handlers for service endpoints
│   │   ├── interceptor/              // gRPC interceptors for Product Service
│   │   │   └── interceptor.go        // Interceptors for Product Service
│   │   ├── middleware/               // Middleware for Product Service
│   │   │   └── middleware.go         // Middleware functions
│   │   ├── repository/               // Repository layer for Product Service
│   │   │   └── repository.go         // Data access layer
│   │   └── proto/                    // gRPC proto definitions
│   │       ├── product.proto         // Proto file defining gRPC services
│   │       └── product.pb.go         // Generated Go code from proto
│   ├── migration/                    // Database migration scripts
│   │   └── 001_create_product_table.up.sql  // Example migration script
│   └── config/                       // Configuration files
│       └── config.yaml               // Configuration specific to Product Service
├── order-service/
│   ├── cmd/
│   │   └── main.go                   // Main entry point for Order Service
│   ├── internal/
│   │   ├── config/                   // Configuration package
│   │   │   └── config.go             // Configuration handling
│   │   ├── handler/                  // Request handlers for Order Service
│   │   │   └── handler.go            // Handlers for service endpoints
│   │   ├── interceptor/              // gRPC interceptors for Order Service
│   │   │   └── interceptor.go        // Interceptors for Order Service
│   │   ├── middleware/               // Middleware for Order Service
│   │   │   └── middleware.go         // Middleware functions
│   │   ├── repository/               // Repository layer for Order Service
│   │   │   └── repository.go         // Data access layer
│   │   └── proto/                    // gRPC proto definitions
│   │       ├── order.proto           // Proto file defining gRPC services
│   │       └── order.pb.go           // Generated Go code from proto
│   ├── migration/                    // Database migration scripts
│   │   └── 001_create_order_table.up.sql    // Example migration script
│   └── config/                       // Configuration files
│       └── config.yaml               // Configuration specific to Order Service
└── user-service/
    ├── cmd/
    │   └── main.go                   // Main entry point for User Service
    ├── internal/
    │   ├── config/                   // Configuration package
    │   │   └── config.go             // Configuration handling
    │   ├── handler/                  // Request handlers for User Service
    │   │   └── handler.go            // Handlers for service endpoints
    │   ├── interceptor/              // gRPC interceptors for User Service
    │   │   └── interceptor.go        // Interceptors for User Service
    │   ├── middleware/               // Middleware for User Service
    │   │   └── middleware.go         // Middleware functions
    │   ├── repository/               // Repository layer for User Service
    │   │   └── repository.go         // Data access layer
    │   └── proto/                    // gRPC proto definitions
    │       ├── user.proto            // Proto file defining gRPC services
    │       └── user.pb.go            // Generated Go code from proto
    ├── migration/                    // Database migration scripts
    │   └── 001_create_user_table.up.sql      // Example migration script
    └── config/                       // Configuration files
        └── config.yaml               // Configuration specific to User Service
