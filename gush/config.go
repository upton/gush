package gush

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
)

var Logger seelog.LoggerInterface

func init() {
	initLogger()
}

func initLogger() {
	log, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	Logger = log
}
