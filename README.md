# RTainer 

A lightweight container and infrastructure management platform

### Built with:
- Go (backend API)
- React (Vite frontend)
- PostgreSQL (database)
- Docker API integration
- JWT authentication
- Database migrations (golang-migrate)

### Features 

- User authentication (JWT-based)
- Role-based access control (RBAC)
- Manage Docker environments
- Container management (start/stop/list)
- Stack (Docker Compose) deployment
- Audit logging
- PostgreSQL-backed persistence
- Auto database migrations on startup

### Architecture

React (Vite)
↓
HTTP API (Go)
↓
PostgreSQL
↓
Docker Engine (via API)

### Tech Stack

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

### Database Schema

Core entities:
- users
- roles
- user_roles
- environments
- stacks
- audit_logs

Migrations run automatically on startup.

### Docker Integration (gRPC agent endpoint)

#### CONTAINERS 

grpcurl -plaintext localhost:50051 list   

 ```list```
> grpcurl -plaintext -d '{}' localhost:50051 agent.AgentService/ListContainers

```start / stop```
> grpcurl -plaintext -d '{"id":"<container_id>"}' localhost:50051 agent.AgentService/StopContainer

```delete```
> grpcurl -plaintext -d '{"id":"5ad2dd0ea57b152ded0be23a713d2d10a8cdff0953f87110d0a68dc628039d93", "force": "true", "remove_volumes": "false"}' localhost:50051 agent.AgentService/RemoveContainer

```run```
> grpcurl -plaintext -d '{"image_id":"5ad2dd0ea57b152ded0be23a713d2d10a8cdff0953f87110d0a68dc628039d93","name":"my-test-container"}' localhost:50051 agent.AgentService/RunContainer

>grpcurl -plaintext -d '{
>  "image_id": "sha256:a97d82f709e2e0ef35e48a697aec860e12cbf2a0ffbfd95d7701976e81d470ed",
>  "name": "my-test-nginx",
>  "command": [
>    "nginx",
>    "-g", 
>    "daemon off;"
>  ],
>  "env": [
>    "ENV=dev"
>  ],
>  "ports": [
>    {
>      "container_port": "80/tcp",
>      "host_port": "8019"
>    }
>  ],
>  "volumes": [
>    {
>      "source": "/tmp",
>      "target": "/tmp"
>    }
>  ],
>  "tty": false,
>  "detach": true
>}' localhost:50051 agent.AgentService/RunContainer


```logs```
> grpcurl -plaintext -d '{"container_id":"abc123","tail":20}' localhost:50051 agent.AgentService/ViewLogs

> grpcurl -plaintext -d '{"container_id":"abc123","follow":true}' localhost:50051 agent.AgentService/ViewLogs


```inspect container```
> grpcurl -plaintext -d '{"id":"cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca"}' localhost:50051 agent.AgentService/InspectContainer

##### response: 
> {
>   "id": "cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca",
>   "name": "/my-test-nginx",
>   "image": "sha256:a97d82f709e2e0ef35e48a697aec860e12cbf2a0ffbfd95d7701976e81d470ed",
>   "status": "running",
>   "ports": [
>     {
>       "containerPort": "80/tcp",
>       "hostIp": "0.0.0.0",
>       "hostPort": "8019"
>     },
>     {
>       "containerPort": "80/tcp",
>       "hostIp": "::",
>       "hostPort": "8019"
>     }
>   ],
>   "mounts": [
>     {
>       "source": "/tmp",
>       "target": "/tmp"
>     }
>   ],
>   "env": [
>     "ENV=dev",
>     "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
>     "NGINX_VERSION=1.28.0",
>     "PKG_RELEASE=1",
>     "DYNPKG_RELEASE=1",
>     "NJS_VERSION=0.8.10",
>     "NJS_RELEASE=1"
>   ]
> }


```container stats```
> grpcurl -plaintext -d '{"id":"cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca"}' localhost:50051 agent.AgentService/ContainerStats

##### response:
>{
>  "memoryUsage": "17657856",
>  "memoryLimit": "16081117184",
>  "networkRx": "16799",
>  "networkTx": "126",
>  "pids": 13
>}

#### IMAGES 

```list```
> grpcurl -plaintext -d '{}' localhost:50051 agent.AgentService/ListImages


```remove```
> grpcurl -plaintext -d '{"id":"sha256:a97d82f709e2e0ef35e48a697aec860e12cbf2a0ffbfd95d7701976e81d470ed", "force":true }' localhost:50051 agent.AgentService/RemoveImage


### API Endpoints

GET /health

GET /containers
POST /containers/:id/start
POST /containers/:id/stop

POST /auth/login
POST /auth/register

### Design Principles

- Migration-based schema (no auto ORM sync)
- Clean layered architecture
- Stateless authentication (JWT)
- PostgreSQL relational model
- Extensible agent-based future design

