package goSquare

import "sync"

func GetSquare(jobs chan int, square chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	go func() {
		for job := range jobs {
			square <- job * job
		}
		close(square)
	}()

	go func() {
		for i := 0; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()
}
