# RTainer 

A lightweight container and infrastructure management platform

Built with:
- Go (backend API)
- React (Vite frontend)
- PostgreSQL (database)
- Docker API integration
- JWT authentication
- Database migrations (golang-migrate)

🚀 Features (MVP)

- User authentication (JWT-based)
- Role-based access control (RBAC)
- Manage Docker environments
- Container management (start/stop/list)
- Stack (Docker Compose) deployment
- Audit logging
- PostgreSQL-backed persistence
- Auto database migrations on startup

🏗️ Architecture

React (Vite)
↓
HTTP API (Go)
↓
PostgreSQL
↓
Docker Engine (via API)

📦 Tech Stack

Backend
- Go 1.22+
- net/http
- PostgreSQL
- golang-migrate
- lib/pq
- golang-jwt/jwt/v5

Frontend
- React + Vite
- TypeScript
- React Router
- TailwindCSS (optional)

Infra
- Docker
- Docker Compose

📁 Project Structure

/config
  config.go

/db
  postgres.go
  migrate.go
  health

/migrations
  000001_init_schema.up.sql
  000001_init_schema.down.sql

/transport
  handlers.go

main.go

🗄️ Database Schema

Core entities:
- users
- roles
- user_roles
- environments
- stacks
- audit_logs

Migrations run automatically on startup.

🐳 Docker Integration

Capabilities:
- List containers
- Start/stop containers
- View logs (planned)
- Execute commands (planned)

📡 API Endpoints

GET /health

GET /containers
POST /containers/:id/start
POST /containers/:id/stop

POST /auth/login
POST /auth/register

🧠 Design Principles

- Migration-based schema (no auto ORM sync)
- Clean layered architecture
- Stateless authentication (JWT)
- PostgreSQL relational model
- Extensible agent-based future design

📈 Roadmap

MVP
- DB setup
- Migrations
- Auth foundation
- Basic API structure

Next
- Container management
- WebSocket logs
- Docker stats streaming
- React dashboard

Future
- Multi-node agents
- Kubernetes support
- Keycloak integration
- Redis event bus
- GraphQL (optional)
