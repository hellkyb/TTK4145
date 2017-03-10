package main

import (
	"fmt"
)

func mapManipulation(someMap map[string]bool) {
	someMap["IP:666"] = true

}

func main() {
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
