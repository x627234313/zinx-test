package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// 定义全局对象，存储Zinx的全部参数，供其他模块使用
//也可以通过 zinx.json 配置文件 用户自定义
type GlobalObj struct {
	Name          string
	Host          string
	TcpPort       int
	Version       string
	MaxConn       int
	MaxPacketSize int
}

var GlobalObject *GlobalObj

// 读取用户自定义的配置文件
func (g *GlobalObj) Reload() {
	conf, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		fmt.Println("ReadFile zinx.json error: ", err)
		return
	}

	err = json.Unmarshal(conf, GlobalObject)
	if err != nil {
		fmt.Println("Unmarshal zinx.json to GlobalObj error: ", err)
		return
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "Zinx-Test Server",
		Host:          "0.0.0.0",
		TcpPort:       9999,
		MaxConn:       10000,
		MaxPacketSize: 4096,
	}

	GlobalObject.Reload()
}
