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
   git clone https://github.com/IslahulKafa/serra-backend.git
   cd serra-backend
   ```

2. Create a `.env` file in the project root (see format below).

3. Install dependencies:

   ```bash
   go mod download
   ```

4. Run the server:
   ```bash
   make run
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

## Database Schema

### Users Table

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### Prekeys Table

```sql
CREATE TABLE IF NOT EXISTS prekeys (
    user_id BIGINT UNSIGNED NOT NULL,
    identity_key TEXT NOT NULL,
    signed_prekey TEXT NOT NULL,
    signed_prekey_signature TEXT NOT NULL,
    one_time_prekeys JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
```

## API Documentation

See [API.md](API.md) for detailed endpoints.

## License

MIT
