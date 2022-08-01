/*******
* @Author:qingmeng
* @Description:
* @File:main
* @Date2022/7/14
 */

package main

import (
	"user-center/server/api"
	"user-center/server/conf"
)

func main() {
	conf.Init()
	api.UserServer(conf.UserTcpPort)
}
