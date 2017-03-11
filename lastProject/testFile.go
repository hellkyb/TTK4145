package main

import (
	"fmt"
	"time"
)
type Order struct {
	Floor  int
	Button int //0 is callDown, 1 callUp, 2 callInside
}

func HandleTimeOutOrder(hallButtonsMap map[Order]int64){
	timeNow := time.Now().Unix()
	if len(hallButtonsMap) > 0{
		for _,value := range hallButtonsMap{
			if timeNow > value+1 {
				fmt.Println("Emergency Order handling")
			}
		}
	}
}



func main() {
	orderMap := make(map[Order]int64)
	currenTime := time.Now().Unix()
	anOrder := Order{2,1}
	orderMap[anOrder] = currenTime
	fmt.Print(orderMap)
	time.Sleep(4*time.Second)
	HandleTimeOutOrder(orderMap)



}
