package backup

import (
  "bufio"
  "os"
  "strconv"
  "fmt"
  //"time"
)

func FileBackup(backupCh chan int,backupRemoveOrderCh chan int){
  backupQueue := [4]int{-1,-1,-1,-1}
  for{
    select{
    case newBackupItem := <- backupCh:
      switch newBackupItem{
      case 1:
        if backupQueue[0]==-1{
          backupQueue[0] = 1
        }
      case 2:
        if backupQueue[1]==-1{
          backupQueue[1] = 2
        }
      case 3:
        if backupQueue[2]==-1{
          backupQueue[2] = 3
        }
      case 4:
        if backupQueue[3]==-1{
          backupQueue[3] = 4
        }
      }

      case removeBackupItem := <- backupRemoveOrderCh:
        switch removeBackupItem{
        case 1:
          if backupQueue[0]==1{
            backupQueue[0] = -1
          }
        case 2:
          if backupQueue[1]==2{
            backupQueue[1] = -1
          }
        case 3:
          if backupQueue[2]==3{
            backupQueue[2] = -1
          }
        case 4:
          if backupQueue[3]==4{
            backupQueue[3] = -1
          }
        }
    }
    writeBackupQueueToFile(backupQueue)
    fmt.Println(backupQueue)
  }
}

func checkError(e error) {
    if e != nil {
        panic(e)
    }
}
func writeBackupQueueToFile(localInsideQueue [4]int){
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

func ReadBackupFromFile()[]int{
  var queueToReturn []int
	var lines []string
	var queueElement int
	//var queueElement []byte

  f,err := os.Open("queueBackup.txt")
	if err != nil{
		f,err = os.Create("queueBackup.txt")
    emptySlice := make([]int, 0)
    return emptySlice
	}

 scanner := bufio.NewScanner(f)
 for scanner.Scan() {
	 lines = append(lines, scanner.Text())
 }
 for i := range(lines){
 	 queueElement, err = strconv.Atoi(lines[i])
	 checkError(err)
   if queueElement != -1{
	    queueToReturn = append(queueToReturn, queueElement)
    }
 }

 return queueToReturn
}
