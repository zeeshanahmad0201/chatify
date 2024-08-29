# Chatify System Architecture

## Overview

Chatify is a real-time chat application built with Golang. It is designed with scalability, modularity, and maintainability in mind. The application leverages MongoDB for data storage and Gorilla WebSocket for real-time communication, while JWT-based authentication ensures secure access.

## High-Level Architecture

The architecture of Chatify follows a layered approach, separating the application into distinct components:

1. **API Layer**: Handles HTTP requests and responses, mapping routes to controllers.
2. **Service Layer**: Contains business logic and interacts with the data layer.
3. **Data Layer**: Responsible for database operations, abstracting MongoDB interactions.
4. **WebSocket Layer**: Manages real-time communication between clients and the server.

### Component Diagram

```plaintext
┌────────────────────┐
│     Web Client     │
│  (Browser/Flutter) │
└────────────────────┘
          │
          ▼
┌────────────────────┐
│    API Gateway     │
│ (HTTP Server - mux)│
└────────────────────┘
          │
          ▼
┌───────────────────────────────┐
│        API Controllers        │
│ - AuthController              │
│ - MessageController           │
└───────────────────────────────┘
          │
          ▼
┌───────────────────────────────┐
│        Service Layer          │
│ - AuthService                 │
│ - MessageService              │
│ - UserService                 │
└───────────────────────────────┘
          │
          ▼
┌───────────────────────────────┐
│        Data Layer             │
│ - MongoDB (CRUD Operations)   │
│ - Models                      │
└───────────────────────────────┘
          │
          ▼
┌───────────────────────────────┐
│        WebSocket Layer        │
│ - Real-Time Messaging         │
└───────────────────────────────┘
```

### Detailed Architecture

## API Layer

The API layer is responsible for handling incoming HTTP requests and delegating them to the appropriate controllers. The controllers map the API routes to service methods that contain the business logic.

	-	**Controllers**: Each controller corresponds to a specific domain (e.g., AuthController, MessageController) and contains methods for handling related API endpoints.

## Service Layer

The service layer contains the core business logic of the application. It is responsible for processing requests, interacting with the data layer, and ensuring that the business rules are followed.

	-	**AuthService**: Manages user authentication, including login, signup, and token validation.
	-	**MessageService**: Handles the creation, deletion, and updating of messages, including status tracking (sent, delivered, read).
	-	**UserService**: Manages user-related operations, such as fetching user details.

## Data Layer

The data layer abstracts the database interactions, making it easier to swap out the underlying database technology if needed. It uses MongoDB as the primary database for storing users and messages.

	-	**MongoDB**: Stores user data and chat messages. Collections are designed to optimize read and write operations for real-time messaging.

## WebSocket Layer

The WebSocket layer enables real-time communication between the server and clients. It allows for instant message delivery and status updates without requiring constant HTTP requests.

	-	**Gorilla WebSocket**: The WebSocket library used for managing real-time connections, enabling features like message broadcasting and user presence.

## Data Flow

	1. **User Authentication:**
	-	A user submits their credentials via the /login or /signup endpoint.
	-	The request is handled by AuthController, which calls AuthService to validate the user and generate a JWT.
	-	The JWT is sent back to the client for use in subsequent requests.
	2.	**Message Sending:**
	-	A user sends a message via the /message/send endpoint.
	-	The MessageController processes the request, and MessageService stores the message in the MongoDB collection.
	-	If the recipient is online, the message is pushed via WebSocket; otherwise, it is marked for delivery.
	3.	**Message Status Update:**
	-	When a message is read or delivered, the client sends a request to /message/status/read or /message/status/delivered.
	-	The MessageController updates the message status via MessageService, which persists the change in MongoDB.

### Future Enhancements

## Microservices Architecture

As the application scales, the monolithic structure can be refactored into microservices. This will allow individual components to scale independently, improving performance and reliability.

 - **End-to-End Encryption**

Implementing end-to-end encryption for messages would enhance the security of user data, making Chatify suitable for more sensitive communications.

- **Caching Layer**

Adding a caching layer (e.g., Redis) can improve the performance of frequently accessed data, such as user sessions or recent messages.

### Conclusion

Chatify’s architecture is designed to be scalable and maintainable, making it a solid foundation for building a robust chat application. By following best practices in API design, service separation, and real-time communication, it offers a reliable platform for future growth and enhancements.