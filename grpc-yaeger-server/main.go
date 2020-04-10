package main

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-jaeger-server/account"
	"io"
	"log"
	"net"
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

	// 注册gRpc account服务
	server := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer))))
	account.RegisterAccountServer(server, &AccountServer{})
	reflection.Register(server)
	log.Println("succeed to register account service")

	// 监听gRpc account服务端口
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("starting register account service")

	// 开启gRpc account服务
	if err := server.Serve(listener); err != nil {
		log.Println(err)
		return
	}
}
