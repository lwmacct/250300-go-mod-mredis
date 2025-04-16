package mredis

import (
	"context"
	"encoding/json"
	"errors"

	asRedis "github.com/redis/go-redis/v9"
)

type Client struct {
	*asRedis.Client
}

// JSONGetArray 获取 JSON 数据, []string
func (t *Client) JSONGetArray(key string) []string {
	ret := []string{}
	val := t.JSONGet(context.Background(), key)
	if val.Err() != nil {
		return ret
	}
	json.Unmarshal([]byte(val.Val()), &ret)
	return ret
}

type clientOpts func(*Client)

func NewClient(url string, opts ...clientOpts) (*Client, error) {
	t := &Client{}
	if url == "" {
		return nil, errors.New("url is required")
	}

	for _, opt := range opts {
		opt(t)
	}

	// 使用 url 解析
	opt, err := asRedis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// 启用 RESP3 响应
	opt.UnstableResp3 = true

	// 创建 Redis 客户端
	client := asRedis.NewClient(opt)
	// 测试 Redis 连接
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	t.Client = client
	return t, nil
}
