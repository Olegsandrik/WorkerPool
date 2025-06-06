package main

import (
	"fmt"
)

const (
	workersNum      = 5
	closeWorkersNum = 2
)

func main() {
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
