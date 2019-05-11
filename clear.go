package main

import (
	"os"
	"os/exec"
	"runtime"
)

var clearMap map[string]func() //create a map for storing clear funcs

func init() {
	clearMap = make(map[string]func()) //Initialize it
	clearMap["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["darwin"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearMap["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

//clear 清屏
func clear() {
	value, ok := clearMap[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                             //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
	return
}
