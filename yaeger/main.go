package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"io"
	"log"
	"time"
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
			Param: 1,       // 开启状态
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

// 检测用户余额，模拟子span任务
func CheckBalance(request string, ctx context.Context) {
	// 创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "CheckBalance")

	// 模拟系统进行一系列的操作，耗时1/3秒
	time.Sleep(time.Second / 3)

	// 示例：将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", "CheckBalance reply")

	// 结束当前span
	span.Finish()

	log.Println("CheckBalance is done")
}

// 从用户账户扣款，模拟子span任务
func Reduction(request string, ctx context.Context) {
	// 创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "Reduction")

	// 模拟系统进行一系列的操作，耗时1/2秒
	time.Sleep(time.Second / 2)

	// 示例：将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", "Reduction reply")

	// 结束当前span
	span.Finish()

	log.Println("Reduction is done")
}

func main() {
	// 初始化jaeger，创建一条新tracer
	tracer, closer, err := initJaeger()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	defer closer.Close()

	// 创建一个新span，作为父span，开始计费过程
	span := tracer.StartSpan("CalculateFee")

	// 生成父span的context
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// 示例：设置一个span标签信息
	span.SetTag("db.instance", "customers")
	// 示例：输出一条span日志信息
	span.LogKV("event", "timed out")

	// 将父span的context作为参数，调用检测用户余额函数
	CheckBalance("CheckBalance request", ctx)

	// 将父span的context作为参数，调用扣款函数
	Reduction("Reduction request", ctx)

	// 结束父span
	span.Finish()
}
