package fsm

import(
	"../elevatorHW"
	"time"
)

type queue struct{
	queueMatrix[elevatorHW.NFloors][elevatorHW.NButtons]orderStatus
}

type orderStatus struct{
	active bool
	assignedElevator string 
}

var localQueue queue 
var globalQueue queue 


