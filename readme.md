🚀 Task Scheduler (Go)

An asynchronous task scheduler built with Go, supporting:

delayed task execution
message queue (RabbitMQ)
retry mechanism with exponential backoff
task prioritization


🧠 Architecture
Client → API → 
   ├── RabbitMQ (instant tasks)
   └── Redis (delayed tasks)

Scheduler → Redis → RabbitMQ

Workers → RabbitMQ → processing → retry / fail


⚙️ Tech Stack
Go
RabbitMQ
Redis
Docker / Docker Compose



🚀 Getting Started
1. Clone the repository

bash
git clone <your-repo>
cd goproject7


2. Run with Docker
bash 
docker-compose up --build



📂 Project Structure
goproject7/
  app/              # main.go (entry point)
  internal/
    api/            # HTTP handlers
    domain/         # models
    queue/          # RabbitMQ logic
    storage/        # Redis logic
    scheduler/      # delayed task scheduler
    worker/         # task processors
  docker-compose.yml
  Dockerfile

🧠 Summary

This project demonstrates:

asynchronous task processing
message queues
delayed execution using Redis
retry mechanisms with backoff