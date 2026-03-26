package worker

import (
	"fmt"
	"time"

	"task-scheduler/internal/domain"
)

type DefaultProcessor struct{}

// Process executes the task based on its type, simulating different processing times for each type of task and returning an error if the task type is unknown.
func (p *DefaultProcessor) Process(task domain.Task) error {
	switch task.Type {

	case "email":
		return sendEmail(task.Payload)

	case "payment":
		return processPayment(task.Payload)

	case "report":
		return generateReport(task.Payload)

	default:
		return fmt.Errorf("unknown task type: %s", task.Type)
	}
}

func sendEmail(payload string) error {
	fmt.Println("📧 Sending email to:", payload)
	time.Sleep(500 * time.Millisecond)
	return nil
}

func processPayment(payload string) error {
	fmt.Println("💳 Processing payment:", payload)
	time.Sleep(700 * time.Millisecond)
	return nil
}

func generateReport(payload string) error {
	fmt.Println("📊 Generating report:", payload)
	time.Sleep(1 * time.Second)
	return nil
}
