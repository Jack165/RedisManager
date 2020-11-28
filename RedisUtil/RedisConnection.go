package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {

	//获取redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     "139.196.38.232:6379",
		Password: "adminfeng@.", // no password set
		DB:       0,             // use default DB
	})

	//获取key的数量
	keysize := rdb.DBSize(ctx)
	print("数量：" + keysize.String())
	//获取所有key的值，游标设置0
	val, _ := rdb.Scan(ctx, 0, "*", keysize.Val()).Val()
	for i := 0; i < len(val); i++ {
		//查询key，打印
		fmt.Println("key--->", val[i])
		//获取key对应值的的类型
		valuetype := rdb.Type(ctx, val[i])
		ts, _ := valuetype.Result()
		//如果是list类型就遍历显示
		if ts == "list" {
			fmt.Println("类型是list")
			len := rdb.LLen(ctx, val[i]).Val()
			res := rdb.LRange(ctx, val[i], 0, len).Val()
			for _, i := range res {
				fmt.Println(i) // [val5 val4 val3 val2 val1 val99 val100]
			}
		} else {
			//如果不是list就直接打印
			value := rdb.Get(ctx, val[i])

			fmt.Println("value-->", value)
		}

	}
}
