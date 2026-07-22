# Task Management API (Go + Echo)

This project implements the backend technical assessment for a multi-user task management API.

## Features
- Register and login with JWT authentication
- Task CRUD scoped by task owner
- List query support: status filter, title search, pagination
- Transactional task assignment with task log insert and rollback on failure
- Structured error responses
- Structured request logging with request_id and latency
- Unit tests without database/external dependencies

## Tech Stack
- Go
- Echo
- PostgreSQL

## Project Structure
- cmd/server: application entrypoint
- internal/config: environment configuration
- internal/auth: JWT token manager
- internal/usecase: application layer (business use cases and ports)
- internal/domain: entities and shared error models
- internal/repository: PostgreSQL queries and transactional operations
- internal/httpapi: handlers, middleware, and global error handler
- internal/notify: mock notifier for assignment flow
- migrations: database schema
- test/unit: unit tests

## Prerequisites
- Go 1.22 or newer
- PostgreSQL 14 or newer

## Environment Configuration
Copy .env.example to .env and adjust values.

Required variables:
- APP_PORT
- APP_ENV
- DB_HOST
- DB_PORT
- DB_NAME
- DB_USER
- DB_PASSWORD
- DB_SSLMODE
- JWT_SECRET
- JWT_EXPIRES_IN

## Running the Project
1. Start PostgreSQL (or run docker compose up -d).
2. Create database schema using migrations/001_init.sql.
3. Install module dependencies.
4. Start server.

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

## Endpoint Notes

GET /tasks query params:
- status: todo | doing | done
- search: case-insensitive title search
- limit: page size
- page: page number

POST /tasks/:id/assign transaction behavior:
- Lock task row
- Validate same-team rule between owner and assignee
- Update task assignee
- Insert task_logs row
- Trigger mock notification
- Rollback all changes if any step fails

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

## Testing
Run all tests:
- go test ./...

Run with race detector:
- go test -race ./...

## Brief Architecture Explanation
This code uses clean boundaries while staying practical for a coding assessment:
- HTTP layer in internal/httpapi handles routing, auth middleware, request parsing, and response formatting.
- Application layer in internal/usecase contains business rules orchestration (auth flow, transactional assignment flow).
- Domain layer in internal/domain defines core entities and structured API errors.
- Infrastructure layer includes internal/repository (PostgreSQL), internal/auth (JWT/password), and internal/notify (mock notifier).

This separation keeps handlers thin, transactions explicit, and race-sensitive behavior testable.

## Documentation
- Technical_Test_Implementation_Plan_Go_Echo.md: implementation plan document


## Tested Scenario Using Postman
- [x] Register User : Positive Case (Success Register) 
- [x] Login User : Positive Case (Success Login)
- [x] Create Task : Positive Case (Success Create Task)
- [x] List Task : Positive Case (Only Task Assigned To Current User And Created by Current User)
- [x] List Task : Filter By Title
- [x] List Task : Filter By Status
- [x] List Task : Filter By Limit and Page
- [x] Detail Task : Positive Case (Only Current User Assigned Task And Current User Owned Task)
- [x] Detail Task : Negative Case (Cannot See Another User Task)
- [x] Update Task : Positive Case (Success Update Current User Owned Task)
- [x] Update Task : Positive Case (Success Update Current User Assigned Task)
- [x] Update Task : Negative Case (Cannot Update Other User Owned or Assigned Task)
- [x] Delete Task : Positive Case (Success Delete Owned Task)
- [x] Delete Task : Negative Case (Cannot Delete Assigned Task)
- [x] Delete Task : Negative Case (Cannot Delete Other User Task)
- [x] Assign Task : Positive Case (Success Assign Task To Own Team)
- [x] Assign Task : Negative Case (Cannot Assign Other User Owned Task)
- [x] Assign Task : Negative Case (Cannot Assign Task To Another Team)