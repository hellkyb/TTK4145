package elevatorHW

import(
	"driver"
	"time"	
)

func SetDoorOpen(){
	driver.ClearBit(lightDoorOpen)
	time.Sleep(3*time.Second)	
	driver.SetBit(lightDoorOpen)
	time.Sleep(3*time.Second)
	driver.ClearBit(lightDoorOpen)
}
