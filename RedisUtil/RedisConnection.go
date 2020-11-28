package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	//ExampleClient()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "139.196.38.232:6379",
		Password: "adminfeng@.", // no password set
		DB:       0,             // use default DB
	})
	//val,error :=rdb.Get(ctx,"*").Result()
	keysize := rdb.DBSize(ctx)
	print("数量：" + keysize.String())
	val, _ := rdb.Scan(ctx, 10, "*", keysize.Val()).Val()
	for i := 0; i < len(val); i++ {
		//查询key
		fmt.Println("key--->", val[i])
		valuetype := rdb.Type(ctx, val[i])
		ts, _ := valuetype.Result()
		if ts == "list" {
			fmt.Println("类型是list")
			len := rdb.LLen(ctx, val[i]).Val()
			res := rdb.LRange(ctx, val[i], 0, len).Val()
			for _, i := range res {
				fmt.Println(i) // [val5 val4 val3 val2 val1 val99 val100]
			}
		} else {
			value := rdb.Get(ctx, val[i])

			fmt.Println("value-->", value)
		}

	}

}

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "139.196.38.232:6379",
		Password: "adminfeng@.", // no password set
		DB:       0,             // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
