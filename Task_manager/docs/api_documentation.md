# Task Management API Documentation

1. Get all tasks

    Endpoint : "GET/tasks"
    Description: Retrieves a list of all tasks.
    Success Response:
        Code: 200 OK
        Content: [{"id":1,"title":"Learn Go","description":"Complete the Go tutorial.","due_date":"2024-12-01T00:00:00Z","status":"Pending"}]

2. Get a specific task

    Endpoint : "GET/tasks/:id"
    Description: Retrieves a single task by its unique ID.
    Success Response:
       Code: 200 OK
       Content: {"id":1,"title":"Learn Go","description":"Complete the Go tutorial.","due_date":"2024-12-01T00:00:00Z","status":"Pending"}
    Error Response:
        Code: 404 Not Found
        Content: {"error":"task not found"}


3. Create a new task

    Endpoint : "POST/tasks"
    Description : Creates a new task.
    Success Response:
       Code : 201 Created
       Content: {"id":2,"title":"Build API","description":"...","due_date":"...","status":"In Progress"}
    Error Response :
        Code: 400 Bad Request`
        Content: {"error":"Invalid request data"}

4. Update a task

    Endpoint : "PUT/tasks/:id"
    Description : Updates an existing task's details.
    Success Response :
        Code : 200 OK
        Content: (The full, updated task object)
    Error Response:
        Code : 404 Not Found
        Content : {"error":"task not found"}

5. Delete a task
    
    Endpoint : "DELETE/tasks/:id"
    Description : Deletes a task by its ID.
    Success Response:
        Code: 204 No Content
    Error Response:
        Code : 404 Not Found
        Content : {"error":"task not found"}