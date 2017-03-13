package main

import (
	"runtime"
	"sync"
	"time"
)

var mu sync.Mutex

var i int = 0

func someGoroutine() {
	for j := 0; j < 1000000; j++ {
		mu.Lock()

		i++
		mu.Unlock()
	}
}

func somesome() {
	for j := 0; j < 1000001; j++ {

		mu.Lock()

		i--
		mu.Unlock()
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go someGoroutine()
	go somesome()
	time.Sleep(1000 * time.Millisecond)
	print(i, "\n")
}
