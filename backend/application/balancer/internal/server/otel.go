package server

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// http 服务初始化 TraceProvider
func initHttpTracerProvider(ctx context.Context, res *resource.Resource, conn string) (func(context.Context) error, error) {
	// Set up a trace exporter

	// 服务端的Jaeger支持HTTPS时使用
	// traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(conn))

	// 服务端的Jaeger不支持HTTPS时使用otlptracehttp.WithInsecure()显式声明只使用HTTP不安全的连接
	traceExporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create http trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// 将基于父span的采样率设置
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		sdktrace.WithBatcher(traceExporter),
		// 在资源中记录有关此应用程序的信息
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}

// gRPC 服务初始化 TraceProvider
func initGrpcTracerProvider(ctx context.Context, res *resource.Resource, conn string) (func(context.Context) error, error) {
	// 服务端的Jaeger使用HTTPS时不需要指定WithInsecure()
	// traceExporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(conn))

	// 服务端的Jaeger只支持HTTP时需要指定WithInsecure() 来显式声明只使用HTTP不安全的连接
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(conn),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc trace exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		// 将基于父span的采样率设置
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		sdktrace.WithBatcher(traceExporter),
		// 在资源中记录有关此应用程序的信息
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider.Shutdown, nil
}
