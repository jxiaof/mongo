/*
 * @Descripttion:
 * @version:
 * @Author: hujianghong
 * @Date: 2022-03-09 11:50:20
 * @LastEditTime: 2022-03-09 22:46:15
 */
package curd

import (
	"mongo/conf"

	"go.mongodb.org/mongo-driver/mongo"
)

type SDao struct {
	Cli *mongo.Client
	A   *mongo.Collection
	B   *mongo.Collection
	C   *mongo.Collection
}

// NewCurdDao 初始化
func NewCurdDao() *SDao {
	var err error
	dao := &SDao{}
	dao.Cli, err = GetCli()
	if err != nil {
		panic(err)
	}
	dao.A = dao.Cli.Database(conf.CollA).Collection(conf.CollA)
	dao.B = dao.Cli.Database(conf.CollB).Collection(conf.CollB)
	dao.C = dao.Cli.Database(conf.CollC).Collection(conf.CollC)
	return dao
}
