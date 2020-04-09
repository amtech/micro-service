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

// 检测用户余额，模拟子span任务
func checkBalance(request string, ctx context.Context) {
	// 创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "checkBalance")

	// 模拟系统进行一系列的操作，耗时1/3秒
	time.Sleep(time.Second / 3)

	// 生成处理结果reply
	reply := "checkBalance reply"

	// 将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", reply)

	// 结束当前span
	span.Finish()

	// 打印输入数据结果
	log.Println(request)
	log.Println(reply)
}

// 从用户账户扣款，模拟子span任务
func reduction(request string, ctx context.Context) {
	// 创建子span

	span, _ := opentracing.StartSpanFromContext(ctx, "reduction")
	// 模拟系统进行一系列的操作，耗时1/2秒
	time.Sleep(time.Second / 2)

	// 生成处理结果reply
	reply := "reduction reply"

	// 将需要追踪的信息放入tag
	span.SetTag("request", request)
	span.SetTag("reply", reply)

	// 结束当前span
	span.Finish()

	// 打印输入数据和输出结果
	log.Println(request)
	log.Println(reply)
}

func main() {

	// 初始化jaeger，创建一条新tracer
	tracer, closer, err := initJaeger()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	defer closer.Close()

	// 创建一个新span，作为父span
	span := tracer.StartSpan("CalculateAmount")

	// 设置一个标签信息tag
	span.SetTag("db.instance", "customers")

	// 设置一条日志信息log
	span.LogKV("event", "timed out")

	// 生成span的context
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// 将父span的context作为参数，调用检测用户余额函数，其中生成子span
	checkBalance("checkBalance request", ctx)

	// 将父span的context作为参数，调用扣款函数，生成子span
	reduction("reduction request", ctx)

	// 结束父span
	span.Finish()
}
