# Chatify

Chatify is a Golang-based real-time chat application that supports user authentication, messaging, and real-time communication via WebSocket. It is built with a clean and modular structure, making it easy to maintain and extend.

## Features

- **User Authentication**: Secure login and signup functionalities.
- **Messaging**: Send and receive messages with status tracking (sent, delivered, read).
- **WebSocket Communication**: Real-time messaging using WebSocket.
- **MongoDB Integration**: Store and retrieve messages and user data.

## API Endpoints

### Authentication
- `POST /login` - User login
- `POST /signup` - User signup

### Messages
- `GET /messages` - Retrieve all messages
- `POST /message/send` - Send a message
- `DELETE /message/delete/{messageId}` - Delete a message
- `PUT /message/status/read` - Mark a message as read
- `PUT /message/status/delivered` - Mark a message as delivered

### WebSocket
- `GET /messages/listen` - WebSocket endpoint for real-time messaging

## Getting Started

### Prerequisites

- Go 1.16 or higher
- MongoDB

### Installation

1. **Clone the repository:**
    ```bash
    git clone https://github.com/zeeshanahmad0201/chatify.git
    ```

2. **Navigate to the project directory:**
    ```bash
    cd chatify
    ```

3. **Install dependencies:**
    ```bash
    go mod tidy
    ```

4. **Set up your MongoDB database** and update the connection string in `pkg/database/db_connection.go`.

5. **Run the application:**
    ```bash
    go run cmd/chatify/main.go
    ```

### Usage

After running the application, you can interact with the API using tools like [Postman](https://www.postman.com/) or [cURL](https://curl.se/).

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.
