# HTL-Defense Server Documentation

## Overview
The HTL-Defense server is the central component of the HTL-Defense cybersecurity platform. It is responsible for handling threat analysis and providing REST/gRPC APIs for interaction with the frontend and other services.

## Directory Structure
- **api/**: Contains the REST/gRPC APIs for the server.
- **core/**: Holds the core detection logic, including signatures, heuristics, and behavior analysis.
- **db/**: Includes the database schema and models.
- **logs/**: Manages the logging system for monitoring and debugging.
- **tests/**: Contains unit tests for the server to ensure functionality and reliability.

## Getting Started
To get started with the HTL-Defense server, follow these steps:

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yourusername/HTL-Defense.git
   cd HTL-Defense/backend/server
   ```

2. **Install Dependencies**
   Ensure you have Go installed, then run:
   ```bash
   go mod tidy
   ```

3. **Configuration**
   Edit the `config.yaml` file to set up your server configuration.

4. **Run the Server**
   Start the server using:
   ```bash
   go run main.go
   ```

## API Documentation
Refer to the `api/README.md` file for detailed information on the available APIs, including endpoints, request/response formats, and authentication methods.

## Testing
To run the unit tests for the server, navigate to the `tests` directory and execute:
```bash
go test ./...
```

## Contribution
For guidelines on contributing to the HTL-Defense project, please refer to the `docs/contributing.md` file.

## License
This project is licensed under the MIT License. See the `LICENSE` file for more details.