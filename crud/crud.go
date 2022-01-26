/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-20 11:25:35
 * @LastEditTime: 2022-01-26 23:34:01
 */
package crud

// 思路: bson在golang中编写比较繁琐,我的思路是将json转换为bson,贴近原生的mongodb操作,只要会写mongo查询语句就能很快很方便的操作
// 功能: 提供对于mongo的简单增删改查操作
// 注意: Mutil是一个操作类型,可以是insert,update,delete一个或者多个document
// 注意: 针对mongo-driver操作做了简化,只保留经常使用的mongo操作,如果需要更复杂的操作需求,请使用mongo-driver提供的的接口
// TODO: 1): 增加聚合查询 (aggregate)

import (
	"context"
	"fmt"

	"mongo/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// 列表获取接口的最大返回条数
	ReturnMaxDocs int64 = 100
)

// 提供了一个简单的操作mongodb的接口,检查document是否存在
func (coll *Scoll) CheckExist(colName, filter string) (bool, error) {
	// 存在返回true,不存在返回false, 出错返回error
	var err error
	var res bson.M
	if colName == "" || filter == "" {
		return false, fmt.Errorf("collection name or filter is empty")
	}
	opts := options.FindOneOptions{}
	fil, err := Json2Bson(filter)
	if err != nil {
		return false, err
	}
	err = coll.ColMap[colName].FindOne(context.TODO(), fil, &opts).Decode(&res)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// 提供一个简单的列表数据获取接口
func (coll *Scoll) List(colName string, p *SMongoListParams) (*SMongoListRes, error) {
	var err error
	var res SMongoListRes
	res.Count = 0
	res.Data = make([]bson.M, 0)
	if colName == "" {
		return nil, fmt.Errorf("collection name is empty")
	}
	opts := options.FindOptions{}
	if p.Filter == "" {
		return nil, fmt.Errorf("filter is empty")
	}
	if p.Options.Skip > 0 {
		opts.Skip = &p.Options.Skip
	}
	if p.Options.Limit > 0 {
		// 设置最大返回条数为100
		if p.Options.Limit > ReturnMaxDocs {
			log.Log.Warnf("return doc limit is too large, set to %v", ReturnMaxDocs)
			p.Options.Limit = ReturnMaxDocs
		}
		opts.Limit = &p.Options.Limit
	}
	if p.Options.Sorts != "" {
		sort, err := Json2Bson(p.Options.Sorts)
		if err != nil {
			err = fmt.Errorf("parser sorts json to bson error: %v", err)
			return nil, err
		}
		opts.Sort = sort
	}
	if p.Options.Fields != "" {
		fields, err := Json2Bson(p.Options.Fields)
		if err != nil {
			err = fmt.Errorf("parser fields json to bson error: %v", err)
			return nil, err
		}
		opts.Projection = fields
	}
	// 解析过滤条件
	fil, err := Json2Bson(p.Filter)
	if err != nil {
		err = fmt.Errorf("parser filter json to bson error: %v", err)
		return nil, err
	}
	collection := coll.ColMap[colName]
	cur, err := collection.Find(context.TODO(), &fil, &opts)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, nil
	}
	defer cur.Close(context.TODO())
	if err := cur.All(context.TODO(), &res.Data); err != nil {
		return nil, err
	}
	// 获取数据总数
	count, err := collection.CountDocuments(context.Background(), &fil)
	if err != nil {
		return nil, err
	}
	res.Count = count
	return &res, nil
}

// 获取一个document
func (coll *Scoll) Get(colName string, p *SMongoGetParams) (*primitive.M, error) {
	var err error
	var res bson.M
	if colName == "" {
		return nil, fmt.Errorf("collection name is empty")
	}

	if p.Filter == "" {
		return nil, fmt.Errorf("filter is empty")
	}
	opts := options.FindOneOptions{}
	if p.Options.Skip > 0 {
		opts.Skip = &p.Options.Skip
	}
	if p.Options.Sorts != "" {
		sort, err := Json2Bson(p.Options.Sorts)
		if err != nil {
			err = fmt.Errorf("parser sorts json to bson error: %v", err)
			return nil, err
		}
		opts.Sort = sort
	}
	if p.Options.Fields != "" {
		fields, err := Json2Bson(p.Options.Fields)
		if err != nil {
			err = fmt.Errorf("parser fields json to bson error: %v", err)
			return nil, err
		}
		opts.Projection = fields
	}
	// 解析过滤条件
	fil, err := Json2Bson(p.Filter)
	if err != nil {
		err = fmt.Errorf("json to bson error: %v", err)
		return nil, err
	}
	collection := coll.ColMap[colName]
	err = collection.FindOne(context.TODO(), fil, &opts).Decode(&res)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, nil
	}
	return &res, nil
}

// 创建一个或者多个document
func (coll *Scoll) Create(colName string, doc []interface{}) ([]interface{}, error) {
	if len(doc) == 0 {
		return nil, fmt.Errorf("doc is empty")
	}
	collection := coll.ColMap[colName]
	if len(doc) == 1 {
		res, err := collection.InsertOne(context.TODO(), doc[0])
		if err != nil {
			return nil, err
		}
		if res.InsertedID == nil {
			return nil, fmt.Errorf("return insert id is empty")
		}

		log.Log.Infof("inserted one document with ID %v\n", res.InsertedID)
		return []interface{}{res.InsertedID}, nil
	} else {
		res, err := collection.InsertMany(context.TODO(), doc)
		if err != nil {
			return nil, err
		}
		if res.InsertedIDs == nil {
			return nil, fmt.Errorf("return insert id is empty")
		}
		log.Log.Infof("inserted many documents with IDs %v\n", res.InsertedIDs)
		return res.InsertedIDs, nil
	}
}

// 更新一个或者多个document
func (coll *Scoll) Update(colName string, p *SMongoUpdateParams) (*mongo.UpdateResult, error) {
	var err error
	var res *mongo.UpdateResult
	collection := coll.ColMap[colName]
	if p.Filter == "" {
		return nil, fmt.Errorf("filter is empty")
	}
	if p.Update == "" {
		return nil, fmt.Errorf("update is empty")
	}
	// 解析过滤条件
	fil, err := Json2Bson(p.Filter)
	if err != nil {
		err = fmt.Errorf("parser filter json to bson error: %v", err)
		return nil, err
	}
	// 解析更新条件
	upd, err := Json2Bson(p.Update)
	if err != nil {
		err = fmt.Errorf("parser update json to bson error: %v", err)
		return nil, err
	}

	opts := options.UpdateOptions{}
	if p.Options.Upsert {
		opts.SetUpsert(true)
	}
	if p.Options.Multi {
		res, err = collection.UpdateMany(context.TODO(), fil, upd, &opts)
	} else {
		res, err = collection.UpdateOne(context.TODO(), fil, upd, &opts)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 删除一个或者多个document
func (coll *Scoll) Delete(colName string, p *SMongoDeleteParams) (*mongo.DeleteResult, error) {
	var err error
	var res *mongo.DeleteResult
	collection := coll.ColMap[colName]
	if p.Filter == "" {
		return nil, fmt.Errorf("filter is empty")
	}
	// 解析过滤条件
	fil, err := Json2Bson(p.Filter)
	if err != nil {
		err = fmt.Errorf("parser filter json to bson error: %v", err)
		return nil, err
	}

	opts := options.DeleteOptions{}
	if p.Options.Multi {
		res, err = collection.DeleteMany(context.TODO(), fil, &opts)
	} else {
		res, err = collection.DeleteOne(context.TODO(), fil, &opts)
	}
	if err != nil {
		return nil, err
	}
	log.Log.Infof("deleted %v documents\n", res.DeletedCount)
	return res, nil
}
