package main

import (
	"github.com/upton/gush/gush"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gush.Run()
}
