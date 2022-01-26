/*
 * @Descripttion:
 * @version:
 * @Author: 江小凡
 * @Date: 2022-01-26 22:39:39
 * @LastEditTime: 2022-01-26 23:53:08
 */
package main

import (
	"fmt"
	"mongo/conf"
	"mongo/crud"
	"mongo/log"
)

var dao crud.Crud

type doc struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
}

func main() {
	var err error
	fmt.Println(conf.LOG_LEVEL)
	log.Log.Info("----- start mongo -----")

	dao = crud.NewCrud()
	f, err := dao.CheckExist(conf.CollA, fmt.Sprintf(`{"name":"%s"}`, "江小凡"))
	if err != nil {
		log.Log.Error("check exist error:", err)
	}
	fmt.Println("文档存在:", f)

	var d doc
	var docs []interface{}
	d.Name = "江小凡"
	d.Age = 18
	docs = append(docs, d)
	dao.Create(conf.CollA, docs)

	var lp *crud.SMongoListParams
	lp.Filter = fmt.Sprintf(`{"name":"%s"}`, "江小凡")
	lp.Options.Sorts = fmt.Sprintf(`{"age":%d}`, 1)
	lp.Options.Limit = 1
	lp.Options.Fields = `{"name":1}`
	lp.Options.Skip = 0

	res, err := dao.List(conf.CollA, lp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("文档数量", res.Count, "文档", res.Data)

	var gp *crud.SMongoGetParams
	gp.Filter = fmt.Sprintf(`{"name":"%s"}`, "江小凡")
	gp.Options.Sorts = fmt.Sprintf(`{"age":%d}`, 1)
	gp.Options.Skip = 0
	gp.Options.Fields = `{"name":1}`

	res2, err := dao.Get(conf.CollA, gp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("获取单个文档", res2)

	var up *crud.SMongoUpdateParams
	up.Filter = fmt.Sprintf(`{"name":"%s"}`, "江小凡")
	up.Update = fmt.Sprintf(`{"$set":{"age":%d}}`, 19)
	up.Options.Multi = true

	res3, err := dao.Update(conf.CollA, up)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("更新文档", res3)

	var dp *crud.SMongoDeleteParams
	dp.Filter = fmt.Sprintf(`{"name":"%s"}`, "江小凡")
	dp.Options.Multi = true

	res4, err := dao.Delete(conf.CollA, dp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("删除文档", res4)
}
