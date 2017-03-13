package main

import (
	"fmt"
	//"time"
)

func mapManipulation(someMap map[string]bool) {
	someMap["IP:666"] = true

}



func main() {
	fmt.Println("ok")
	myMap := make(map[string]bool)

	myMap["IP:123"] = true
	myMap["IP:666"] = false
	myMap["IP:Dildo"] = true

	
	
	fmt.Print("len function on the map:  ")
	fmt.Println(len(myMap))
	lengthOfMap := len(myMap)
	for i := 0 ; i < lengthOfMap; i++ {
		fmt.Println(i)

	if "IP:123" == myMap["IP:123"]{
		fmt.Println("Correct")
	
	}
	}
}
