<!--
 * @Descripttion: 
 * @version: 
 * @Author: hujianghong
 * @Date: 2022-01-26 18:08:29
 * @LastEditTime: 2022-03-09 22:57:29
-->

// Idea: bson is tedious to write in golang, my idea is to convert json to bson, close to the native mongodb operation, as long as you can write mongo query statement can be very fast and convenient operation

// Function: Provide simple create, update, read one or read list , delete and check operations for mongo

// Note: Mutil is an operation type, can be insert, update, delete one or more documents

// Note: For mongo-driver operations are simplified, only frequently used mongo operations are retained, if you need more complex operations, please use the interface provided by mongo-driver

// TODO: 1): Add aggregate query

Usage: as shown below or view the main.go file

```go
func main() {
	var err error
	fmt.Println(conf.LOG_LEVEL)
	log.Log.Info("----- start mongo -----")

	dao = crud.NewCrud()
	f, err := dao.CheckExist(conf.CollA, fmt.Sprintf(`{"name":"%s"}`, "bob"))
	if err != nil {
		log.Log.Error("check exist error:", err)
	}
	fmt.Println("doc exsits:", f)

	var d doc
	var docs []interface{}
	d.Name = "bob"
	d.Age = 18
	docs = append(docs, d)
	dao.Create(conf.CollA, docs)

	var lp *crud.SMongoListParams
	lp.Filter = fmt.Sprintf(`{"name":"%s"}`, "bob")
	lp.Options.Sorts = fmt.Sprintf(`{"age":%d}`, 1)
	lp.Options.Limit = 1
	lp.Options.Fields = `{"name":1}`
	lp.Options.Skip = 0

	res, err := dao.List(conf.CollA, lp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("doc count", res.Count, "doc", res.Data)

	var gp *crud.SMongoGetParams
	gp.Filter = fmt.Sprintf(`{"name":"%s"}`, "bob")
	gp.Options.Sorts = fmt.Sprintf(`{"age":%d}`, 1)
	gp.Options.Skip = 0
	gp.Options.Fields = `{"name":1}`

	res2, err := dao.Get(conf.CollA, gp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("fetch one doc", res2)

	var up *crud.SMongoUpdateParams
	up.Filter = fmt.Sprintf(`{"name":"%s"}`, "bob")
	up.Update = fmt.Sprintf(`{"$set":{"age":%d}}`, 19)
	up.Options.Multi = true

	res3, err := dao.Update(conf.CollA, up)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("update doc", res3)

	var dp *crud.SMongoDeleteParams
	dp.Filter = fmt.Sprintf(`{"name":"%s"}`, "doc")
	dp.Options.Multi = true

	res4, err := dao.Delete(conf.CollA, dp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("delete doc", res4)
}
```
