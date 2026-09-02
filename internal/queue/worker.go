package queue

import (
	"log"
	"sync"
	"time"
)

type NotificationTask struct {
	ProductID string
	Action    string // e.g., "CREATED", "UPDATED", "DELETED"
}

type TaskQueue struct {
	taskChan chan NotificationTask
	wg       sync.WaitGroup
}

func NewTaskQueue(bufferSize int, numWorkers int) *TaskQueue {
	q := &TaskQueue{
		taskChan: make(chan NotificationTask, bufferSize),
	}

	// Start background workers
	for i := 1; i <= numWorkers; i++ {
		go q.startWorker(i)
	}

	return q
}

// Enqueue adds a task without blocking the HTTP handler.
// If the buffer is full, it drops or logs instead of slowing the API down.
func (q *TaskQueue) Enqueue(task NotificationTask) {
	select {
	case q.taskChan <- task:
		// Successfully queued!
	default:
		// Buffer is full: Log warning instead of blocking customer request
		log.Printf("[QUEUE WARNING] Buffer full! Dropped task for product %s", task.ProductID)
	}
}

// startWorker loops indefinitely pulling tasks off the channel
func (q *TaskQueue) startWorker(workerID int) {
	log.Printf("Worker %d listening for background tasks...", workerID)

	for task := range q.taskChan {
		q.processTask(workerID, task)
	}
}

func (q *TaskQueue) processTask(workerID int, task NotificationTask) {
	q.wg.Add(1)
	defer q.wg.Done()

	// Simulate sending a WhatsApp message or external API sync (takes 2 seconds)
	time.Sleep(2 * time.Second)

	log.Printf("[Worker %d] Successfully processed background %s event for Product ID: %s",
		workerID, task.Action, task.ProductID)
}

func (q *TaskQueue) Stop() {
	log.Println("Closing background job channel...")
	close(q.taskChan) // Signals workers: "No more new tasks are coming!"

	log.Println("Waiting for active background workers to finish...")
	q.wg.Wait() // Blocks until all active processTask() calls call wg.Done()
	log.Println("All background workers drained successfully.")
}
