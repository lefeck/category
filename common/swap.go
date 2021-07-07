package common

import "encoding/json"

//通过json tag进行结构体赋值，转化成struct格式
func SwapTo(request,category interface{}) (err error)  {
	dataByte,err :=json.Marshal(request)
	if err !=nil {
		return
	}
	return json.Unmarshal(dataByte,category)
}
