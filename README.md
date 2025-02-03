# Naboobase

Naboobase is an innovative Backend-as-a-Service (BaaS) platform inspired by [Pocketbase](https://github.com/pocketbase/pocketbase). It is designed to be scalable and easy to maintain while providing powerful tools for rapid backend development in Go.

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Project Architecture](#project-architecture)
- [Installation and Setup](#installation-and-setup)
- [Usage](#usage)
  - [CLI: JSON to Go Struct Conversion](#cli-json-to-go-struct-conversion)
  - [Example Server](#example-server)
- [Code Structure](#code-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Naboobase aims to simplify backend development by providing a robust, modular framework built with Go. It includes utilities for:

- Converting JSON models to Go structs
- Managing API endpoints and routing with the Gin framework
- Integrating authentication and database connectors (MongoDB is used in the example)
- A clean, organized project structure for scalable development

---

## Key Features

- **JSON to Go Struct Conversion**: Quickly generate Go model structs from JSON definitions.
- **Modular Architecture**: Separates concerns (e.g., controllers, core functionality, utilities) to ease maintenance and scalability.
- **API Server Example**: Includes an example server (`example/server.go`) demonstrating endpoint creation, health checks, and database integration.
- **Built on Gin**: Leverages the Gin framework for fast and efficient HTTP server operations.
- **Database Integration**: Includes a core MongoDB connector for database operations.

---

## Project Architecture

The repository is organized as follows:

```
.
├── cli
│   └── main.go            # CLI tool to convert JSON to Go structs
├── configs
│   └── env.go             # Environment configurations
├── constants              # Project constants
├── controllers
│   └── user_controller.go # User controller for handling user-related endpoints
├── core
│   ├── authentication.go  # Authentication layer integration
│   ├── db.go              # Database connection logic
│   ├── route.go           # Routing utilities
│   └── server.go          # Server configuration and startup
├── example
│   └── server.go          # Example of how to use Naboobase to create a server
├── go.mod                 # Module definitions
├── go.sum                 # Module checksums
├── json
│   ├── user.json          # JSON definition for the User model
│   ├── user_request.json  # Example JSON request payload for user endpoints
│   └── user_response.json # Example JSON response payload for user endpoints
├── main.go                # Entry point (could be used for additional bootstrapping)
├── models
│   ├── error_response.go  # Error response model
│   ├── login_request.go   # Login request model
│   ├── login_response.go  # Login response model
│   ├── refreshToken_request.go # Refresh token request model
│   ├── user.go            # User model (includes BSON/JSON tags and validations)
│   ├── user_request.go    # User request model
│   └── user_response.go   # User response model
├── README.md              # Project documentation
├── routes                 # Route definitions for endpoints
└── utils
    ├── auth_utils.go      # Authentication utility functions
    ├── getStruct.go       # Utility for extracting struct information
    ├── jsonToGo.go        # Logic for converting JSON definitions to Go structs
    ├── struct_utils.go    # General utilities for struct handling
    └── user_utils.go      # Utility functions for user model operations
```

---

## Installation and Setup

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/naboobase.git
   cd naboobase
   ```

2. **Install Dependencies**

   Ensure you have [Go installed](https://golang.org/doc/install) (version 1.23.4 or higher is recommended). Then, run:

   ```bash
   go mod download
   ```

3. **Configure the Environment**

   Customize any configuration parameters in `configs/env.go` as needed (for example, setting database URIs, server host/port, etc.).

---

## Usage

### CLI: JSON to Go Struct Conversion

The CLI tool converts JSON model definitions into Go struct code. For example, given a JSON like:

```json
{
  "_id": {
    "value": "cjhjvzivfvbsif",
    "db": "autogenerate"
  },
  "name": {
    "value": "John Doe"
  },
  "email": {
    "value": "johndoe@example.com",
    "db": "unique"
  },
  "passwordHashed": {
    "value": "5f4dcc3b5aa765d61d8327deb882cf99"
  },
  "is_verified": {
    "value": false
  },
  "is_superuser": {
    "value": false
  },
  "role": {
    "value": "user"
  },
  "created_at": {
    "value": "2024-12-23T00:00:00Z"
  },
  "updated_at": {
    "value": "2024-12-23T12:00:00Z"
  }
}
```

Running the following command:

```bash
go run cli/main.go generate
```

will generate a Go struct similar to:

```go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    PasswordHashed string             `json:"passwordHashed" bson:"passwordHashed" validate:"max=255"`
    IsVerified     bool               `json:"is_verified" bson:"is_verified"`
    IsSuperuser    bool               `json:"is_superuser" bson:"is_superuser"`
    Role           string             `json:"role" bson:"role" validate:"max=255"`
    CreatedAt      string             `json:"created_at" bson:"created_at" validate:"max=255"`
    ID             primitive.ObjectID `json:"_id" bson:"_id" db:"autogenerate" validate:"max=255"`
    Name           string             `json:"name" bson:"name" validate:"max=255"`
    Email          string             `json:"email" bson:"email" db:"unique" validate:"email"`
    UpdatedAt      string             `json:"updated_at" bson:"updated_at" validate:"max=255"`
}
```

This tool automatically maps the JSON definitions (including additional metadata such as database constraints) to the appropriate Go struct with JSON, BSON, and validation tags.

---

### Example Server

An example of how to use Naboobase to create an HTTP server is provided in `example/server.go`. This server demonstrates:

- Setting up a MongoDB connection
- Initializing the API server on a specified host and port
- Attaching API endpoints including user creation and a health check
- Integrating an authentication layer

Below is the core snippet from `example/server.go`:

```go
package main

import (
	"naboobase/controllers"
	"naboobase/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Good")
	}
}

var dbConnector = core.MongoDBconnector{}

func main() {
	// Connect to the MongoDB database (replace "naboobase" with your database name)
	_ = dbConnector.Connect("naboobase")

	// Initialize the API server on localhost:1555
	myApi := core.Server{}
	myApi.Init("localhost", 1555)

	// Attach endpoints to the server
	myApi.AttachEndpoints([]core.Endpoint{
		{
			Method:  "POST",
			Path:    "/user",
			Handler: controllers.CreateUser(dbConnector),
		},
		{
			Method:  "GET",
			Path:    "/health",
			Handler: HealthCheck(),
		},
	})

	// Attach authentication middleware
	myApi.AttachAuthenticationLayer(dbConnector)

	// Run the server
	myApi.RunServer()
}
```

To run the example server:

```bash
go run example/server.go
```

This will start the server on `localhost:1555`. You can then test the endpoints (e.g., `POST /user` or `GET /health`) using tools like [Postman](https://www.postman.com) or `curl`.

---

## Code Structure

- **cli**: Contains a command-line tool that converts JSON definitions to Go structs.
- **configs**: Holds environment and configuration settings.
- **constants**: Global constants used across the project.
- **controllers**: Contains the logic for handling API requests (e.g., user operations).
- **core**: Implements the core server functionality, database connection, routing, and authentication mechanisms.
- **example**: Provides an example of setting up and running the server.
- **json**: Example JSON files that define models, requests, and responses.
- **models**: Defines the data models (structs) with proper JSON/BSON tags and validations.
- **routes**: (If needed) Can contain route definitions and grouping logic.
- **utils**: Utility functions for common tasks such as authentication, struct handling, and JSON conversion.

---

## Contributing

Contributions are welcome! If you’d like to contribute, please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes and ensure tests pass.
4. Submit a pull request describing your changes.

For major changes, please open an issue first to discuss what you would like to change.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

Happy coding! If you have any questions or need further assistance, feel free to open an issue or reach out via our [GitHub Discussions](https://github.com/lodjim/naboobase/discussions).

---


