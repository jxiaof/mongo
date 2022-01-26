/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-20 11:25:35
 * @LastEditTime: 2022-01-26 23:34:20
 */
package crud

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// curd 接口定义
type Crud interface {
	CheckExist(colName, filter string) (bool, error)
	List(colName string, p *SMongoListParams) (*SMongoListRes, error)
	Get(colName string, p *SMongoGetParams) (*primitive.M, error)
	Create(colName string, doc []interface{}) ([]interface{}, error)
	Update(colName string, p *SMongoUpdateParams) (*mongo.UpdateResult, error)
	Delete(colName string, p *SMongoDeleteParams) (*mongo.DeleteResult, error)
}

// mongo集合定义
type Scoll struct {
	Cli    *mongo.Client
	ColMap map[string]*mongo.Collection
}

// 获取列表数据options参数定义
type SListParams struct {
	Skip   int64
	Limit  int64
	Sorts  string
	Fields string
}

// 列表接口入参定义
type SMongoListParams struct {
	Filter  string
	Options SListParams
}

// 列表接口出参定义
type SMongoListRes struct {
	Count int64
	Data  []primitive.M
}

// 获取单个数据options参数定义
type SGetParams struct {
	Skip   int64
	Sorts  string
	Fields string
}

// 获取单个数据入参定义
type SMongoGetParams struct {
	Filter  string
	Options SGetParams
}

// 更新操作options定义
type SUpdateParams struct {
	Multi  bool
	Upsert bool
}

// 更新操作入参定义
type SMongoUpdateParams struct {
	Filter  string
	Update  string
	Options SUpdateParams
}

// 删除操作options定义
type SDeleteParams struct {
	Multi bool
}

// 删除操作入参定义
type SMongoDeleteParams struct {
	Filter  string
	Options SUpdateParams
}
