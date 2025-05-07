package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lwmacct/250300-go-mod-mredis/pkg/mredis"
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	vredis, err := mredis.NewClient(redisURL)
	if err != nil {
		fmt.Println(err)
	}
	res := vredis.Raw.Ping(context.Background())
	fmt.Printf("%s\n", res.Val())
}
