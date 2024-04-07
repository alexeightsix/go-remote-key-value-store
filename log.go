package main

import (
	"fmt"
	"time"
)

func Log(str string) {
	fmt.Println(time.Now().String() + "\t" + str)
}
