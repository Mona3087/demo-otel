// This file contains the OpenTelemetry setup code.

package main

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

//Uncomment below to use stdout exporter
// func newTraceProvider() (*trace.TracerProvider, error) {
// 	traceExporter, err := stdouttrace.New(
// 		stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		return nil, err
// 	}

// 	traceProvider := trace.NewTracerProvider(
// 		trace.WithBatcher(traceExporter,
// 			// Default is 5s. Set to 1s for demonstrative purposes.
// 			trace.WithBatchTimeout(time.Second)),
// 	)
// 	return traceProvider, nil
// }

func newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(10*time.Minute))))
	//)
	return meterProvider, nil
}

func newTraceProvider() (*trace.TracerProvider, error) {
	// Create the Jaeger exporter
	// exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("https:///api/trace")))
	// if err != nil {
	// 		return nil, err
	// }
	// tp := trace.NewTracerProvider(
	// 		// Always be sure to batch in production.
	// 		trace.WithBatcher(exp),
	// 		// Record information about this application in a Resource.
	// 		trace.WithResource(resource.NewWithAttributes(
	// 				semconv.SchemaURL,
	// 				semconv.ServiceNameKey.String("demo-app"),
	// 		)),
	// )
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx)
	if err != nil {
		panic(err)
	}

	tracerProvider := trace.NewTracerProvider(trace.WithBatcher(exp))
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()
	//otel.SetTracerProvider(tracerProvider)
	return tracerProvider, nil
}
