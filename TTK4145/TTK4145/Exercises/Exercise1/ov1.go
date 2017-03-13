package main

import (
	. "fmt"
	"runtime"
	"time"
)

var i int = 0

func someGoroutine() {
	println("From the routine")
	for j := 0; j < 10000; j++ {
		i++
	}
}

func somesome() {
	println("Prinsss")
	for j := 0; j < 10000; j++ {
		i--
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go someGoroutine()
	go somesome()
	time.Sleep(100 * time.Millisecond)
	Println("Print")
	print(i, "\n")
}
