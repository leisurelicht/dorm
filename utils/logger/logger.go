package logger

import (
	"github.com/chuck1024/gd"
	"os"
	"runtime"
	"strings"
)

var dev bool

func init() {
	idc := os.Getenv("IDC")
	if idc == "dev" || idc == "test" {
		dev = true
	}
}

func getCallerName() (string, string, string) {
	pc, _, _, _ := runtime.Caller(2)
	ca := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	cd := strings.Split(ca[len(ca)-1], ".")
	return cd[0], cd[1], cd[2]
}

func PrintCallerInfo() {
	if !dev {
		return
	}

	_, structName, callName := getCallerName()
	gd.Info("[%s] Call %s", structName, callName)
}
