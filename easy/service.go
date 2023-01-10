package main

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"io"
	"log"
)

const name = "easy"

type EasyService struct {
	r io.Reader
	l *log.Logger
}

func (s *EasyService) EasyServiceHandler(ctx context.Context) error {
	// Tracer创建一个命名为name的tracer
	// Start创建一个span
	newCtx, span := otel.Tracer(name).Start(ctx, "Run")
	defer span.End()

	// 调用服务CallGetNumber获取数值
	n, err := s.CallGetNumber(newCtx)
	if err != nil {
		return err
	}
	// 调用服务CallComputeSqrt计算结果
	r, err := s.CallComputeSqrt(newCtx, n)
	if err != nil {
		return err
	}
	fmt.Println("输出结果：", r)
	return nil
}

func (s *EasyService) CallComputeSqrt(ctx context.Context, n float64) (float64, error) {
	newCtx, span := otel.Tracer(name).Start(ctx, "CallComputeSqrt")
	defer span.End()

	// 这模拟调用其他服务
	return CallOtherService(newCtx, n)
}

// CallGetNumber 假设这是去调用一个获取数值的服务
func (s *EasyService) CallGetNumber(ctx context.Context) (float64, error) {
	_, span := otel.Tracer(name).Start(ctx, "CallGetNumber")
	defer span.End()

	// 这模拟调用其他服务
	return s.CallOtherServiceAndGetANumber(ctx)
}
