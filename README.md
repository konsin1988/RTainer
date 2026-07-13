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

### Docker Integration (gRPC agent endpoints)

#### CONTAINERS 

grpcurl -plaintext localhost:50051 list   

##### list
> grpcurl -plaintext -d '{}' localhost:50051 agent.AgentService/ListContainers

##### start / stop
> grpcurl -plaintext -d '{"id":"<container_id>"}' localhost:50051 agent.AgentService/StopContainer

##### delete
> grpcurl -plaintext -d '{"id":"5ad2dd0ea57b152ded0be23a713d2d10a8cdff0953f87110d0a68dc628039d93", "force": "true", "remove_volumes": "false"}' localhost:50051 agent.AgentService/RemoveContainer

##### run
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


##### logs
> grpcurl -plaintext -d '{"container_id":"abc123","tail":20}' localhost:50051 agent.AgentService/ViewLogs

> grpcurl -plaintext -d '{"container_id":"abc123","follow":true}' localhost:50051 agent.AgentService/ViewLogs


##### inspect container
> grpcurl -plaintext -d '{"id":"cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca"}' localhost:50051 agent.AgentService/InspectContainer

###### response: 
```
{
  "id": "cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca",
  "name": "/my-test-nginx",
  "image": "sha256:a97d82f709e2e0ef35e48a697aec860e12cbf2a0ffbfd95d7701976e81d470ed",
  "status": "running",
  "ports": [
    {
      "containerPort": "80/tcp",
      "hostIp": "0.0.0.0",
      "hostPort": "8019"
    },
    {
      "containerPort": "80/tcp",
      "hostIp": "::",
      "hostPort": "8019"
    }
  ],
  "mounts": [
    {
      "source": "/tmp",
      "target": "/tmp"
    }
  ],
  "env": [
    "ENV=dev",
    "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "NGINX_VERSION=1.28.0",
    "PKG_RELEASE=1",
    "DYNPKG_RELEASE=1",
    "NJS_VERSION=0.8.10",
    "NJS_RELEASE=1"
  ]
}
```


##### container stats
> grpcurl -plaintext -d '{"id":"cbae295e09a0f533d2e1b9a1060de036df854b9730e7a19183cb4da377dad5ca"}' localhost:50051 agent.AgentService/ContainerStats

###### response:
>{
>  "memoryUsage": "17657856",
>  "memoryLimit": "16081117184",
>  "networkRx": "16799",
>  "networkTx": "126",
>  "pids": 13
>}


##### execute command
grpcurl -plaintext -d '{ "container_id":"2eeebdcd26c91ac36cde73609b6991538d51172e29d77d2c8b14bf758c8366ea", "command":["ls","-la","/"] }' localhost:50051 agent.AgentService/ExecuteCommand
grpcurl -plaintext -d '{ "container_id":"2eeebdcd26c91ac36cde73609b6991538d51172e29d77d2c8b14bf758c8366ea", "command":["env"] }' localhost:50051 agent.AgentService/ExecuteCommand

###### response:
```
...
{
  "line": "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
  "stream": "STDOUT"
}
{
  "line": "HOSTNAME=ed4e3dc26b29",
  "stream": "STDOUT"
}
{
  "line": "ENV=dev",
  "stream": "STDOUT"
}
...
```

#### IMAGES 

##### list
> grpcurl -plaintext -d '{}' localhost:50051 agent.AgentService/ListImages


##### remove
> grpcurl -plaintext -d '{"id":"sha256:a97d82f709e2e0ef35e48a697aec860e12cbf2a0ffbfd95d7701976e81d470ed", "force":true }' localhost:50051 agent.AgentService/RemoveImage

#### NETWORKS 

##### list
> grpcurl -plaintext -d '{}' localhost:50051 agent.AgentService/ListNetworks

###### response: 
```
{
  "networks": [
    {
      "id": "d0fd8f8a4c998c2cfed8610a19862ddf4733c7836c8338715dec17d38cdcab10",
      "name": "none",
      "driver": "null",
      "scope": "local"
    },
    {
      "id": "d322401759b1a25b698253deba793673c6246b5bb289428951d94f6e16611cd8",
      "name": "rtainer",
      "driver": "bridge",
      "scope": "local",
      "containers": [
        {
          "id": "ae89ee158e15e6ad502ea635845909f1a790850bc16987e4730cd538be6a9632",
          "name": "postgres",
          "ipv4Address": "172.22.0.3/16"
        },
        {
          "id": "d982b751aef9d3c9177945f3afdf22f9445649bb66add464bda1d10c5b1eedf7",
          "name": "keycloak",
          "ipv4Address": "172.22.0.2/16"
        },
        {
          "id": "e2ba9fe693b6f3d317225d350469cbd4c76c8b92d91de3294c631580cbd998ec",
          "name": "rtainer-agent",
          "ipv4Address": "172.22.0.5/16"
        },
        {
          "id": "fabcc775a2fae80f65ceea5b4dc323d7cc2e601df76cd4027c494b2923a5292e",
          "name": "rtainer-dev",
          "ipv4Address": "172.22.0.4/16"
        }
      ]
    },
...
}
```

##### create
> grpcurl -plaintext -d '{ "name":"my-test-network", "driver":"bridge" }' localhost:50051 agent.AgentService/CreateNetwork

##### remove
> grpcurl -plaintext -d '{ "id":"8e4f5f..." }' localhost:50051 agent.AgentService/RemoveNetwork


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

