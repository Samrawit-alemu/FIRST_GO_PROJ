Task Management API (Clean Architecture)

This project is a RESTful API for managing tasks, built with Go and the Gin framework. It is designed following Clean Architecture principles to ensure a clear separation of concerns, high testability, and maintainability.

The architecture is divided into distinct layers: Domain, Usecases, Repositories, Infrastructure, and Delivery, with a strict dependency rule pointing inwards. Data Transfer Objects (DTOs) are used at the API boundary to decouple the internal domain models from the external world.

Features

User Management: Secure user registration and login.
JWT Authentication: Protected endpoints using JSON Web Tokens.
Role-Based Access Control (RBAC): Distinction between admin and user roles.

The first user to register automatically becomes an admin.
Admins can create, update, and delete tasks.
All authenticated users can view their own tasks.
Task Management: Full CRUD (Create, Read, Update, Delete) operations for tasks, respecting user ownership.
Persistent Storage: Uses MongoDB for data persistence.

Architectural Layers

Domain: Contains the pure, core business entities (User, Task) with no external dependencies or tags.
Usecases: Holds the application-specific business rules and orchestrates the flow between other layers. Depends only on repository interfaces.
Repositories: Defines the data access contracts (interfaces) and provides the concrete MongoDB implementations. Handles mapping between domain models and database models.
Infrastructure: Implements low-level, external concerns like password hashing (bcrypt), JWT services, and middleware.
Delivery: The entrypoint to the application, responsible for handling HTTP requests. Contains the Gin router, controllers, and Data Transfer Objects (DTOs).


1. Setup and Configuration

Prerequisites

Go (version 1.21+ recommended)
MongoDB Server (local or Docker)

A .env file in the project root.
Database Setup
This API requires a running MongoDB instance.
Connection String: The application connects to mongodb://localhost:27017 by default.
Using Docker :
docker run --name my-mongo -p 27017:27017 -d mongo

Environment Variables
Create a .env file in the project's root directory and add the following key:
JWT_SECRET=a_super_secret_key_that_is_long_and_random

Important: Add .env to your .gitignore file to prevent committing secrets.
Running the API
Navigate to the project's root directory.
Install dependencies:
go mod tidy

Run the server:
go run delivery/main.go

The server will start on http://localhost:8080.


2. API Endpoints

Base URL: http://localhost:8080
Authentication Endpoints (Public)
Register a New User
Endpoint: POST /auth/register
Description: Creates a new user. The first user becomes an admin.
Request Body (dto.RegisterRequest):
{
    "username": "someuser",
    "password": "a_strong_password"
}

Success Response (201 Created):
{
    "message": "User registered successfully",
    "user": {
        "id": "655a8c1f...",
        "username": "someuser",
        "role": "user"
    }
}


Login

Endpoint: POST /auth/login
Description: Authenticates a user and returns a JWT.
Request Body (dto.LoginRequest):
{
    "username": "someuser",
    "password": "a_strong_password"
}

Success Response (200 OK):
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ..."
}

Protected Task Endpoints

All endpoints below require a valid JWT in the format Authorization: Bearer <token>.
Get All Tasks for the Logged-in User

Endpoint: GET /tasks
Authorization: user or admin.
Description: Retrieves a list of tasks created by the authenticated user.
Success Response (200 OK, []dto.TaskResponse): An array of task objects.

Create a New Task

Endpoint: POST /tasks
Authorization: admin role only.
Request Body (dto.TaskRequest):

{
    "title": "A New Admin Task",
    "description": "Details about the task.",
    "due_date": "2025-10-25T15:00:00Z",
    "status": "Pending"
}
Success Response (201 Created, dto.TaskResponse): The newly created task object.
Error Response (403 Forbidden): If a non-admin user attempts this action.
(... and so on for GET by ID, PUT, and DELETE task endpoints, explaining their authorization rules)

Protected Admin Endpoints

Promote a User to Admin
Endpoint: PUT /admin/promote/:id
Authorization: admin role only.
Description: Promotes a regular user to have the admin role.
Success Response (200 OK):
{
    "message": "User promoted successfully",
    "user": {
        "id": "...",
        "username": "promoteduser",
        "role": "admin"
    }
}

3. Testing
This project includes a comprehensive suite of unit and integration tests.

Prerequisites

A running MongoDB instance is required for the repository integration tests.
mockery must be installed for regenerating mocks:
go install github.com/vektra/mockery/v2@latest
Running Tests
To run all tests for the project, navigate to the root directory and execute:
go test ./... -v

Checking Test Coverage

To run the tests and generate a coverage report, use the -cover flag:
                                                        go test ./... -cover

To generate a visual HTML report of which lines are covered:
                            go test ./... -coverprofile=coverage.out
                            go tool cover -html=coverage.out