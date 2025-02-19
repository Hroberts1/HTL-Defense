# HTL-Defense

HTL-Defense is a comprehensive security solution designed to protect systems from various threats. This project includes both agent implementations for different operating systems and a backend server that manages logs and provides a REST API for interaction.

## Project Structure

```
HTL-Defense
├── agent
│   ├── windows
│   │   └── agent.go          # Windows-specific agent implementation
│   ├── linux
│       └── agent.go          # Linux-specific agent implementation
├── server
│   ├── api
│   │   ├── server.go         # Entry point for the REST API
│   │   ├── handlers.go       # HTTP handlers for API endpoints
│   ├── logs                  # Directory for JSON-based logs
│   ├── database
│   │   └── log_manager.go    # Log management and database interactions
│   └── main.go               # Main entry point for the backend application
├── sandbox
│   └── qemu                  # QEMU sandboxing system files
├── docs                      # Documentation for the project
├── config
│   └── config.yaml           # Configuration settings for the application
├── .gitignore                # Files and directories to ignore in version control
└── README.md                 # Project documentation
```

## Getting Started

To get started with HTL-Defense, follow these steps:

1. **Clone the repository**:
   ```
   git clone https://github.com/yourusername/HTL-Defense.git
   cd HTL-Defense
   ```

2. **Install dependencies**:
   Make sure you have Go installed. Then, navigate to the `server` directory and run:
   ```
   go mod tidy
   ```

3. **Configure the application**:
   Edit the `config/config.yaml` file to set your desired configuration options.

4. **Run the server**:
   From the `server` directory, execute:
   ```
   go run main.go
   ```

5. **Start the agent**:
   Depending on your operating system, navigate to either the `agent/windows` or `agent/linux` directory and run:
   ```
   go run agent.go
   ```

## Usage

Once the server is running, you can interact with the API using tools like `curl` or Postman. Refer to the documentation in the `docs` directory for detailed API usage instructions.

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.