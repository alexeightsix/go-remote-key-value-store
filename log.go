package main

import (
	// "fmt"
	// "time"
)

type Log struct {
}

func (l Log) Log(str string) {
	// fmt.Println(time.Now().String() + "\t" + str)
}

func (l Log) Notice(str string) {
	l.Log(str)
}

func (l Log) Error(str string) {
	l.Log(str)
}

func (l Log) Warning(str string) {
	l.Log(str)
}
