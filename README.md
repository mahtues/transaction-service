# Transaction Service

## Introduction

The project is implemented in Golang, using only the standard libraries. Not using third-party libraries was a choice. This approach is more challenging in the sense that a framework's basic features are not available and must be implemented or at least taken into consideration. It also forces an exploration and deep understanding of the technology.

The project structure was inspired by the Hexagonal Architecture. It comprises a core of business logic independent of technology and infrastructure. And attachable driver and driven actors.

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

Otherwise, there is a status code error and a JSON with an error message. e.g.

```json
{
  "error": "transaction id tr-frnshjdacxdnjhsczqq not found"
}
```

### Transaction Service

Transaction Service contains the business logic for transactions in general. It routes creation requests to the repository and uses the conversion service to retrieve requests.

### Transaction Repository

The only responsibility of the Transaction Repository is to access the persistence. In the current implementation, it is an in-memory key-value database (a Hash Map). Other persistence implementations can easily replace it. For example, the dynamo repository is functional and can be used if properly configured.

### Conversion Service

Conversion Service consumes the Treasury Reporting Rates of Exchange API.

## Tests

A non-existent set of unit tests are implemented. In the project's root, running

```
$ go test -v ./...
```

Unit tests are implemented, creating the service to be tested and injecting mocks as dependencies. Those mocks can check for method calls and return distinct values to increase the test coverage.

## Running

Running the server can be done by executing from the project's root

```
$ go run app/avocado/main.go
```

The server runs on port `8000`. Some `curl` examples:

```bash
$ curl -s http://localhost:8000/transaction -d 'description=desc&amountUs=123.45&date=2022-07-31T20:04:59Z' | jq
{
  "id": "tr-iympvyxnfwqvvuprkzku"
}
$ curl -s http://localhost:8000/transaction/tr-iympvyxnfwqvvuprkzku?country=Brazil | jq
{
  "id": "tr-iympvyxnfwqvvuprkzku",
  "description": "message",
  "date": "2022-07-31T20:04:59Z",
  "amountUs": "123.45",
  "conversionRate": "5.182",
  "amountConverted": "639.72"
}
```

## Github Workflow

Github Actions is used to run tests in PRs.

## Future Improvements

Some requirements were not implemented due to time constraints. This is a small project I will keep developing as a way to increase my Golang knowledge.

The current implementation uses `float64` to convert currencies. This is **WRONG**. Native floating point types must be avoided when decimal precision is required, such as monetary representation. Research for a decimal library will be conducted.

Upgrading my account to a paid tier, I would use Github actions to

- block commits to main;
- merges to main only from PRs with at least some amount of approvers and passing all actions (just unit tests for now);
- add integration tests. It is possible to spin up some infrastructure using docker-compose;
- create staging and development environments and a CD to deliver code to those places.
