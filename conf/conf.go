/*******
* @Author:qingmeng
* @Description:
* @File:conf
* @Date2022/7/14
 */

package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
	"user-center/cache"
	"user-center/dao"
)

var (
	UserTcpPort string
	Db          string
	DbHost      string
	DbPort      string
	DbUser      string
	DbPassWord  string
	DbName      string
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
	TokenSecret string
)

func Init() {
	//从本地读取环境
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("ini load failed", err)
	}

	LoadService(file)
	LoadSecret(file)

	//初始化数据库
	LoadMySQL(file) //读取配置信息
	dsn := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	dao.InitDatabase(dsn)

	//初始化redis
	LoadRedis(file)
	cache.InitRedis(RedisDbName, RedisAddr) //redis连接
}

func LoadSecret(file *ini.File) {
	TokenSecret = file.Section("secret").Key("TokenSecret").String()
}

func LoadService(file *ini.File) {
	UserTcpPort = file.Section("service").Key("userTcp").String()
}

func LoadMySQL(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func LoadRedis(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()

}
