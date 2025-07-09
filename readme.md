# Serra Chat App Backend

Serra is a chat application backend built with Go. This service provides RESTful APIs and WebSocket endpoints for real-time messaging.

## Features

- User authentication (JWT)
- Real-time chat with WebSockets
- Message history and persistence
- Scalable architecture

## Getting Started

### Prerequisites

- Go 1.20+
- Docker (optional, for database)
- MySQL

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/serra.git
   cd serra
   ```

2. Create a `.env` file in the project root (see format below).

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

### .env File Format

Create a `.env` file in the project root with the following content:

```env
PUBLIC_HOST=http://localhost
PORT=8080
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=serra
JWT_SECRET=your_jwt_secret
```

- `PUBLIC_HOST`: Base URL for the server.
- `PORT`: Port for the server to listen on.
- `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`: MySQL database connection settings.
- `JWT_SECRET`: Secret key for JWT authentication.

## API Documentation

See [API.md](API.md) for detailed endpoints.

## License

MIT
