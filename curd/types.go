/*
 * @Descripttion:
 * @version:
 * @Author: hujianghong
 * @Date: 2022-01-20 11:25:35
 * @LastEditTime: 2022-03-09 19:17:02
 */
package curd

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type MgoCrud interface {
	CheckExist(colName, filter string) (bool, error)
	Get(coll *mongo.Collection, p *SMongoGetParams) error
	List(coll *mongo.Collection, p *SMongoListParams) error
	Create(coll *mongo.Collection, p *SMongoCreateParams) error
	Update(coll *mongo.Collection, p *SMongoUpdateParams) error
	Delete(coll *mongo.Collection, p *SMongoDeleteParams) error
}

type SMongoListParams struct {
	Flt string
	Opt SListParams
	Rsp interface{} // return struct pointer
	Cnt int64
}

type SMongoGetParams struct {
	Flt string
	Opt SGetParams
	Rsp interface{} // return struct pointer
}

type SMongoCreateParams struct {
	Docs *[]interface{}
	Res  *[]interface{}
}
type SMongoUpdateParams struct {
	Flt string
	Upd string
	Opt SUpdateParams
	Res *mongo.UpdateResult
	// Befor   interface{}
	// After   interface{}
}

type SMongoDeleteParams struct {
	Flt string
	Opt SUpdateParams
	Res *mongo.DeleteResult
}

type SListParams struct {
	Skip   int64
	Limit  int64
	Sorts  string
	Fields string
}

type SGetParams struct {
	Skip   int64
	Sorts  string
	Fields string
}

type SUpdateParams struct {
	Multi  bool
	Upsert bool
}

type SDeleteParams struct {
	Multi bool
}
