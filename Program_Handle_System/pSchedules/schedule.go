// 调度管理

package pSchedules

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/QHtzs/appCallSys/Program_Handle_System/myLog"
	"github.com/QHtzs/appCallSys/Program_Handle_System/mytypes"
)

type ShceduleData struct {
	Confs               map[string]mytypes.ProJectConf
	Lock                sync.Mutex
	MaxRunningTaskLimit int                          //限制最大运行任务数
	CanAdd              bool                         //是否可以添加新任务
	IsRunning           []string                     //正在运行
	IsFinished          []mytypes.RunStatus          //运行完毕程序详情
	ReadyRan            map[string]mytypes.RunStatus //正在运行程序详情
}

func empty_N(slice []mytypes.RunStatus, N int) []mytypes.RunStatus {
	log := "\t\n"
	for i, v := range slice {
		if i >= N-1 {
			break
		}
		log = fmt.Sprintf(`%[1]s \t\n ProjectName:%[2]s; 
		Description: %[3]s;
		StartTime:%[4]s;
		EndTime: %[5]s;
		ExitStatus: %[6]s;
		LogFile:%[7]s; `,
			log, v.ProjectName, v.Description, v.StartTime, v.EndTime, v.ExitStatus, v.LogFile)
	}
	myLog.TaskLog("finish_task_log", log)
	length := len(slice)
	return append(slice[N:length-2], slice[length-2:]...) //a new slice copy
}

func callApp(commandName string, path string, logFile string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("unhandled error raised while execute command.")
		}
	}()
	cmd := exec.Command(commandName, path)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		return errors.New("Failed to create log file:" + logFile)
	}
	cmd.Stdout = file // 输出定向到file文件
	cmd.Start()
	err = cmd.Wait()
	file.Close() //关闭文件，不然无法写入
	return err
}

func (scd *ShceduleData) add(appName string) (string, error) {
	scd.Lock.Lock()
	defer scd.Lock.Unlock()
	if !scd.CanAdd {
		return "", errors.New("Running Task reach Limit,Please wait!")
	}
	if len(scd.IsRunning) >= scd.MaxRunningTaskLimit {
		scd.CanAdd = false
		return "", errors.New("Running Task reach Limit,Please wait!")
	}
	scd.IsRunning = append(scd.IsRunning, appName)
	runStatus := mytypes.RunStatus{}
	runStatus.ProjectName = appName
	runStatus.StartTime = time.Now().Format("2006-01-02 15:04:05")
	runStatus.EndTime = "running..."
	runStatus.Description = scd.Confs[appName].Description
	runStatus.ExitStatus = "running..."
	runStatus.LogFile = fmt.Sprintf("%[1]s_t%[2]d.log", appName, time.Now().Unix())
	scd.ReadyRan[appName] = runStatus
	return runStatus.LogFile, nil
}

func (scd *ShceduleData) remove(appName string, exitStatus string) {
	scd.Lock.Lock()
	defer scd.Lock.Unlock()
	var swapSlice []string
	for i, v := range scd.IsRunning {
		if v == appName {
			swapSlice = append(scd.IsRunning[:i], scd.IsRunning[i+1:]...) //len -1
		}
	}
	scd.IsRunning = swapSlice //从正在运行中删除
	//添加进运行完毕
	runStatus := scd.ReadyRan[appName]
	runStatus.EndTime = time.Now().Format("2006-01-02 15:04:05")
	runStatus.ExitStatus = exitStatus
	scd.IsFinished = append(scd.IsFinished, runStatus)
	//从正在运行详情中删除
	delete(scd.ReadyRan, appName)
	scd.CanAdd = true
	if len(scd.IsFinished) >= 2*scd.MaxRunningTaskLimit {
		isFinish := empty_N(scd.IsFinished, scd.MaxRunningTaskLimit)
		scd.IsFinished = isFinish
	}

}

func (scd *ShceduleData) include(appName string) bool {
	var res bool = false
	for _, v := range scd.IsRunning {
		if v == appName {
			res = true
			break
		}
	}
	return res
}

func (scd *ShceduleData) CallApp(appName string) string {
	commandName := scd.Confs[appName].Execute
	path := scd.Confs[appName].FilePath
	if scd.include(appName) {
		return "有同名程序正在运行,请待程序运行完后再启动..."
	}
	logFile, err := scd.add(appName)
	if err != nil {
		return "管理系统中正在运行的程序已达上限，请待部分程序完成后再启动..."
	}
	go func() {
		err = callApp(commandName, path, logFile)
		exitStatus := ""
		if err != nil {
			exitStatus = err.Error()
		} else {
			exitStatus = "sucessed"
		}
		scd.remove(appName, exitStatus)
	}()
	return "程序启动成功，运行信息请查看相关页面..."
}
