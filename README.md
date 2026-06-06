# shortlink-backend

Backend service for a URL shortening platform written in Go with Fiber. Handles
user authentication, URL shortening, click tracking, JWT token management, and
session storage. Built with security, performance, and scalability in mind.

## Stack

- **[Go](https://go.dev/)** and **[Fiber](https://gofiber.io/)** —
  high-performance HTTP server
- **[PostgreSQL](https://www.postgresql.org/)** — user accounts, links, and
  analytics storage
- **[Redis](https://redis.io/)** — cache and session storage
- **[JWT](https://jwt.io/)** — token-based authentication
- **[bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)** — password hashing
- **[Swagger](https://swagger.io/)** — API documentation
- **[Docker](https://www.docker.com/)** — containerized deployment

## Endpoints

## Endpoints

| Method | Path                           | Description                                     | Auth Required |
| ------ | ------------------------------ | ----------------------------------------------- | ------------- |
| `POST` | `/api/v1/auth/register`        | Register a new user                             | ❌            |
| `POST` | `/api/v1/auth/login`           | Login and receive access/refresh tokens         | ❌            |
| `POST` | `/api/v1/auth/refresh`         | Refresh access token                            | ❌            |
| `POST` | `/api/v1/auth/logout`          | Logout and revoke current session               | ✅            |
| `POST` | `/api/v1/user/profile`         | Get information about your profile              | ✅            |
| `POST` | `/api/v1/link/shortlink`       | Create a new short link                         | ✅            |
| `GET`  | `/api/v1/link/shortlink/:link` | Resolve short link and redirect to original URL | ❌            |
| `GET`  | `/api/v1/link/activate`        | Get link activity and visit statistics          | ✅            |

## Running locally

### Requirements

- Go 1.26.2+
- PostgreSQL
- Redis

### Setup

```bash
git clone https://github.com/whymewaomi/shortlink-backend.git
cd shortlink-backend

cp .env.example .env

go mod tidy
go run cmd/main.go
```

## Running with Docker

```bash
docker compose up -d --build
```

## Features

- JWT authentication
- Refresh token rotation
- URL shortening
- Click tracking
- User activity analytics
- Redis caching
- PostgreSQL persistence
- Swagger documentation
- Docker deployment
- Fiber middleware support
- Secure password hashing with bcrypt
