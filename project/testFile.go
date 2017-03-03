package main

import ("fmt")
type struktur struct{
	heltall int
	streng string
}

func main() {
	a:= struktur{1,"hei"}
 	fmt.Println(a.heltall)
}
