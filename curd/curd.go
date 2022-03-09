/*
 * @Descripttion:
 * @version:
 * @Author: hujianghong
 * @Date: 2022-03-08 17:21:38
 * @LastEditTime: 2022-03-09 22:55:00
 */
package curd

// Idea: bson is tedious to write in golang, my idea is to convert json to bson, close to the native mongodb operation, as long as you can write mongo query statement can be very fast and convenient operation
// Function: Provide simple create, update, read one or read list , delete and check operations for mongo
// Note: Mutil is an operation type, can be insert, update, delete one or more documents
// Note: For mongo-driver operations are simplified, only frequently used mongo operations are retained, if you need more complex operations, please use the interface provided by mongo-driver
// TODO: 1): Add aggregate query

import (
	"context"
	"fmt"

	"mongo/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const ReturnMaxDocs int64 = 100

// Return true for existence, false for non-existence, error for error
func (dao *SDao) CheckExist(coll *mongo.Collection, flt string) (bool, error) {
	var err error
	var res bson.M
	if coll == nil || flt == "" {
		return false, fmt.Errorf("collection is nil or filter is empty")
	}
	opts := options.FindOneOptions{}
	fil, err := Json2Bson(flt)
	if err != nil {
		return false, err
	}
	err = coll.FindOne(context.TODO(), fil, &opts).Decode(&res)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// fetch one document
func (dao *SDao) Get(coll *mongo.Collection, p *SMongoGetParams) error {
	var err error
	if coll == nil {
		return fmt.Errorf("collection is nil")
	}

	if p.Flt == "" {
		return fmt.Errorf("filter is empty")
	}
	opts := options.FindOneOptions{}
	if p.Opt.Skip > 0 {
		opts.Skip = &p.Opt.Skip
	}
	if p.Opt.Sorts != "" {
		sort, err := Json2Bson(p.Opt.Sorts)
		if err != nil {
			err = fmt.Errorf("parser sorts json to bson error: %v", err)
			return err
		}
		opts.Sort = sort
	}
	if p.Opt.Fields != "" {
		fields, err := Json2Bson(p.Opt.Fields)
		if err != nil {
			err = fmt.Errorf("parser fields json to bson error: %v", err)
			return err
		}
		opts.Projection = fields
	}
	// parser filter
	fil, err := Json2Bson(p.Flt)
	if err != nil {
		err = fmt.Errorf("json to bson error: %v", err)
		return err
	}
	err = coll.FindOne(context.TODO(), fil, &opts).Decode(p.Rsp)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		return nil
	}
	return nil
}

// fetch documents list
func (dao *SDao) List(coll *mongo.Collection, p *SMongoListParams) error {
	var err error
	count := int64(0)
	if coll == nil {
		return fmt.Errorf("collection is nil")
	}
	opts := options.FindOptions{}
	if p.Flt == "" {
		return fmt.Errorf("filter is empty")
	}
	if p.Opt.Skip > 0 {
		opts.Skip = &p.Opt.Skip
	}
	if p.Opt.Limit > 0 {
		// 设置最大返回条数
		if p.Opt.Limit > ReturnMaxDocs {
			log.Log.Warnf("return doc limit is too large, set to %v", ReturnMaxDocs)
			p.Opt.Limit = ReturnMaxDocs
		}
		opts.Limit = &p.Opt.Limit
	}
	if p.Opt.Sorts != "" {
		sort, err := Json2Bson(p.Opt.Sorts)
		if err != nil {
			err = fmt.Errorf("parser sorts json to bson error: %v", err)
			return err
		}
		opts.Sort = sort
	}
	if p.Opt.Fields != "" {
		fields, err := Json2Bson(p.Opt.Fields)
		if err != nil {
			err = fmt.Errorf("parser fields json to bson error: %v", err)
			return err
		}
		opts.Projection = fields
	}
	// parser filter
	fil, err := Json2Bson(p.Flt)
	if err != nil {
		err = fmt.Errorf("parser filter json to bson error: %v", err)
		return err
	}

	cur, err := coll.Find(context.TODO(), &fil, &opts)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
		return nil
	}
	defer cur.Close(context.TODO())
	if p.Rsp == nil {
		log.Log.Warnf("result is nil ...")
		result := make([]bson.M, 0)
		p.Rsp = &result
	}
	if err := cur.All(context.TODO(), p.Rsp); err != nil {
		return err
	}
	// count
	count, err = coll.CountDocuments(context.TODO(), &fil)
	if err != nil {
		return err
	}
	p.Cnt = count
	return nil
}

//	insert one or many document
func (dao *SDao) Create(coll *mongo.Collection, p *SMongoCreateParams) error {
	if coll == nil {
		return fmt.Errorf("collection is nil")
	}
	if p.Docs == nil {
		return fmt.Errorf("docs is nil")
	}
	docs := *p.Docs
	if len(docs) < 1 {
		return fmt.Errorf("doc is empty")
	}
	if len(docs) == 1 {
		res, err := coll.InsertOne(context.TODO(), docs[0])
		if err != nil {
			return err
		}
		if res.InsertedID == nil {
			return fmt.Errorf("return insert id is empty")
		}
		p.Res = &[]interface{}{res.InsertedID}
		log.Log.Infof("inserted one document with ID [%v]\n", res.InsertedID)
		return nil
	} else {
		res, err := coll.InsertMany(context.TODO(), docs)
		if err != nil {
			return err
		}
		if res.InsertedIDs == nil {
			return fmt.Errorf("return insert id is empty")
		}
		r := make([]interface{}, 0)
		cnt := 0
		for _, id := range res.InsertedIDs {
			r = append(r, id)
			cnt += 1
		}
		p.Res = &r
		log.Log.Infof("inserted %d documents with IDs [%v]\n", cnt, res.InsertedIDs)
		return nil
	}
}

// update one or many document
func (dao *SDao) Update(coll *mongo.Collection, p *SMongoUpdateParams) error {
	var err error
	var res *mongo.UpdateResult
	if coll == nil {
		return fmt.Errorf("collection is nil")
	}
	if p.Flt == "" {
		return fmt.Errorf("filter is empty")
	}
	if p.Upd == "" {
		return fmt.Errorf("update is empty")
	}
	// parser filter
	fil, err := Json2Bson(p.Flt)
	if err != nil {
		return fmt.Errorf("parser filter json to bson error: %v", err)
	}
	// parser update
	upd, err := Json2Bson(p.Upd)
	if err != nil {
		return fmt.Errorf("parser update json to bson error: %v", err)
	}
	opts := options.UpdateOptions{}
	if p.Opt.Upsert {
		opts.SetUpsert(true)
	}
	if p.Opt.Multi {
		res, err = coll.UpdateMany(context.TODO(), fil, upd, &opts)
	} else {
		res, err = coll.UpdateOne(context.TODO(), fil, upd, &opts)
	}
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("no documents matched")
	}
	p.Res = res
	return nil
}

// delete one or many document
func (dao *SDao) Delete(coll *mongo.Collection, p *SMongoDeleteParams) error {
	var err error
	var res *mongo.DeleteResult
	if coll == nil {
		return fmt.Errorf("collection is nil")
	}
	if p.Flt == "" {
		return fmt.Errorf("filter is empty")
	}
	// parser filter
	fil, err := Json2Bson(p.Flt)
	if err != nil {
		err = fmt.Errorf("parser filter json to bson error: %v", err)
		return err
	}

	opts := options.DeleteOptions{}
	if p.Opt.Multi {
		res, err = coll.DeleteMany(context.TODO(), fil, &opts)
	} else {
		res, err = coll.DeleteOne(context.TODO(), fil, &opts)
	}
	if err != nil {
		return err
	}
	p.Res = res
	log.Log.Warnf("deleted %v documents\n", res.DeletedCount) // The number of documents deleted.
	return nil
}
