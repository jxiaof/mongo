/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-26 22:48:07
 * @LastEditTime: 2022-01-26 23:10:28
 */
package crud

import (
	"mongo/conf"

	"go.mongodb.org/mongo-driver/mongo"
)

func GetColl() (*Scoll, error) {
	var err error
	var scol Scoll
	scol.ColMap = make(map[string]*mongo.Collection)
	scol.Cli, err = GetCli()
	if err != nil {
		return nil, err
	}
	scol.ColMap[conf.CollA] = scol.Cli.Database(conf.DB_NAME).Collection(conf.CollA)
	scol.ColMap[conf.CollB] = scol.Cli.Database(conf.DB_NAME).Collection(conf.CollB)
	scol.ColMap[conf.CollC] = scol.Cli.Database(conf.DB_NAME).Collection(conf.CollC)
	return &scol, nil
}

// NewCrud 初始化
func NewCrud() *Scoll {
	coll, err := GetColl()
	if err != nil {
		panic(err)
	}
	return coll
}
