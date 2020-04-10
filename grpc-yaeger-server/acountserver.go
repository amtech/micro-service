package main

import (
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"grpc-jaeger-server/account"
	"time"
)

type AccountServer struct{}

// 检测用户余额微服务，模拟子span任务
func (s *AccountServer) CheckBalance(ctx context.Context, request *account.CheckBalanceRequest) (response *account.CheckBalanceResponse, err error) {
	response = &account.CheckBalanceResponse{
		Reply: "CheckBalance Reply", // 处理结果
	}

	// 创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "CheckBalance")

	// 模拟系统进行一系列的操作，耗时1/3秒
	time.Sleep(time.Second / 3)

	// 将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", response)

	// 结束当前span
	span.Finish()

	return response, err
}

// 从用户账户扣款微服务，模拟子span任务
func (s *AccountServer) Reduction(ctx context.Context, request *account.ReductionRequest) (response *account.ReductionResponse, err error) {
	response = &account.ReductionResponse{
		Reply: "Reduction Reply", // 处理结果
	}

	// 创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "Reduction")

	// 模拟系统进行一系列的操作，耗时1/3秒
	time.Sleep(time.Second / 3)

	// 将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", response)

	// 结束当前span
	span.Finish()
	return response, err
}
