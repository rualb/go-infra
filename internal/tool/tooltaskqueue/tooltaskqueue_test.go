package tooltaskqueue

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// Task for testing purposes
type TestTask struct {
	value int32
}

// Test panic recovery inside handler
func TestTaskQueue_PanicRecovery(t *testing.T) {
	var processed atomic.Int32

	handler := func(task *TestTask) error {
		if task.value == 2 {
			panic("test panic")
		}
		processed.Add(task.value)
		return nil
	}

	queue := NewTaskQueue("panicQueue", handler, 2)
	queue.SetActive(true)

	// Enqueue tasks, one of which will cause a panic
	_ = queue.Enqueue(&TestTask{value: 1})
	_ = queue.Enqueue(&TestTask{value: 2}) // This should cause a panic
	_ = queue.Enqueue(&TestTask{value: 3})

	time.Sleep(500 * time.Millisecond)

	// Check that processing continued for non-panic tasks
	if processed.Load() != 4 {
		t.Errorf("Expected processed to be 4, got %d", processed.Load())
	}
}

// Test task queue error handling in handler
func TestTaskQueue_ErrorHandling(t *testing.T) {

	handler := func(task *TestTask) error {
		if task.value == 2 {
			return errors.New("test error")
		}
		return nil
	}

	queue := NewTaskQueue("errorQueue", handler, 2)
	queue.SetActive(true)

	// Enqueue tasks, one of which will cause an error
	_ = queue.Enqueue(&TestTask{value: 1})
	_ = queue.Enqueue(&TestTask{value: 2}) // This will cause an error
	_ = queue.Enqueue(&TestTask{value: 3})

	time.Sleep(500 * time.Millisecond)

	// No panic should have occurred, just error handling in logs
	// You could also capture logs if needed, but for simplicity, it's not done here.
}
