package main

import (
		"bufio"
    "fmt"
    "os"
		"strconv"
)
func checkError(e error) {
    if e != nil {
			fmt.Println(e)
        panic(e)
    }
}
func writeBackupQueueToFile(localInsideQueue []int){
  var queueToBackup string
  for i:=range(localInsideQueue){
    queueToBackup += strconv.Itoa(localInsideQueue[i]) +"\n"
  }
	f, err := os.Create("queueBackup.txt")
	checkError(err)

	defer f.Close()

	_, err = f.WriteString(queueToBackup)
	checkError(err)
	f.Sync()
}

func readBackupFromFile()[]int{
  var queueToReturn []int
	var lines []string
	var queueElement int
	//var queueElement []byte

  f,err := os.Open("queueBackup.txt")
	if err != nil{
		_,err = os.Create("queueBackup.txt")
	}

 scanner := bufio.NewScanner(f)
 for scanner.Scan() {
	 lines = append(lines, scanner.Text())
 }
 for i := range(lines){
 	 queueElement, err = strconv.Atoi(lines[i])
	 checkError(err)

	 queueToReturn = append(queueToReturn, queueElement)
 }
 return queueToReturn
}


func main() {
	localTestQueue := []int{1,2,3,4}
	writeBackupQueueToFile(localTestQueue)
	// localQueue := readBackupFromFile()
	// fmt.Println(localQueue)
	// writeBackupQueueToFile([]int{1,2,3})
	localQueue := readBackupFromFile()
	fmt.Println(localQueue)
}
