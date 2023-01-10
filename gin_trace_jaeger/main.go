package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	jg "github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"net/http"
)

const serviceName = "gin-service"

func main() {
	r := gin.Default()

	/*	如果要上报handler内部的调用情况就得加上这些
		tp, err2 := tracerProvider("http://localhost:14268/api/traces")
		if err2 != nil {
			panic(err2)
		}
		ctx, cancelFunc := context.WithCancel(context.Background())
		defer cancelFunc()
		defer func(c context.Context) {
			tp.Shutdown(c)
		}(ctx)
		otel.SetTracerProvider(tp)*/

	// 设置jaeger配置信息, 这里和grpc那边使用jaeger是一样的
	cfg := &jaegerConfig.Configuration{
		ServiceName: serviceName, //对其发起请求的的调用链，叫什么服务
		Sampler: &jaegerConfig.SamplerConfig{ //采样策略的配置
			Type:  "const",
			Param: 1,
			/**
			"const" : 0 or 1 for always false/true respectively
			"probabilistic" sampler, a probability between 0 and 1
			"rateLimiting" sampler, the number of spans per second
			"remote" sampler, param is the same as for "probabilistic"
			*/
		},
		Reporter: &jaegerConfig.ReporterConfig{ //配置客户端如何上报trace信息，所有字段都是可选的
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831", // jaeger本地agent的地址,span会被send到此agent
		},
		//Token配置
		Tags: []opentracing.Tag{ //设置tag，token等信息可存于此
			//opentracing.Tag{Key: "token", Value: token}, //设置token
		},
	}
	// 根据jaeger配置来创建得到tracer
	tracer, closer, err := cfg.NewTracer(jaegerConfig.Logger(jg.StdLogger))
	defer closer.Close()
	if err != nil {
		panic(fmt.Sprintf("ERROR: fail init Jaeger: %v\n", err))
	}
	//用中间件实现拦截器
	r.Use(ginhttp.Middleware(tracer, ginhttp.OperationNameFunc(func(r *http.Request) string {
		/**
		官方默认 OperationName 是 HTTP + HttpMethod,
		    //建议使用 HTTP + HttpMethod + URL 可以分析到具体接口，具体用法如下
		    //PS：Restful 接口主要URL应该是参数名，不是具体参数值。 如： 正确：/user/{id}， 错误：/user/1
		*/
		return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.String()) // OperationName
	})))

	// router
	r.GET("/ping", MyHandler)

	// run
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}

}
func MyHandler(c *gin.Context) {
	newCtx, span := otel.Tracer("").Start(c, "MyHandler")
	defer span.End()
	hello := SayHello(newCtx)

	c.JSON(http.StatusOK, gin.H{
		"message": hello,
	})
	return
}

func SayHello(c context.Context) string {
	//_, span := otel.Tracer("").Start(c, "SayHello")
	//defer span.End()
	return "hello guys"
}

func tracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url))) // span会被发送到url
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL, // 使用semconv包为资源属性提供常规名称。
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	return tp, nil
}
