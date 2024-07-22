# Maglo: Financial Banking System

Maglo is a financial banking system that allows users to perform basic banking operations such as creating accounts, logging in, managing user information, and handling transactions. The system uses Google OAuth for authentication and integrates with PostgreSQL using GORM for data storage.

## Features

- User Authentication (Google OAuth)
- User Management
- Transaction Management

## Setup

### Prerequisites

- Go 1.16 or later
- PostgreSQL

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/alade-dev/maglo.git
    cd maglo
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

3. Configure environment variables:

    Create a `.env` file and add the following environment variables:

    ```
    Client=<your_google_client_id>
    Secret=<your_google_client_secret>
    redirect_url=http://localhost:3000/api/auth/google/callback # note change the port to your own respective port
    DB_HOST=
    DB_PORT=
    DB_USER=
    DB_PASSWORD=
    DB_NAME=
    Secret= #jwt-secret
    PORT=
    
    ```

4. Run the application:

    ```sh
    go run main.go
    ```

## API Endpoints

### Base URL

```
http://localhost:8000/
```

### Routes

#### Authentication

- **Login**

  ```http
  POST /api/auth/login
  ```

- **Google OAuth**

  ```http
  GET /api/auth/google
  ```

- **Google OAuth Callback**

  ```http
  GET /api/auth/google/callback/
  ```

- **Logout**

  ```http
  POST /api/auth/logout
  ```

#### User Management

- **Get User by ID**

  ```http
  GET /api/user/:id
  ```

- **Create User**

  ```http
  POST /api/user
  ```

- **Update User**

  ```http
  PATCH /api/user/:id
  ```

- **Delete User**

  ```http
  DELETE /api/user/:id
  ```

#### Transaction Management

- **Get All Transactions**

  ```http
  GET /api/transaction
  ```

- **Get Transaction by ID**

  ```http
  GET /api/transaction/:id
  ```

- **Create Transaction**

  ```http
  POST /api/transaction
  ```

- **Delete Transaction**

  ```http
  DELETE /api/transaction/:id
  ```

## Middleware

- **Logger**: Logs request details.
- **Protected**: Ensures the route is accessible only by authenticated users.

## Project Structure

```
maglo/
├── auth
│   └── auth.go
├── config
│   └── database.go
├── handler
│   └── handler.go
├── middleware
│   └── middleware.go
├── model
│   └── model.go
├── router
│   └── router.go
├── main.go
├── go.mod
└── go.sum
```

## Contribution

Contact: ibisomimayowa@gmail.com