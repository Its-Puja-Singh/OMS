# Order Management Service

This is a simple order management service that allows you to add and update orders, and fetch them in a sorted and filtered way. The service is built using Go, MySQL DB and JSON-HTTP REST APIs.

## Prerequisites

Before you can run this service, you need to have the following installed on your system:

```
Go
MySQL
```
## Getting Started

To get started with running the service, follow these steps:


## Create A MySQL Database:

```
CREATE DATABASE order_management_service;
```

## Create Table orders the same databse

```
CREATE TABLE orders (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    items JSON NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    currency_unit VARCHAR(3) NOT NULL
  );
```

## Set The Database Credentials:

```
db, err = sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/order_management_service")
```

## Run The Service:

```
go mod init oms

go get

go run .
```

## Test The Endpoint:

You can use a tool like curl or Postman to test the service. Here are some sample requests:

## Add An Order:

```
curl --location 'http://localhost:8000/createOrder' \
--header 'Content-Type: application/json' \
--data '{
    "id": "abcdef-123457",
    "status": "PENDING_INVOICE",
    "items": [{
    "id": "123457 ",
    "description": "a product description",
    "price": 12.40,
    "quantity": 1
    }],
    "total": 12.40,
    "currencyUnit": "USD"
}'
```

## Get Orders:

```
curl --location 'http://localhost:8000/getOrdersById'
```

## Get Order By Params Status Or Currency Unit
```
curl --location 'http://localhost:8000/getOrdersByParams?status=Delivered'
```

## Update Order

```
curl --location --request PUT 'http://localhost:8000/updateOrder' \
--header 'Content-Type: application/json' \
--data '{
    "Status":"Not Delivered",
    "ID":"abcdef-123457"
}'
```

## To Run The Test Cases

```
go test -v
```