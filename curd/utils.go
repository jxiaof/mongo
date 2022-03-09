/*
 * @Descripttion:
 * @version:
 * @Author: hujianghong
 * @Date: 2022-03-09 11:50:20
 * @LastEditTime: 2022-03-09 11:52:06
 */
package curd

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// 将json字符串转换为bson
func Json2Bson(s string) (interface{}, error) {
	var doc interface{}
	err := bson.UnmarshalExtJSON([]byte(s), true, &doc)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// bson to struct
func Bson2Struct(bs *bson.M, res interface{}) error {
	data, err := bson.Marshal(bs)
	if err != nil {
		err = fmt.Errorf("bson error marshalling: %v, bson: %v  ", err, bs)
		return err
	}
	err = bson.Unmarshal(data, res)
	if err != nil {
		err = fmt.Errorf("json error from bson bytes: %v  ", err)
		return err
	}
	return nil
}
