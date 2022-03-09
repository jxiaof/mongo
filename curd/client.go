/*
 * @Descripttion:
 * @version:
 * @Author: hujianghong
 * @Date: 2022-03-09 11:50:20
 * @LastEditTime: 2022-03-09 22:38:05
 */
package curd

import (
	"context"

	"mongo/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCli() (*mongo.Client, error) {
	var err error
	var cli *mongo.Client
	cli, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MONGO_URI))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
