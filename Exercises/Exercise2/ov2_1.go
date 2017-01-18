package main

import (
	"runtime"
	"sync"
)

var mu sync.Mutex

var i int = 0

func someGoroutine(done chan bool) {
	for j := 0; j < 1000000; j++ {
		mu.Lock()

		i++
		mu.Unlock()
	}
	done <- true
}

func somesome(done chan bool) {
	for j := 0; j < 1000001; j++ {

		mu.Lock()

		i--
		mu.Unlock()
	}
	done <- true
}

func main() {

	done := make(chan bool)

	runtime.GOMAXPROCS(runtime.NumCPU())
	go someGoroutine(done)
	go somesome(done)
	<-done
	<-done
	print(i, "\n")
}
