package elevatorHW

//in port 4
const port4Subdevice = 3
const port4ChannelOffset = 16

//const port4Direction		= COMEDI_INPUT
const obstruction = (0x300 + 23)
const stop = (0x300 + 22)
const buttonCommand1 = (0x300 + 21)
const buttonCommand2 = (0x300 + 20)
const buttonCommand3 = (0x300 + 19)
const buttonCommand4 = (0x300 + 18)
const buttonUp1 = (0x300 + 17)
const buttonUp2 = (0x300 + 16)

//in port 1
const port1Subdevice = 2
const port1ChannelOffset = 0

//const port1Direction		= COMEDI_INPUT
const buttonDown2 = (0x200 + 0)
const buttonUp3 = (0x200 + 1)
const buttonDown3 = (0x200 + 2)
const buttonDown4 = (0x200 + 3)
const sensorFloor1 = (0x200 + 4)
const sensorFloor2 = (0x200 + 5)
const sensorFloor3 = (0x200 + 6)
const sensorFloor4 = (0x200 + 7)

//out port 3
const port3Subdevice = 3
const port3ChannelOffset = 8

//const port3Direction		= COMEDI_OUTPUT
const Motordir = (0x300 + 15)
const lightStop = (0x300 + 14)
const lightCommand1 = (0x300 + 13)
const lightCommand2 = (0x300 + 12)
const lightCommand3 = (0x300 + 11)
const lightCommand4 = (0x300 + 10)
const lightUp1 = (0x300 + 9)
const lightUp2 = (0x300 + 8)

//out port 2
const port2Subdevice = 3
const port2ChannelOffset = 0

//const port2Direction		= COMEDI_OUTPUT
const lightDown2 = (0x300 + 7)
const lightUp3 = (0x300 + 6)
const lightDown3 = (0x300 + 5)
const lightDown4 = (0x300 + 4)
const lightDoorOpen = (0x300 + 3)
const lightFloorInd2 = (0x300 + 1)
const lightFloorInd1 = (0x300 + 0)

//out port 0
const motor = (0x100 + 0)

//non-existing ports = (for alignment)
const buttonDown1 = -1
const buttonUp4 = -1
const lightDown1 = -1
const lightUp4 = -1
