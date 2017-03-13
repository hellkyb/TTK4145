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
		cmd := exec.Command("bash", "-c", "gnome-terminal -x go run Ex6.go backup")
		out, err := cmd.Output()
		fmt.Println(string(out))
		check(err)
	}
	i := 0
	for {
		if backup {
			info, err := os.Stat("tall")
			check(err)
			if time.Since(info.ModTime()).Seconds() > 3 {
				data, err := ioutil.ReadFile("tall")
				check(err)
				err = ioutil.WriteFile("tall", data, 0644)
				i, _ = strconv.Atoi(string(data))
				backup = false
				cmd := exec.Command("bash", "-c", "gnome-terminal -x go run Ex6.go backup")
				_, err = cmd.Output()
				check(err)
			}
		} else {
			fmt.Printf("%d\n", i)
			err := ioutil.WriteFile("tall", []byte(strconv.Itoa(i)), 0644)
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