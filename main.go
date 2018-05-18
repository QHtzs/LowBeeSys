// 程序运行管理系统
package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/QHtzs/appCallSys/Program_Handle_System/Parse_Conf"
	"github.com/QHtzs/appCallSys/Program_Handle_System/generateXML"
	"github.com/QHtzs/appCallSys/Program_Handle_System/myLog"
	"github.com/QHtzs/appCallSys/Program_Handle_System/mytypes"
	"github.com/QHtzs/appCallSys/Program_Handle_System/pSchedules"
)

type MyHandler struct {
	schedule pSchedules.ShceduleData
}

func init_handle(confFileName string, maxApp int) MyHandler {
	ret := MyHandler{}
	conf := Parse_Conf.ConfigMap{}
	conf.Initialization(confFileName)
	ret.schedule.Confs = conf.ConfMap
	ret.schedule.MaxRunningTaskLimit = maxApp
	ret.schedule.IsRunning = make([]string, 0, maxApp+10)
	ret.schedule.IsFinished = make([]mytypes.RunStatus, 0, 2*maxApp+10)
	ret.schedule.ReadyRan = make(map[string]mytypes.RunStatus, maxApp+10)
	ret.schedule.CanAdd = true
	return ret
}

func (handle *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ruri := r.URL.RequestURI()
	if "/favicon.ico" != ruri {
		myLog.ToLogFile(ruri)
	}
	if "/list" == ruri {
		page := generateXML.GenerateViewPage("callApp", "appName", handle.schedule.Confs)
		io.WriteString(w, page)
	} else if "/favicon.ico" == ruri {
		return
	} else if strings.Index(ruri, "rundetail") == 1 {
		r.ParseForm()
		logFile := r.Form["logfile"]
		if len(logFile) == 0 {
			io.WriteString(w, "请求无效,程序出现未知错误")
		} else {
			page := generateXML.GenerateLogPage(logFile[0])
			io.WriteString(w, page)
		}
		return

	} else if "/run" == ruri {
		page := generateXML.GenerateStatusPage(handle.schedule.ReadyRan, handle.schedule.IsFinished)
		io.WriteString(w, page)
	} else if strings.Index(ruri, "callApp") == 1 {
		r.ParseForm()
		appNames := r.Form["appName"]
		if len(appNames) == 0 {
			io.WriteString(w, "请求错误，没有重要参数")
			return
		}
		appName := appNames[0]
		if appName == "" {
			io.WriteString(w, "请求错误，所要调用的程序名为空")
		}
		path := handle.schedule.Confs[appName].FilePath
		commandName := handle.schedule.Confs[appName].Execute
		if path == "" && commandName == "" {
			io.WriteString(w, "配置文件填写错误，请务必填写EXECUTE,FILEPATH项")
			return
		} else if commandName == "" {
			io.WriteString(w, "配置文件填写错误，请指定运行程序所需的 指令/程序")
			return
		}
		info := handle.schedule.CallApp(appName)
		io.WriteString(w, info)
	} else {
		msg := `
		<html>
		<head>
		<meta http-equiv="content-type" content="text/html;charset=utf-8">
		<title>404</title>
		<style>
            *{
                margin:0;
                padding:0;
            }
            #div0{
            width: 400px;
            height: 400px;
        	}
			div{
            width: 400px;
            height: 20px;
        	}
        	.center-in-center{
            	position: absolute;
            	top: 50%;
            	left: 50%;
        	}
  
        </style>
		</head>
		<body>
		<div class="center-in-center" id="div0">
		<div>
		<span>
		<a data-click="{
			'F':'778317EA',
			'F1':'9D73F1E4',
			'F2':'4CA6DE6B',
			'F3':'54E5343F',
			'T':'1526612023',
			'y':'E7C7DB9B'
			 }" href="/list">
			可调度程序列表
		</a>
		</span>
		</div>
		<div>
		<span>
		<a data-click="{
			'F':'778317EA',
			'F1':'9D73F1E4',
			'F2':'4CA6DE6B',
			'F3':'54E5343F',
			'T':'1526612023',
			'y':'E7C7DB9B'
			 }" href="/run">
			程序运行状态信息							
		</a>
		</span>
		</div>
		</div>
		</body>
		</html>
		`
		io.WriteString(w, msg)
	}

}

func main() {
	handle := init_handle("conf.xml", 100)
	mux := http.NewServeMux()
	mux.Handle("/", &handle)
	http.ListenAndServe(":13145", mux)
}
