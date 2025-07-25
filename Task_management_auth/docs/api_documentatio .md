 Task Management API Documentation

This API provides a secure, role-based system for managing tasks. It includes user registration, login, and protected endpoints for task management.

1. Setup and Configuration
Prerequisites

- Go (version 1.18+ recommended)
- MongoDB Servers
- A .env file in the project root.

Database Setup

This API requires a running MongoDB instance to function.

1. Install MongoDB: Download and install MongoDB Community Server from the official MongoDB website. Follow the installation instructions for your operating system.
2. Ensure MongoDB is Running: Make sure the MongoDB service (mongod) is running on your machine. You can typically check this in your system's "Services" panel (Windows) or by using brew services list (macOS with Homebrew).
3. Connection String: The application is configured to connect to mongodb://localhost:27017 by default, which is the standard for a local installation. This can be changed in data/db.go.

Environment Variables

Create a .env file in the project's root directory and add the following key. This key is used to sign the authentication tokens.

JWT_SECRET=a_super_secret_key_that_is_long_and_random

Important: Add .env to your .gitignore file to prevent committing secrets.

Running the API

1. Navigate to the project's root directory (task_manager/).
2. Install dependencies:
                    go mod tidy

Run the server:
                go run main.go

The server will start on http://localhost:8080.


2. Authentication Endpoints

These endpoints are public and do not require an authentication token.

Register a New User

- Endpoint: POST /auth/register
- Description: Creates a new user account. The first user to register will automatically be assigned the admin role. All subsequent users will be assigned the user role.
- Request Body:

{
    "username": "someuser",
    "password": "a_strong_password"
}

- Success Response (201 Created):

{
    "message": "User registered successfully",
    "user_id": "655a8c1f3b3a4e5e6f7g8h9i"
}

- Error Response (400 Bad Request): For invalid data or if the username already exists.

Login

- Endpoint: POST /auth/login
- Description: Authenticates a user and returns a JSON Web Token (JWT) for accessing protected routes.
- Request Body:

{
    "username": "someuser",
    "password": "a_strong_password"
}

- Success Response (200 OK):

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAi..."
}

- Error Response (401 Unauthorized): For invalid credentials.


3. Protected Task Endpoints

All endpoints below require a valid JWT to be sent in the Authorization header.

Format: Authorization: Bearer <your_jwt_token>

Get All Tasks for the Logged-in User

- Endpoint: GET /tasks
- Authorization: user or admin role.
- Description: Retrieves a list of all tasks created by the authenticated user.
- Success Response (200 OK): An array of task objects.

Get a Specific Task

- Endpoint: GET /tasks/:id
- Authorization: user or admin role.
- Description: Retrieves a single task by its ID, but only if it was created by the authenticated user.
- Success Response (200 OK): A single task object.
- Error Response (404 Not Found): If the task does not exist or does not belong to the user.

Create a New Task

- Endpoint: POST /tasks
- Authorization: admin role only.
- Description: Creates a new task and assigns it to the authenticated admin user.
- Request Body:

{
    "title": "A New Admin Task",
    "description": "Details about the task.",
    "due_date": "2025-10-25T15:00:00Z",
    "status": "Pending"
}

- Success Response (201 Created): The newly created task object.
- Error Response (403 Forbidden): If a non-admin user attempts this action.
Update a Task

- Endpoint: PUT /tasks/:id
- Authorization: admin role only.
- Description: Updates an existing task. An admin can only update tasks they created.
- Success Response (200 OK): The full, updated task object.
- Error Response (403 Forbidden): If a non-admin user 
attempts this action.

Delete a Task

- Endpoint: DELETE /tasks/:id
- Authorization: admin role only.
- Description: Deletes a task. An admin can only delete tasks they created.
- Success Response: 204 No Content.
- Error Response (403 Forbidden): If a non-admin user attempts this action.


4. Protected Admin Endpoints

These endpoints are for administrative purposes and require an admin role.

Promote a User to Admin

- Endpoint: PUT /admin/promote/:id
- Authorization: admin role only.
- Description: Promotes a regular user to have the admin role.
- URL Parameter: id (string, required) - The ObjectID of the user to promote.
- Success Response (200 OK):

{
    "message": "User promoted to admin successfully",
    "user": {
        "id": "...",
        "username": "promoteduser",
        "role": "admin"
    }
}

- Error Response (404 Not Found): If the user ID does not exist.