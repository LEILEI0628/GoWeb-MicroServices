package ioc

import (
	"github.com/LEILEI0628/GinPro/middleware/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitLimiter(redisClient redis.Cmdable) limiter.Limiter {
	// 滑动窗口算法 1000/1s 使用redis统计请求数量
	return limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)

}
