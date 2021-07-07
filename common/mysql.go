package common

import "github.com/micro/go-micro/v2/config"

type MysqlConfig struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pwd string `json:"pwd"`
	Database string `json:"database"`
	Port int64 `json:"port"`
}

//获取 mysql 的配置从配置中心
func GetMysqlFromConsul(config config.Config,path ...string) *MysqlConfig{
	mysqlConfig := &MysqlConfig{}
	//从mysql中扫描获取配置信息，返回
	config.Get(path...).Scan(mysqlConfig)
	return mysqlConfig
}