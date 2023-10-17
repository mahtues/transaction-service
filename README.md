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

### Transaction Service

### Transaction Repository

### Conversion Service

## Requests

## Future Improvements
