package jorm

import (
	"errors"
	"log"
	"runtime"
	"strings"
)

//获取程序执行路径,default skip is 1
//skip:0_调用runtime.Caller的地方，1_调用GetRuntimePath的地方,2_函数执行的路径
func GetRuntimePath(skip int) (string, error) {
	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		log.Println("error")
		return "", errors.New("get parent path fail")
	}
	//log.Println("path = ", file)
	lastIndex := strings.LastIndex(file, "/")
	if lastIndex == -1 {
		return "", errors.New("path not contain substr")
	}
	//log.Println("lastIndex = ", lastIndex)
	path := string([]byte(file)[0:lastIndex])
	log.Println("path = ", path)
	return path, nil
}
