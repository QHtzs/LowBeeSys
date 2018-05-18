// 配置文件解析
package Parse_Conf

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/QHtzs/appCallSys/Program_Handle_System/mytypes"
)

/*
type ProJectConf struct {
	XMLName     xml.Name `xml:"server"`
	ProjectName string   `xml:"ProjectName"`
	Execute     string   `xml:"Execute"`
	FilePath    string   `xml:"FilePath"`
	Description string   `xml:"Description"` //描述
	Schedule    string   `xml:"Schedule"`    //调度
}
*/

type XMLSruct struct {
	XMLName     xml.Name              `xml:"servers"`
	Version     string                `xml:"version,attr"`
	Projects    []mytypes.ProJectConf `xml:"server"`
	Description string                `xml:",innerxml"`
}

type ConfigMap struct {
	ConfMap map[string]mytypes.ProJectConf
}

func (cfmp *ConfigMap) Initialization(path string) error {
	cfmp.ConfMap = make(map[string]mytypes.ProJectConf)
	var key string
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("open file error:%v", err)
		return err
	}
	defer file.Close()

	XmlData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("read error :%v", err)
		return err
	}

	xmlstruct := XMLSruct{}
	err = xml.Unmarshal(XmlData, &xmlstruct)
	if err != nil {
		fmt.Printf("deseries xml error:%v", err)
		return err
	}

	for _, projConf := range xmlstruct.Projects {
		key = projConf.ProjectName
		cfmp.ConfMap[key] = projConf
	}
	return nil

}
