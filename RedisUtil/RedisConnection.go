package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "strings"
)

var ctx = context.Background()

func main() {

	fmt.Println(buildDbStr("139.196.38.232:6379", "adminfeng@.", 0))
}

func buildDbStr(address, password string, db int) string {

	//获取redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	rdb.ConfigGet(ctx, "databases").Val()
	//获取key的数量
	keysize := rdb.DBSize(ctx)
	//获取所有key的值，游标设置0
	val, _ := rdb.Scan(ctx, 0, "*", keysize.Val()).Val()
	var resultStr = "{"
	for i := 0; i < len(val); i++ {
		//获取key对应值的的类型
		valuetype := rdb.Type(ctx, val[i])
		ts, _ := valuetype.Result()
		key := val[i]
		switch ts {

		case "list": //list类型
			valueLen := rdb.LLen(ctx, key).Val()
			res := rdb.LRange(ctx, key, 0, valueLen).Val()
			slice := make([]string, valueLen)
			var listStr = "["
			for _, i := range res {
				slice = append(slice, i)
				listStr += "\"" + i + "\","
			}
			listStr = listStr[0 : len(listStr)-1]
			listStr += "],"
			resultStr += "\"" + key + "\"" + ":" + listStr
			break
		case "set": //set类型
			setLen := rdb.LLen(ctx, key).Val()
			setList := rdb.SMembers(ctx, key).Val()
			setSlice := make([]string, setLen)
			var str = "["
			for _, i := range setList {
				setSlice = append(setSlice, i)
				str += "\"" + i + "\"" + ","
			}
			str = str[0 : len(str)-1]
			str += "],"
			resultStr += "\"" + key + "\"" + ":" + str
			break
		case "hash": //hash类型
			hashStr := ""
			hashKeys := rdb.HKeys(ctx, key).Val()
			for _, i := range hashKeys {
				//fmt.Println(i)
				hashValues := rdb.HGetAll(ctx, key).Val()
				hashStr += "\"" + i + "\":["
				for _, j := range hashValues {
					hashStr += "\"" + j + "\","
					//fmt.Println( j)
				}
				hashStr = hashStr[0 : len(hashStr)-1]
				hashStr += "],"
			}

			resultStr += hashStr
			break
		case "zset":
			zsetStr := "\"" + key + "\":["
			zsetlen := rdb.LLen(ctx, key).Val()
			zsetValue := rdb.ZRange(ctx, key, 0, zsetlen).Val()
			for i, _ := range zsetValue {
				zsetStr += "\"" + zsetValue[i] + "\","
			}
			zsetStr = zsetStr[0 : len(zsetStr)-1]
			zsetStr += "],"
			resultStr += zsetStr
			break
		default:
			value := rdb.Get(ctx, key).Val()
			resultStr += "\"" + key + "\"" + ":" + "\"" + value + "\","
		}
	}
	resultStr = resultStr[0:len(resultStr)-1] + "}"
	return resultStr
}
