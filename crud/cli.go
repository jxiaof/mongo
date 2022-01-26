package crud

import (
	"context"
	"mongo/conf"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCli() (*mongo.Client, error) {
	var err error
	var cli *mongo.Client
	cli, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(conf.MONGO_URL))
	if err != nil {
		return nil, err
	}
	return cli, nil
}
