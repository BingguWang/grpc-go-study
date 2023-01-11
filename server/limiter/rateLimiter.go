package limiter

import (
	"context"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/tap"
	"log"
	"sync"
)

/**
该函数从context中获取用户信息，然后为该用户创建一个限流器并注册到全局的map中，
这样每次收到该用户的请求时都去检查一下限流器，如果超过限流的阈值就直接返回错误拒绝这次请求。
*/

var (
	myTimeRateLimiterOnce sync.Once
	myTimeRateLimiter     *rate.Limiter
)

func NewTimeRateLimiter() *rate.Limiter {
	myTimeRateLimiterOnce.Do(func() {
		// r 也就是放入令牌的速率，如果传入时间,可以用rate.Every,多久放入一个令牌
		// b 是令牌桶的大小
		myTimeRateLimiter = rate.NewLimiter(50, 50)
	})
	return myTimeRateLimiter
}
func RateLimiter(ctx context.Context, info *tap.Info) (context.Context, error) {
	log.Println("...........获取限流器...........", info.FullMethodName)
	limiter := NewTimeRateLimiter()
	if !limiter.Allow() {
		return nil, status.Errorf(codes.ResourceExhausted, "client exceeded rate limit")
	}
	return ctx, nil
}
