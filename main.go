package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	workTime        = 2 * time.Second
	workersNum      = 5
	closeWorkersNum = 2
)

type WorkerPool struct {
	input   chan string
	closers []chan struct{}
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	wp := &WorkerPool{
		input:   make(chan string),
		closers: make([]chan struct{}, 0),
		mutex:   sync.Mutex{},
		wg:      sync.WaitGroup{},
	}

	return wp
}

func (wp *WorkerPool) Stop() {
	close(wp.input)
	wp.wg.Wait()
	for _, ch := range wp.closers {
		close(ch)
	}
}

func (wp *WorkerPool) AddWork(work string) {
	wp.input <- work
}

func (wp *WorkerPool) CloseWorker() error {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if len(wp.closers) > 1 {
		close(wp.closers[len(wp.closers)-1])
		wp.closers = wp.closers[:len(wp.closers)-1]
		return nil
	} else {
		return errors.New("cannot close worker when in worker pool only one worker")
	}
}

func (wp *WorkerPool) AddWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	closeChan := make(chan struct{})
	wp.closers = append(wp.closers, closeChan)

	wp.wg.Add(1)
	go func(closeChan chan struct{}, workerID int) {
		wp.startWorkByWorker(workerID)
		for {
			select {
			case <-closeChan:
				wp.endWorkByWorkerWithClose(workerID)
				wp.wg.Done()
				return
			default:
				if work, ok := <-wp.input; ok {
					wp.workByWorker(workerID, work)
				} else {
					wp.endWorkByWorker(workerID)
					wp.wg.Done()
					return
				}
			}
		}
	}(closeChan, len(wp.closers))
}

func (wp *WorkerPool) startWorkByWorker(workerID int) {
	fmt.Printf("worker № %d start work\n", workerID)
}

func (wp *WorkerPool) workByWorker(workerID int, input string) {
	time.Sleep(workTime)
	fmt.Printf("worker № %d recive %s\n", workerID, input)
}

func (wp *WorkerPool) endWorkByWorker(workerID int) {
	fmt.Printf("worker № %d end work\n", workerID)
}

func (wp *WorkerPool) endWorkByWorkerWithClose(workerID int) {
	fmt.Printf("worker № %d end work with close\n", workerID)
}

func main() {
	fmt.Printf("Start worker pool with %d workers...\n", workersNum)
	wp := NewWorkerPool()
	defer wp.Stop()

	for i := 0; i < workersNum; i++ {
		wp.AddWorker()
	}

	for i := 0; i < closeWorkersNum; i++ {
		err := wp.CloseWorker()
		if err != nil {
			fmt.Println(err)
		}
	}

	names := []string{"Alex", "Misha", "Dmitry", "Vladimir", "Petr", "Nikolay", "Diana",
		"Oleg", "Pasha", "Victor", "Artem", "George", "Vasily", "Anton",
	}

	for _, name := range names {
		wp.AddWork(name)
	}
}
