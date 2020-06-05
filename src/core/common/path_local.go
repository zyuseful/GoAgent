package common

import "os"

//获取当前路径
func GetLocalPath() string {
	dir,_ := os.Getwd()
	return dir
}