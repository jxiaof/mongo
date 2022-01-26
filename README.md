<!--
 * @Descripttion: 
 * @version: 
 * @Author: hujianghong
 * @Date: 2022-01-26 18:08:29
 * @LastEditTime: 2022-01-26 23:49:28
-->
- 思路: bson在golang中编写比较繁琐.
  我的思路是将json转换为bson,贴近原生的mongodb操作
  只要会写mongo查询语句就能很快很方便的操作mongodb
- 功能: 提供对于mongo的简单增删改查操作
- 注意: Mutil是一个操作类型,可以是insert,update,delete一个或者多个document
- 注意: 针对mongo-driver操作做了简化,只保留经常使用的mongo操作,如果需要更复杂的操作需求,请使用mongo-driver提供的的接口
- 放在github目的主要是为了方便查看.mongo-driver的接口比较简单,相比关系型数据库没有更复杂的操作,没必要再封装mongo-driver.
我这样主要是写起来比较舒服熟悉,有那么一点mvc整套的舒畅感觉(.List / .Creat / .Get ...),好处就是不用去写繁琐的bson.
-  TODO: 1): 增加聚合查询 (aggregate)


使用方式: 如下所示或者查看main.go文件

```golang
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
```