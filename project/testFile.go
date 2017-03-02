package main

import (
	"fmt"
	"time"
)

var i int32 = 0
var slice []int

func timerFunction() {
	fmt.Println("Start")
	time.Sleep(2000 * time.Millisecond)
	fmt.Println("Timer finished")
}

func Appender() {

}

func krobar() {
	i = int32(time.Now().Unix())

}

func main() {
	x := 5
	y := 10
	var xPtr *int
	xPtr := &x
	i = int32(time.Now().Unix())
	fmt.Println("Hello")
	fmt.Println(i)

	timerFunction()
	i = int32(time.Now().Unix())
	fmt.Println(i)
}
