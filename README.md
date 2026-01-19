Distributed WebSocket Application with UI and Backend

This project provides a full-stack WebSocket-based application with two WebSocket servers, a Redis instance, and an Nginx reverse proxy. It includes a web-based UI and backend services, all containerized using Docker Compose.

Table of Contents 

	•	Features￼
	•	Architecture￼
	•	Prerequisites￼
	•	Getting Started￼
	•	Usage￼
	•	Docker Compose Services￼
	•	Configuration￼
	•	Troubleshooting￼
	•	License￼

⸻

Features
	•	Two WebSocket servers: Handles real-time communication.
	•	Redis instance: Used as a shared store or message broker.
	•	Nginx reverse proxy: Serves both frontend and backend on separate ports. Load balances backend servers.
	•	Web UI: Accessible via http://localhost:3030.
	•	Backend API: Accessible via http://localhost:8080.
	•	Fully containerized with Docker Compose.

⸻

Architecture

          ┌───────────────────┐
          │     Web UI        │
          │  Port: 3030       │
          └─────────┬─────────┘
                    │
                    ▼
              ┌─────────┐
              │ Nginx   │
              │ Reverse │
              │ Proxy   │
              └─────────┬─────────┐
                        │          │
                        ▼          ▼
                 ┌─────────┐  ┌─────────┐
                 │ WS #1   │  │ WS #2   │
                 └─────────┘  └─────────┘
                        │           |
                        ▼           |
                    ┌─────────┐ <---|   
                    │ Redis   │
                    │ Store   │
                    └─────────┘

	•	Nginx routes traffic to the appropriate service.
	•	Web UI communicates with backend via WebSocket servers.
	•	Redis can be used for pub/sub or shared state between WebSocket servers.

⸻

Prerequisites

Before running this project, make sure you have the following installed:
	•	Docker￼ (v20+ recommended)
	•	Docker Compose￼ (v2+ recommended)

⸻

Getting Started
	1.	Clone the repository

git clone <repository-url>
cd <repository-directory>

	2.	Build and start the services

docker compose up --build

This command will start:
	•	Two WebSocket servers
	•	Redis instance
	•	Nginx reverse proxy

⸻

Usage
	•	Access Web UI: http://localhost:3030￼
	•	Access Backend Ws API: http://localhost:8080/ws￼

Interacting with WebSockets
	1.	The UI connects to the WebSocket servers automatically.
	2.	Backend services are exposed via Nginx on port 8080.
	3.	Redis is used for messaging between WebSocket servers if needed.

⸻

Docker Compose Services

Service	Port	Description
ws-server-1	N/A	WebSocket server #1
ws-server-2	N/A	WebSocket server #2
redis	6379	Redis instance
nginx	3030, 8080	Reverse proxy for UI and backend


⸻

Configuration
	•	Nginx: Configured to forward requests to:
	•	/ → Web UI (3030)
	•	/ws → Backend (8080)
	•	Redis: Default port 6379. No authentication configured by default.
	•	WebSocket servers: Configured via environment variables (if applicable, e.g., Redis URL).

⸻

Troubleshooting
	•	Docker Compose fails to start
	•	Make sure Docker and Docker Compose are installed.
	•	Ensure no other services are using ports 3030 or 8080.
	•	UI not loading
	•	Verify Nginx container is running:

docker ps


	•	Check Nginx logs for errors:

docker logs <nginx-container-id>


	•	WebSocket connection issues
	•	Ensure Redis is running and accessible.
	•	Verify WebSocket server logs.
