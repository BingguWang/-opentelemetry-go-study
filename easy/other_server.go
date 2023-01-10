package main

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"math"
	"strconv"
)

/**
假设这里的每个方法代表一个微服务的方法
*/

// CallOtherService 调用其他服务去计算，得到计算结果
func CallOtherService(ctx context.Context, n float64) (float64, error) {
	_, span := otel.Tracer(name).Start(ctx, "CallOtherService")
	defer span.End()

	if n < 0 {
		return -1, errors.New("input value can be a minus")
	}
	ret := math.Sqrt(n)
	return ret, nil
}

type A struct {
	Name string
}

func (a *A) GetName() string {
	return a.Name
}

// CallOtherServiceAndGetANumber 调用其他服务去获取一个数值
func (s *EasyService) CallOtherServiceAndGetANumber(ctx context.Context) (float64, error) {
	_, span := otel.Tracer(name).Start(ctx, "CallOtherServiceAndGetANumber")
	defer span.End()

	var n float64
	_, err := fmt.Fscanf(s.r, "%f\n", &n) // 因为是模拟，这里我们从控制台输入一个数就行了
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	nStr := strconv.FormatFloat(n, 'f', 2, 64)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, nil
}
