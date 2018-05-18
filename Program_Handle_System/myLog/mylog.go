// 简单封装log
package myLog

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
)

func ToLogFile(info interface{}) {
	TaskLog("logs.log", info)
}

func TaskLog(file string, info interface{}) {
	var msg string
	if reflect.TypeOf(info) == reflect.TypeOf("1") {
		msg, _ = info.(string)
	} else if reflect.TypeOf(info) == reflect.TypeOf(errors.New("")) {
		msg = info.(error).Error()
	} else {
		msg = "unkown msg type"
	}
	if runtime.GOOS == "windows" {
		msg += "\r\n"
	} else {
		msg += "\t\n"
	}

	_file, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		fmt.Println(err)
	}
	defer _file.Close()
	logger := log.New(_file, "", log.LstdFlags|log.Llongfile)
	logger.Println(msg)

}
