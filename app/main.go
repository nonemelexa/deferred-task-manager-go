package main

import (
	"fmt"
	"log"
	"os"

	"task-scheduler/internal/api"
	"task-scheduler/internal/queue"
	"task-scheduler/internal/scheduler"
	"task-scheduler/internal/storage"
	"task-scheduler/internal/worker"

	"github.com/joho/godotenv"
)

// Main function initializes the application by loading environment variables, setting up the queue and storage, starting the scheduler and worker goroutines, and finally starting the API server to listen for incoming requests.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	q := queue.NewRabbitMQ()
	s := storage.NewRedisStorage()

	handler := api.NewHandler(q, s)

	sched := scheduler.NewScheduler(s, q)
	go sched.Start()

	processor := &worker.DefaultProcessor{}

	for i := 0; i < 3; i++ {
		w := worker.NewWorker(i, q, s, processor)
		go w.Start()
	}

	fmt.Println("Server starting on port:", port)

	api.StartServer(handler, port)
}
