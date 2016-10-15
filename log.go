package main

import (
	"fmt"
	"github.com/fatih/color"
)

var good, bad string

func init() {
	good = color.GreenString("<:)")
	bad = color.RedString("<:( ERROR!")
}

func badDevsInfo(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf("%v %s", good, format)
	return fmt.Printf(s, a...)
}

func badDevsError(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf("%v %s", bad, format)
	return fmt.Printf(s, a...)
}
