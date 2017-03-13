package main

import(
	"io"
	"log"
	"net"
	//"time"
	"os"

)

func main(){
	conn, err := net.Dial("tcp", "129.241.187.43:34933")
	if err != nil {
		log.Fatal(err)	
	}
	defer conn.Close()
	mustCopy(os.Stdout, conn)
}

func mustCopy(dst io.Writer, src io.Reader){
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}	
}
