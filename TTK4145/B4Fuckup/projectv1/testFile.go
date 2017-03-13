package main

import (
	"fmt"
	"time"
)

func mapManipulation(someMap map[string]bool) {
	someMap["IP:666"] = true

}

var m []int64

func main() {
	m = append(m, time.Now().Unix())
	m = append(m[:0], m[1]...)

	fmt.Println("ok")
	myMap := make(map[string]bool)

	myMap["IP:123"] = true
	myMap["IP:666"] = false

	fmt.Print("Ip 123 :")
	fmt.Println(myMap["IP:123"])
	fmt.Print("Ip 666: ")
	fmt.Println(myMap["IP:666"])
	mapManipulation(myMap)
	fmt.Println(myMap["IP:666"])
	fmt.Println(myMap)
}
