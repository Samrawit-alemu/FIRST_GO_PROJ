Task Management API Documentation

1. Setup and Configuration

This API requires a running MongoDB instance to function. The application is configured to connect to a local MongoDB server by default.

Prerequisites

Go (version 1.18+ recommended)
MongoDB Server

Running the Database

You can run a local MongoDB instance using various methods. If you have MongoDB installed locally, ensure the service is running.

Running the API

a. Clone the repository or navigate to the project's root directory (task_manager/).
b. Install dependencies:
                        go mod tidy


c. Run the server:
                    go run main.go

The server will start on http://localhost:8080.

2. API Endpoints

Base URL: http://localhost:8080

Get All Tasks:

    Endpoint: GET /tasks
    Description: Retrieves a list of all tasks.
    Success Response (200 OK):
Json:
[
  {
    "id": "655a8b5b8f673a3a4b5e6e6a",
    "title": "Learn Go",
    "description": "Complete the Go tutorial.",
    "due_date": "2024-12-01T00:00:00Z",
    "status": "Pending"
  }
]

Get a Specific Task:

    Endpoint: GET /tasks/:id
    Description: Retrieves a single task by its unique ID.
    URL Parameter: id (string, required) - The ObjectID of the task to retrieve.
    Success Response (200 OK):
Json
{
  "id": "655a8b5b8f673a3a4b5e6e6a",
  "title": "Learn Go",
  "description": "Complete the Go tutorial.",
  "due_date": "2024-12-01T00:00:00Z",
  "status": "Pending"
}

    Error Responses:
        404 Not Found: {"error":"task not found"}
        400 Bad Request: {"error":"invalid task ID format"}

Create a New Task

    Endpoint: POST /tasks
    Description: Creates a new task.
    Request Body:
Json
{
    "title": "Build API",
    "description": "Develop a REST API using Go and Gin.",
    "due_date": "2024-10-25T15:00:00Z",
    "status": "In Progress"
}

    Success Response (201 Created): Returns the newly created task object with its generated id.
Json
{
  "id": "655a8c1f3b3a4e5e6f7g8h9i",
  "title": "Build API",
  "description": "Develop a REST API using Go and Gin.",
  "due_date": "2024-10-25T15:00:00Z",
  "status": "In Progress"
}

    Error Response (400 Bad Request): {"error":"Invalid request data: ..."}

Update a Task

    Endpoint: PUT /tasks/:id
    Description: Updates an existing task's details.
    URL Parameter: id (string, required) - The ObjectID of the task to update.
    Request Body:
Json
{
    "title": "Build an Awesome API",
    "description": "Develop a REST API using Go and Gin with documentation.",
    "due_date": "2024-10-26T18:00:00Z",
    "status": "Completed"
}

    Success Response (200 OK): Returns the full, updated task object.
    Error Responses:
        404 Not Found: {"error":"task not found"}
        400 Bad Request: For invalid ID format or invalid request data.

Delete a Task

    Endpoint: DELETE /tasks/:id
    Description: Deletes a task by its ID.
    URL Parameter: id (string, required) - The ObjectID of the task to delete.
    Success Response: 204 No Content (The response will have no body).
    Error Responses:
        404 Not Found: {"error":"task not found"}
        400 Bad Request: {"error":"invalid task ID format"}