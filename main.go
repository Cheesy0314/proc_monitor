package main

import (
	"time"
	"io/ioutil"
	"strings"
)

type process struct {
	name string
}

type processList = []process


func main() {
	valid := processList{}
	path := "//Users/dylan/GoglandProjects/proctologist/processors.txt"
    data,readError := ioutil.ReadFile(path)
    if readError != nil {
    	panic(readError)
	}

	arr := strings.Split(string(data),"\n")
	for _,proc := range arr {
		if proc != "" {
			nextProc := process{name: proc}
			valid = append(valid,nextProc)
		}
	}

	for {
		problem := report(valid)
		if problem != "" {
			sendAlert(problem)
		}

		time.Sleep(time.Second)
	}
}
