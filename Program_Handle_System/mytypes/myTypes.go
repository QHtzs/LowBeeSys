// 项目的可能用于多个包的数据类型

package mytypes

import "encoding/xml"

type RunStatus struct {
	ProjectName string
	Description string
	StartTime   string
	EndTime     string
	ExitStatus  string
	LogFile     string
}

type ProJectConf struct {
	XMLName     xml.Name `xml:"server"`
	ProjectName string   `xml:"ProjectName"`
	Execute     string   `xml:"Execute"`
	FilePath    string   `xml:"FilePath"`
	Description string   `xml:"Description"` //描述
	Schedule    string   `xml:"Schedule"`    //调度
}
