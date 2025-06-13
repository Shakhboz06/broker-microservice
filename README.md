# Go Microservices Sample

A sample microservices ecosystem in Go, demonstrating:

- **Broker Service**  
  Central HTTP/gRPC/RPC entrypoint forwarding to downstream services.
- **Auth Service**  
  User registration & login backed by PostgreSQL.
- **Logger Service**  
  HTTP and RPC logging to MongoDB.
- **Mailer Service**  
  Email dispatch via MailHog.
- **Listener Service**  
  RabbitMQ consumer for asynchronous event handling.
- **Frontend**  
  Minimal Go‐templated UI calling the broker.
- **Caddy**  
  Acts as a reverse‐proxy & HTTPS endpoint.

---

## 📦 Architecture

```text
┌────────┐      ┌────────────┐      ┌──────────────┐
│Client/ │  →   │ Caddy HTTP │  ↔   │ Frontend UI  │
│Browser │      └────────────┘      └──────────────┘
│        │            ↓
│        │       (proxy to)
│        │            ↓
│        │  ┌──────────────────────────┐
│        └─▶│   Broker Service (8080)  │
└────────┘  └──────────────────────────┘
                 │    │    │    │
      ┌──────────┼────┘    │    └─ RPC to Logger Service
      │          │         └───── HTTP to Mailer Service
      │          └─────────────── HTTP to Auth Service
      └─────────────────────────── AMQP ⇒ Listener Service


## 🛠️ Tech Stack
Language: Go 1.24

Messaging: RabbitMQ

Databases: PostgreSQL (Auth), MongoDB (Logger)

Containerization: Docker, Docker Swarm, Kubernetes (optional)

Ingress: Caddy (or NGINX)

Testing: Go’s net/http/httptest, testing pkg

## 📖 API Endpoints

| Path        | Method | Payload                                   | Description                          |
| ----------- | ------ | ----------------------------------------- | ------------------------------------ |
| `/`         | POST   | *empty*                                   | Health check / ping                  |
| `/handle`   | POST   | `{ action: "auth" \| "log" \| "mail", …}` | Dispatches to Auth / Logger / Mailer |
| `/log-grpc` | POST   | `{ action:"log", log:{ name, data } }`    | Log via gRPC → Logger Service        |

## Auth Service

| Path            | Method | Payload               | Resp               |
| --------------- | ------ | --------------------- | ------------------ |
| `/authenticate` | POST   | `{ email, password }` | `{ token }` or 401 |

## Logger Service
HTTP: POST /log { name, data } → 202 Accepted
RPC: Port 5001, methods: LogInfo(ctx, RpcPayload, *string)
