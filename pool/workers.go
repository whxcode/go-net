package pool

import (
	"sync"
)

type Task func(args ...interface{})

type worker struct {
	wait sync.WaitGroup
}

func NewWorkers(size int) *worker {
	w := &worker{}

	// go w.Execute()

	return w
}

func (w *worker) Post(task Task, args ...interface{}) {
	w.wait.Add(1)

	go func() {
		task(args...)
		w.wait.Done()
	}()
}

func (w *worker) Wait() {
	w.wait.Wait()
}
