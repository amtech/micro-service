package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"grpc-jaeger-client/account"
	"io"
	"log"
)

// 初始化jaeger
func initJaeger() (tracer opentracing.Tracer, closer io.Closer, err error) {
	// 构造配置信息
	cfg := &config.Configuration{
		// 设置服务名称
		ServiceName: "ServiceAmount",

		// 设置采样参数
		Sampler: &config.SamplerConfig{
			Type:  "const", // 全采样模式
			Param: 1,       // 开启全采样模式
		},
	}

	// 生成一条新tracer
	tracer, closer, err = cfg.NewTracer()
	if err == nil {
		// 设置tracer为全局单例对象
		opentracing.SetGlobalTracer(tracer)
	}
	return
}

func main() {
	// 初始化jaeger，创建一条新tracer
	tracer, closer, err := initJaeger()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	defer closer.Close()
	log.Println("succeed to init jaeger")

	// 创建一个新span，作为父span
	span := tracer.StartSpan("CalculateFee")

	// 函数返回时关闭span
	defer span.Finish()

	// 生成span的context
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// 连接gRpc server
	conn, err := grpc.Dial("localhost:8080",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer),
		)))
	if err != nil {
		log.Println(err)
		return
	}

	// 创建gRpc计费服务客户端
	client := account.NewAccountClient(conn)

	// 将父span的context作为参数，调用检测用户余额的gRpc微服务
	checkBalanceResponse, err := client.CheckBalance(ctx,
		&account.CheckBalanceRequest{
			Account: "user account",
		})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(checkBalanceResponse)

	// 将父span的context作为参数，调用扣款的gRpc微服务
	reductionResponse, err := client.Reduction(ctx,
		&account.ReductionRequest{
			Account: "user account",
			Amount: 1,
		})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(reductionResponse)
}
