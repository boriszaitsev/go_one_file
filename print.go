package main

import "fmt"
import "os"

//-------------------------------------------------------------------------
var print = fmt.Println
var printf = fmt.Printf
var sprintf = fmt.Sprintf
var err = func(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
var errf = func(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s, a...)
}
