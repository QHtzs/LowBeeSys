// 生成xml文件
//

package generateXML

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/QHtzs/appCallSys/Program_Handle_System/mytypes"
)

func readLog(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return "FAILED TO LOAD RUNNING LOG"
	}
	defer file.Close()
	msg, err := ioutil.ReadAll(file)
	if err != nil {
		return "FAILED TO LOAD RUNNING LOG"
	}
	return string(msg[:])
}

//生成程序运行状态页面
func GenerateStatusPage(mp map[string]mytypes.RunStatus, slice []mytypes.RunStatus) string {
	page := `
	<html>
	<head><meta http-equiv="content-type" content="text/html;charset=utf-8"><title>程序调度运行状态</title></head>
	<body>
	<table  align="center" border="1"><caption>ProjectName</caption> 
	<tr width="80%">
	<th width="14%">ProjectName</th> 
	<th width="14%">Description</th> 
	<th width="14%">StartTime</th> 
	<th width="14%">EndTime</th> 
	<th width="14%">ExitStatus</th> 
	<th width="14%">Run Detail</th> 
	</tr>`
	for _, v := range mp {
		page = fmt.Sprintf(`%[1]s <tr width=\"80%\">
		<td width="14%"> %[2]s </td>
		<td width="14%"> %[3]s </td>
		<td width="14%"> %[4]s </td>
		<td width="14%"> %[5]s </td>
		<td width="14%"> %[6]s </td>
		<td width="14%"><strong>stdout</strong></td>
		</tr>`,
			page, v.ProjectName, v.Description, v.StartTime, v.EndTime, v.ExitStatus)
	}
	for _, v := range slice {
		page = fmt.Sprintf(`%[1]s <tr width=\"80%\">
		<td width="14%"> %[2]s </td>
		<td width="14%"> %[3]s </td>
		<td width="14%"> %[4]s </td>
		<td width="14%"> %[5]s </td>
		<td width="14%"> %[6]s </td>
		<td width="14%"><a href="/rundetail?logfile=%[7]s">stdout</a></td>
		</tr>`,
			page, v.ProjectName, v.Description, v.StartTime, v.EndTime, v.ExitStatus, v.LogFile)
	}
	page += `</table></body></html>`
	return page
}

//读取日志文件，并以页面展示
func GenerateLogPage(file string) string {
	log := readLog(file)
	page := fmt.Sprintf(`
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html;charset=utf-8">
	<title>Run Informatin Detail</title>
	</head>
	<body>
	%[1]s
	</body>
	</html>
	`, log)
	return page
}

//配置文件内容展示在页面，便于通过点击启动程序
func GenerateViewPage(funName string, param_ string, mp map[string]mytypes.ProJectConf) string {
	var url string
	url = fmt.Sprintf("/%[1]s?%[2]s=", funName, param_)
	page := `
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html;charset=utf-8">
	<title>Project List</title>
	</head>
	<body>
	<table  align="center" border="1"><caption>Project</caption>
	<tr width="80%">
		<th width="50%">ProjectName</th>
		<th width="50%">Descript</th>
	</tr>`
	length := len(mp)
	tmp := make([]string, length, length)
	i := 0
	for projName, _ := range mp {
		tmp[i] = projName
		i += 1
	}
	sort.Strings(tmp)
	for _, projName := range tmp {
		conf := mp[projName]
		page = fmt.Sprintf("%[1]s<tr width=\"80%\"><td width=\"50%\"><a href=\"%[2]s%[3]s\"> %[3]s </a></td><td width=\"50%\"> %[4]s </td></tr>",
			page, url, projName, conf.Description)
	}
	page += `</table></body></html>`
	return page
}
