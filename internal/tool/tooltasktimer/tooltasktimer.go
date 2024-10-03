// Package tooltasktimer run task by timer
package tooltasktimer

import (
	xlog "go-infra/internal/tool/toollog"
	"sync"
	"time"
)

// TaskTimer defines a struct that runs a task every N seconds
// and prevents concurrent execution using TryLock.
type TaskTimer struct {
	interval time.Duration // Interval between task executions
	task     func() error  // Task to execute
	mutex    sync.Mutex    // Mutex to lock running state
	stopChan chan struct{} // Channel to signal stopping of the timer
	Debug    bool
	name     string
}

// NewTaskTimer creates a new TaskTimer instance with the given interval and task.
func NewTaskTimer(name string, interval time.Duration, task func() error) *TaskTimer {
	return &TaskTimer{
		interval: interval,
		task:     task,
		stopChan: make(chan struct{}),
		name:     name,
	}
}

// Start begins the task timer, executing the task every interval.
func (t *TaskTimer) Start() {

	go func() {
		ticker := time.NewTicker(t.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Try to lock, if unable, skip the task.
				if t.mutex.TryLock() {
					go func() {
						defer func() {
							// Ensure any panic in the task does not crash the program.
							if r := recover(); r != nil {
								xlog.Info("Recovered from panic: %v", r)
							}
							// Mark the task as completed.
							t.mutex.Unlock()
						}()

						// Execute the task.
						err := t.task()
						if err != nil {
							xlog.Info("Error in task timer %s: %v", t.name, err.Error())
						}
					}()
				} else if t.Debug {
					xlog.Info("Previous task is still running, skipping this step.")
				}

			case <-t.stopChan:
				// Stop the timer.
				if t.Debug {
					xlog.Info("Task Timer stopped.")
				}
				return
			}
		}
	}()
}

// Stop stops the task timer.
func (t *TaskTimer) Stop() {
	close(t.stopChan)
}

// func main() {

// 	//5 ticks, 3 skips
// 	sampleTask := func() error {
// 		fmt.Println("Running task at:", time.Now())
// 		time.Sleep(5 * time.Second) // Simulate a task taking 3 seconds
// 		return fmt.Errorf("error 2: %v", fmt.Errorf("error 1"))
// 	}

// 	// Create a new TaskTimer with an interval of 5 seconds
// 	timer := NewTaskTimer("timerrrr", 2*time.Second, sampleTask)
// 	timer.Debug = true
// 	// Start the timer
// 	timer.Start()

// 	// Run for 20 seconds, then stop the timer
// 	time.Sleep(10 * time.Second)
// 	timer.Stop()
// 	time.Sleep(1 * time.Second)
// 	//fmt.Println("Task Timer stopped.")
// }
