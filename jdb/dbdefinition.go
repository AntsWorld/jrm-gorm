package jdb

import (
	"errors"
	"fmt"
)

//数据库连接信息
type DbConnectInfo struct {
	DriverName  string //数据库驱动名称
	User        string //用户名
	Password    string //密码
	Ip          string //IP地址
	Port        int    //端口
	DBName      string //数据库名称
	Description string //描述
}

//根据结构体参数组装dataSourceName
func (connectInfo *DbConnectInfo) ToDataSourceName() (string, error) {
	//校验参数
	if "" == connectInfo.DriverName {
		return "", errors.New("DriverName can't be nil")
	}
	var dataSourceName string
	switch connectInfo.DriverName {
	case "mysql":
		//username:password@protocol(address)/dbname?param=value
		//"user:password@tcp(127.0.0.1:3306)/jrm_manager"
		dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", connectInfo.User, connectInfo.Password,
			connectInfo.Ip, connectInfo.Port, connectInfo.DBName)
	default:
		return "", errors.New("un support driver")
	}

	return dataSourceName, nil
}
