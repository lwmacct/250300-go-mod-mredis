package mredis

import (
	"context"
	"fmt"

	asRedis "github.com/redis/go-redis/v9"
)

type List struct {
	Raw *asRedis.Client
}

// Ltrim 裁剪列表
func (t *List) Ltrim(ctx context.Context, keyName string, count int64) ([]string, error) {
	// 计算 count-1，提前计算避免重复计算
	trimCount := count - 1

	// Lua 脚本：获取前 count 个元素并裁剪列表
	script := fmt.Sprintf(`
		local values = redis.call('lrange', KEYS[1], 0, %d)
		redis.call('ltrim', KEYS[1], %d + 1, -1)
		return values
	`, trimCount, trimCount)

	// 执行 Lua 脚本，获取值并裁剪列表
	val, err := t.Raw.Eval(ctx, script, []string{keyName}, trimCount).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to execute Lua script on key %s: %v", keyName, err)
	}

	// 处理 Lua 脚本返回的数据
	values, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("error processing data from Redis")
	}

	// 跳过非字符串值并收集结果
	var result []string
	for _, v := range values {
		if str, isString := v.(string); isString {
			result = append(result, str)
		}
	}

	// 返回结果
	return result, nil
}
