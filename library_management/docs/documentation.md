# Library Management System

This is a simple console-based library management system written in Go. It demonstrates the use of structs, interfaces, methods, maps, and slices, organized into a clear, multi-file structure.

## Project Structure

-   `/main.go`: The entry point of the application. It initializes the services and controllers and runs the main menu loop.
-   `/models`: Contains the data structures for the application (`Book`, `Member`).
-   `/services`: Contains the business logic. It defines the `LibraryManager` interface and provides its concrete implementation in the `Library` struct.
-   `/controllers`: Acts as the bridge between the user interface (console) and the service layer. It handles user input and output.
-   `/docs`: Contains project documentation.
-   `/go.mod`: Defines the Go module and its dependencies.

## How to Run

1.  Navigate to the root directory of the project (`library_management/`).
2.  Run the application using the following command:
    ```sh
    go run main.go
    ```
3.  Follow the on-screen menu to interact with the library system.