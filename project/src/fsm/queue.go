package fsm

import (
	//"time"
	//"../elevatorHW"
	"fmt"
	"sort"
)

var localQueue [][]int
var localQueueInside []int
var localQueueUp []int
var localQueueDown []int

// localQueue[{insideOrders},{UpOrders},{DownOrders}]

func CreateQueueSlice() {
	localQueue = append(localQueue, localQueueInside)
	localQueue = append(localQueue, localQueueUp)
	localQueue = append(localQueue, localQueueDown)
}

func AppendUpOrder(upOrder int) {
	if upOrder == 1 || upOrder == 2 || upOrder == 3 {
		if !numInSlice(upOrder, localQueue[1]) {
			localQueue[1] = append(localQueue[1], upOrder)
		}
	}
}

func AppendDownOrder(downOrder int) {
	if downOrder == 2 || downOrder == 3 || downOrder == 4 {
		if !numInSlice(downOrder, localQueue[2]) {
			localQueue[2] = append(localQueue[2], downOrder)
		}
	}
}

func AppendInsideOrder(insideOrder int) {
	if insideOrder == 1 || insideOrder == 2 || insideOrder == 3 || insideOrder == 4 {
		if !numInSlice(insideOrder, localQueue[0]) {
			localQueue[0] = append(localQueue[0], insideOrder)
		}
	}
}

func PrintLocalQueue() { //Debugging function
	fmt.Println(localQueue)
}

func numInSlice(num int, slice []int) bool {
	for i := range slice {
		if slice[i] == num {
			return true
		} else {
			continue
		}
	}
	return false

}

func SortLocalQueue() {
	sort.Ints(localQueue[0])
	sort.Ints(localQueue[1])
	sort.Ints(localQueue[2])
}

func DeleteOldestOrderDown() {
	length := len(localQueue[2])
	if length < 1 {
		return
	}
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[2][i], localQueue[2][opp] = localQueue[2][opp], localQueue[2][i]
	}
	localQueue[2] = localQueue[2][:length-1]
	length = len(localQueue[2])
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[2][i], localQueue[2][opp] = localQueue[2][opp], localQueue[2][i]
	}
}

func DeleteOldestOrderUp() {
	length := len(localQueue[1])
	if length < 1 {
		return //Empty Queue
	}
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[1][i], localQueue[1][opp] = localQueue[1][opp], localQueue[1][i]
	}
	localQueue[1] = localQueue[1][:length-1]
	length = len(localQueue[1])
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[1][i], localQueue[1][opp] = localQueue[1][opp], localQueue[1][i]
	}
}

func DeleteOldestOrderInside() {
	length := len(localQueue[0])
	if length < 1 {
		return
	}
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[0][i], localQueue[0][opp] = localQueue[0][opp], localQueue[0][i]
	}
	localQueue[0] = localQueue[0][:length-1]
	length = len(localQueue[0])
	for i := length/2 - 1; i >= 0; i-- {
		opp := length - 1 - i
		localQueue[0][i], localQueue[0][opp] = localQueue[0][opp], localQueue[0][i]
	}
}
