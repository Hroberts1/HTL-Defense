# Installation Guide for HTL-Defense

## Prerequisites
Before you begin, ensure you have the following installed on your system:
- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 1.27 or higher
- **Go**: Version 1.16 or higher (for backend development)
- **Node.js**: Version 14 or higher (for frontend development)

## Installation Steps

1. **Clone the Repository**
   Open your terminal and run the following command to clone the repository:
   ```
   git clone https://github.com/yourusername/HTL-Defense.git
   cd HTL-Defense
   ```

2. **Set Up Environment Variables**
   Copy the example environment variables file and modify it as needed:
   ```
   cp .env.example .env
   ```

3. **Build the Backend**
   Navigate to the backend directory and build the server:
   ```
   cd backend/server
   go build -o server main.go
   ```

4. **Build the Frontend**
   Navigate to the frontend directory and install dependencies:
   ```
   cd ../../frontend
   npm install
   ```

5. **Run the Application with Docker Compose**
   Go back to the root directory of the project and run:
   ```
   docker-compose up --build
   ```

6. **Access the Application**
   Once the application is running, you can access the web dashboard at `http://localhost:3000`.

## Additional Configuration
- For advanced configurations, refer to the `config.yaml` file located in the `backend/server` directory.
- You can customize the Docker setup by modifying the `docker-compose.yml` file.

## Troubleshooting
- If you encounter issues, check the logs in the `backend/server/logs` directory for more information.
- Ensure all dependencies are correctly installed and up to date.

## Conclusion
You are now ready to use the HTL-Defense cybersecurity platform. For further information, refer to the documentation in the `docs` directory.