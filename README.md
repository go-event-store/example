# GO EventStore Example

## Example Project based on [Gin](https://github.com/gin-gonic/gin)

One more overengineered ToDo API to demonstrate the [GO EventStore](https://github.com/go-event-store/eventstore) Package

## Required

- Go >= 1.13
- Postgres

## Features

- Swagger UI to try the ToDo API
- `docker-compose.yaml` to start the example as Docker Container
- `Makefile` to update SwaggerUI, start or build the App

## Exampled Features

- Define and register Aggregates and Events
- Create and use an Repository to load and save Aggregates
- Create and prepare an EventStore
- Create an EventStream
- Build ReadModel Projections
- Use persisted ReadModel with an custom ReadModelFinder
- Execute an ReadModel in the Background
- Use Gin to perform your actions in form of Queries and Commands

## Start the Project

Use `docker-compose.yaml` to start only the required Postgres Database

```bash
docker-compose up -d postgres
```

Use the Makefile to start the App locally

```bash
make start
```

You could also start the App as Docker Container

```bash
docker-compose up -d app
```

## Try it out

Serve the SwaggerUI under

```
http://localhost:8080/swagger/index.html
```