/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-26 22:31:11
 * @LastEditTime: 2022-01-26 23:10:28
 */
package conf

import "os"

const (
// DB_NAME database name ...
// MONGO_URL = "mongodb://localhost:27017"
)

// 获取env
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var (
	// 日志级别
	LOG_LEVEL = getEnv("LOG_LEVEL", "debug")
)

var (
	// mongo url
	MONGO_URL = getEnv("MONGO_URL", "mongodb://localhost:27017")
	// 数据库名称
	DB_NAME = getEnv("DB_NAME", "db_name")
	// 数据库集合
	CollA = getEnv("CollA", "coll_a")
	CollB = getEnv("CollB", "coll_b")
	CollC = getEnv("CollC", "coll_c")
)
