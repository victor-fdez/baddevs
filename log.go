package main

import (
	"fmt"
	"github.com/fatih/color"
)

var good, bad, debug string

func init() {
	good = color.GreenString("<:)")
	bad = color.RedString("<:( ERROR!")
	debug = color.YellowString("<:| DEBUG ")
}

func badDevsInfo(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf("%v %s", good, format)
	return fmt.Printf(s, a...)
}

func badDevsDebug(format string, a ...interface{}) (n int, err error) {
	if badDevsConfig.verbose {
		s := fmt.Sprintf("%v %s", debug, format)
		return fmt.Printf(s, a...)
	}
	return 0, nil
}

func badDevsError(format string, a ...interface{}) (n int, err error) {
	s := fmt.Sprintf("%v %s", bad, format)
	return fmt.Printf(s, a...)
}
