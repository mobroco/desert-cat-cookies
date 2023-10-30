package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mobroco/desert-cat-cookies/cmd"
)

func main() {
	fmt.Println("hello")
	if len(os.Args) > 1 && os.Args[1] == "redis" {
		ctx := context.Background()
		var db int
		if num := os.Getenv("REDIS_DB"); num != "" {
			db, _ = strconv.Atoi(num)
		}
		rdb := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS12,
			},
		})
		err := rdb.Set(ctx, "desert-cat-cookies", time.Now().String(), 0).Err()
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println(cmd.Execute())
	}

}
