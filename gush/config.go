package gush

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"os"
)

var Logger seelog.LoggerInterface
var config *GushConfig

type GushConfig struct {
	Port_tcp      string
	Port_notify   string
	Read_timeout  int
	Write_timeout int
}

func init() {
	initLogger()
	initGushConf()
}

func initLogger() {
	log, err := seelog.LoggerFromConfigAsFile("conf/seelog.xml")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	Logger = log
}

func initGushConf() {
	file, fileerr := os.Open("conf/app.json")

	if fileerr != nil {
		fmt.Println("error:", fileerr)
	}

	decoder := json.NewDecoder(file)
	config = &GushConfig{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("error:", err)
	}
}
