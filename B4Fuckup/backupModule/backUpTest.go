package main
import "fmt"
import "os/exec"
import "time"
import "os"
import "io/ioutil"
import "strconv"
func main() {
	backup := false
	if len(os.Args) == 2 {
		if os.Args[1:][0] == "backup" {
			backup = true
			fmt.Println("Backup started")
		}
	}
	if !backup {
		fmt.Println("Main started")
    str := "go run backUpTest.go"
    cmd1:= exec.Command("ssh student@129.241.187.159", str)
		cmd := exec.Command("bash", "-c", "gnome-terminal -x go run backUpTest.go backup")
		out, err := cmd.Output()
		fmt.Println(string(out))
		check(err)
	}
	i := 0
	for {
		if backup {
			info, err := os.Stat("localQueue")
			check(err)
			if time.Since(info.ModTime()).Seconds() > 3 {
				data, err := ioutil.ReadFile("")
				check(err)
				err = ioutil.WriteFile("", data, 0644)
				i, _ = strconv.Atoi(string(data))
				backup = false
				cmd := exec.Command("bash", "-c", "gnome-terminal -x go run backUpTest.go backup")
				_, err = cmd.Output()
				check(err)
			}
		} else {
			fmt.Printf("%d\n", i)
			err := ioutil.WriteFile("", []byte(strconv.Itoa(i)), 0644)
			check(err)
			i++
		}
		time.Sleep(1*time.Second)
	}
}
func check(e error) {
    if e != nil {
        panic(e)
    }
}
