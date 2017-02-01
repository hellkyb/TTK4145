package driver 

/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/

import "C"

func Io_init() bool{
	return bool(int(C.io_init()) != 1)
}

func Io_set_bit(channel int){
	C.io_set_bit(channel)
}