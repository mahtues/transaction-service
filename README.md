# Transaction Service

## Introduction

## Architecture

![image](https://github.com/mahtues/transaction-service/assets/14203456/a058a367-ab3f-4d10-bce8-2d55257b8503)

```
transaction-service$ tree
.
├── app
│   └── avocado
│       ├── main.go
│       └── server
│           ├── misc.go
│           └── server.go
├── apperrors
│   ├── errors.go
│   └── errors_test.go
├── conversion
│   ├── conversion.go
│   └── conversion_test.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── random
│   └── random.go
├── README.md
├── support
│   ├── aws.go
│   ├── errors.go
│   └── time.go
└── transaction
    ├── repo_dynamo.go
    ├── repo_dynamo_test.go
    ├── repo.go
    ├── repo_inmem.go
    ├── service.go
    └── service_test.go

```

### Core

The Core type - a.k.a. `Avocado` - glues the entire together. It is responsible for instantiating all application's services, repositories, and resources. And also injecting the dependencies in the modules. This is the place the application grows from.

### Server

The server has two end-points.

```
POST /transaction
Content-Type: application/x-www-form-urlencoded

description=<description>&amountUs=123.45&date=2006-01-02T15:04:05Z
```

Form's fields format:

```
description: "string 1-50 characters"
amountUs: "123.45"
date: "2006-01-02T15:04:05Z"
```

In case of success, return `200 OK` and a JSON with the transaction ID created:

```json
{
  "id": "tr-frnshjdacxdnjhsczqqe"
}
```

Otherwise, a status code error and a JSON with an error message. e.g.

```json
{
  "error": "invalid form: description missing, date missing, amountUs missing"
}
```

And

```
GET /transaction/<transaction-id>?country=<country>
```

In case of success, return `200 OK` and a JSON with the transaction ID created:

```
{
  "id": "tr-frnshjdacxdnjhsczqqe",
  "description": "asd",
  "date": "2022-07-31T20:04:59Z",
  "amountUs": "123.3",
  "conversionRate": "5.182",
  "amountConverted": "638.94"
}
```

Otherwise, a status code error and a JSON with an error message. e.g.

```json
{
  "error": "transaction id tr-frnshjdacxdnjhsczqq not found"
}
```

### Transaction Service

### Transaction Repository

### Conversion Service

## Requests

## Future Improvements
