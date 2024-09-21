// 响应工具包 用于在HTTP服务器中向客户端返回JSON格式的响应

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H struct {
	Code  int         // 状态码
	Msg   string      // 信息字符串，描述响应状态或错误信息
	Data  interface{} // 任意类型的数据
	Rows  interface{} // 用于返回分页数据中的当前页的具体内容
	Total interface{} // 用于返回分页数据时的总记录数
}

// Resp 响应函数，向客户端返回状态码、消息、数据的JSON响应
func Resp(w http.ResponseWriter, code int, data interface{}, msg string) {
	w.Header().Set("Content-Type", "application/json") // 设置响应头，表明返回的数据格式为JSON
	w.WriteHeader(http.StatusOK)                       // 设置响应码
	h := H{                                            // 创建结构体实例，将code、data、msg填充到对应字段中
		Code: code,
		Data: data,
		Msg:  msg,
	}
	ret, err := json.Marshal(h) // 将h序列化为JSON字符串
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret) // 将JSON响应谢日到w中，即返回给客户端
}

// RespList 用于返回分页数据的响应函数
func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h := H{ // 创建结构体实例，将code、data、total填充到对应字段中
		Code:  code,
		Rows:  data,
		Total: total,
	}
	ret, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(ret)
}

// RespFail 快捷响应函数，返回失败信息
func RespFail(w http.ResponseWriter, msg string) {
	Resp(w, -1, nil, msg)
}

// RespOK 快捷响应函数，返回成功信息
func RespOK(w http.ResponseWriter, data interface{}, msg string) {
	Resp(w, 0, data, msg)
}

// RespOKList 快捷响应函数，返回分页成功的信息
func RespOKList(w http.ResponseWriter, data interface{}, total interface{}) {
	RespList(w, 0, data, total)
}
