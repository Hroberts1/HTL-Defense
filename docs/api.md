# API Documentation for HTL-Defense

## Overview
The HTL-Defense platform provides a set of APIs that allow interaction with the backend services. This documentation outlines the available endpoints, their methods, and the expected request and response formats.

## Base URL
```
http://<your-server-address>/api
```

## Authentication
All API requests require authentication. Use the following method to authenticate:

- **Token-Based Authentication**: Include the token in the `Authorization` header as follows:
```
Authorization: Bearer <your-token>
```

## Endpoints

### 1. Get Threat Intelligence
- **Endpoint**: `/threats`
- **Method**: `GET`
- **Description**: Retrieves the latest threat intelligence data.
- **Response**:
  - **200 OK**: Returns a list of threats.
  - **401 Unauthorized**: If authentication fails.

### 2. Submit Incident Report
- **Endpoint**: `/incidents`
- **Method**: `POST`
- **Description**: Submits a new incident report.
- **Request Body**:
```json
{
  "title": "Incident Title",
  "description": "Detailed description of the incident",
  "severity": "high"
}
```
- **Response**:
  - **201 Created**: Returns the created incident details.
  - **400 Bad Request**: If the request body is invalid.

### 3. Get Incident Status
- **Endpoint**: `/incidents/{id}`
- **Method**: `GET`
- **Description**: Retrieves the status of a specific incident by ID.
- **Response**:
  - **200 OK**: Returns the incident details.
  - **404 Not Found**: If the incident ID does not exist.

## Error Handling
All error responses will include a standard error format:
```json
{
  "error": {
    "code": "error_code",
    "message": "Error message describing the issue"
  }
}
```

## Rate Limiting
To ensure fair usage, the API enforces rate limiting. Exceeding the limit will result in a `429 Too Many Requests` response.

## Conclusion
This API documentation provides a high-level overview of the available endpoints and their usage. For further details, please refer to the specific endpoint documentation or contact the development team.