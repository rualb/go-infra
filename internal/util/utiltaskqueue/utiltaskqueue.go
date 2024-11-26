// Package utiltaskqueue task queue tool
package utiltaskqueue

import (
	"container/list"
	"fmt"
	xlog "go-infra/internal/util/utillog"
	"sync"
)

// TaskQueueStats queue stats
type TaskQueueStats struct {
	QueueSize   int
	WorkerCount int
	MaxWorker   int
}

// TaskQueue task queue
type TaskQueue[T any] struct {
	handler       func(*T) error
	list          list.List
	mu            sync.Mutex
	workerCounter int // not atomic allowed
	maxWorker     int
	isActive      bool // not atomic allowed, durty read-write allowed
	name          string
	MaxQueueSize  int // durty read-write allowed
}

// Stats get stats
func (x *TaskQueue[T]) Stats() TaskQueueStats {
	// trigger for processing

	return TaskQueueStats{
		QueueSize:   x.list.Len(),
		WorkerCount: x.workerCounter,
		MaxWorker:   x.maxWorker,
	}

}

// SetActive start-stop queue
func (x *TaskQueue[T]) SetActive(value bool) {
	// trigger for processing

	x.isActive = value

}

// Enqueue add to queue
func (x *TaskQueue[T]) Enqueue(data *T) error {
	// trigger for processing

	x.pushData(data)
	x.tryRunWorker()

	return nil
}

func (x *TaskQueue[T]) tryRunWorker() {

	if !x.isActive {
		return
	}

	// durty read, protect from no-resource run
	if x.hasDataUnsafe() && x.hasFreeWorkerUnsafe() {
		go x.runWorker() // async goroutine
	}
}
func (x *TaskQueue[T]) allocateWorker() bool {
	res := false
	x.mu.Lock()
	defer x.mu.Unlock()

	if x.hasFreeWorkerUnsafe() {
		res = true
		x.workerCounter++ // inc
	}

	return res
}
func (x *TaskQueue[T]) releaseWorker() {
	x.mu.Lock()
	defer x.mu.Unlock()
	// Race Condition

	x.workerCounter-- // decr
}
func (x *TaskQueue[T]) runWorker() {

	{
		gotWorker := x.allocateWorker()
		//
		if gotWorker {

			defer x.releaseWorker()
			defer x.tryRunWorker() // call back

			for x.isActive { // loop if data exists

				data := x.popData()
				if data == nil {
					break
				}

				// Handle potential panic inside task handler
				err := func() (err error) {
					defer func() {

						if r := recover(); r != nil {
							// Log or handle the panic
							err = fmt.Errorf("error panic: %v", r)
							// err = fmt.Errorf("error panic: %v\n%s", r, debug.Stack())
						}
					}()
					return x.handler(data)
				}()

				if err != nil {
					xlog.Error("task queue %s: %v", x.name, err)
				}

			}

		}

	}

}

// hasFreeWorkerUnsafe is usafe and durty resula allowed
func (x *TaskQueue[T]) hasFreeWorkerUnsafe() bool {
	return x.workerCounter < x.maxWorker
}

// hasDataUnsafe is usafe and durty resula allowed
func (x *TaskQueue[T]) hasDataUnsafe() bool {
	return x.list.Len() > 0 // protect from loop
}

func (x *TaskQueue[T]) popData() *T {
	x.mu.Lock()
	defer x.mu.Unlock()

	if el := x.list.Back(); el != nil && el.Value != nil {
		x.list.Remove(el)
		data, _ := el.Value.(*T)

		return data
	}
	return nil
}
func (x *TaskQueue[T]) pushData(data *T) {

	if !x.isActive {
		return
	}

	if data == nil {
		return
	}

	x.mu.Lock()
	defer x.mu.Unlock()

	if x.MaxQueueSize > 0 && x.list.Len() > x.MaxQueueSize {
		xlog.Info("task queue %v  is overloaded", x.name)
		return
	}

	x.list.PushFront(data)

}

func NewTaskQueue[T any](name string, handler func(*T) error, maxWorker int) *TaskQueue[T] {

	return &TaskQueue[T]{
		handler:   handler,
		maxWorker: maxWorker,
		name:      name,
		isActive:  true,
	}
}
