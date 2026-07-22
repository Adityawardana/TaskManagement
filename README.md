# Task Management API (Go + Echo)

This project implements the backend technical assessment for a multi-user task management API.

## Tech Stack
- Go
- Echo
- PostgreSQL

## Prerequisites
- Go 1.22 or newer
- PostgreSQL 14 or newer

## How to running the Project
1. Start PostgreSQL (or run docker compose up -d).
2. Create database schema using migrations/001_init.sql.
3. Run database seeder query for dummy data using migrations/dummy/seeder.sql(optional).
4. Install module dependencies.
5. Copy .env.example to .env and adjust values.
6. Start server.
7. Open Postman and Import Task Management.postman_collection.json to test API.

Commands:
- go mod tidy
- go run main.go

Server address:
- http://localhost:8080

## Endpoints

Authentication:
- POST /auth/register
- POST /auth/login

Tasks:
- POST /tasks
- GET /tasks
- GET /tasks/:id
- PUT /tasks/:id
- DELETE /tasks/:id
- POST /tasks/:id/assign

## Architecture
This project uses a modular Clean Architecture, where each feature (such as Auth and Task) is organized into its own module. Every module has its own domain, use case, repository, and delivery layers, while shared components are reused across the application.

- Delivery Layer (internal/modules/*/delivery/resthandler) handles HTTP requests, routing, request validation, and responses.
- Use Case Layer (internal/modules/*/usecase) contains the application's business logic.
- Domain Layer (internal/modules/*/domain) defines the core entities and business models for each module.
- Repository Layer (internal/modules/*/repository) handles database operations using PostgreSQL.
- Module Initialization (internal/service.go) creates shared dependencies (database, JWT, password hasher, notifier, configuration), injects them into each module, and starts the server.
- Shared Packages (pkg/shared and pkg/helper) contain common utilities such as middleware, error handling, JWT, password hashing, logging, and the mock notifier.

This structure keeps the code organized, makes each module easier to maintain and test, reduces code duplication, and clearly separates business logic from HTTP and database code.


------------------------------------------------------------Additional Information------------------------------------------

## Features
- Register and login with JWT authentication
- Task CRUD scoped by task owner, with read/update access also granted to the assigned user
- List query support: status filter, title search, pagination
- Transactional task assignment with task log insert and rollback on failure
- Structured error responses
- Structured request logging with request_id and latency

## Error Response Format
All API errors use this shape:
- status
- code
- message
- timestamp

## Logging and Observability
Every request emits structured JSON with:
- request_id
- method
- path
- status_code
- latency_ms

Log levels:
- INFO for 2xx
- WARN for 4xx
- ERROR for 5xx