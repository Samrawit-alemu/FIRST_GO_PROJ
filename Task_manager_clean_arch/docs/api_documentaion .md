  Testing

This project includes a comprehensive suite of unit and integration tests to ensure code quality and correctness.

- Prerequisites

- A running MongoDB instance is required for the repository integration tests.
- `mockery` must be installed for regenerating mocks:
  `go install github.com/vektra/mockery/v2@latest`

- Running Tests

To run all tests for the project, navigate to the root directory and execute the following command:
                                    go test ./... -v

Checking Test Coverage

To run the tests and generate a coverage report for each package, use the -cover flag:
            go test ./... -cover

To generate a visual HTML report of which lines are covered:
                        go test ./... -coverprofile=coverage.out

TO View the report in your browser:
                        go tool cover -html=coverage.out