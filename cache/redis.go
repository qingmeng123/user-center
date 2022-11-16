/*******
* @Author:qingmeng
* @Description:
* @File:redis
* @Date2022/7/16
 */

package cache

import (
	"github.com/go-redis/redis"
	"log"
	"strconv"
)

var RedisClient *redis.Client

func InitRedis(redisDbName string, redisAddr string) {
	db, _ := strconv.ParseUint(redisDbName, 10, 64) //string to int
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   int(db),
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Println(err)
		panic(err)
	}
	RedisClient = client
}
