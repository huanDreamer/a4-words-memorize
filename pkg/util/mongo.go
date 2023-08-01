package util

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FormatInsertedID(insertedID interface{}) (insertedIDString string) {
	// 将 _id 值转换为 string 类型
	switch v := insertedID.(type) {
	case primitive.ObjectID:
		insertedIDString = v.Hex()
	case string:
		insertedIDString = v
	default:
		fmt.Println("Unexpected _id type")
		return
	}
	return
}
