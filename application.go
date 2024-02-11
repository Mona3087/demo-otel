// demo application to demonstrate opentelemetry

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	tracer  = otel.Tracer("catalogue")
	meter   = otel.Meter("catalogue")
	reqCnt metric.Int64Counter
	httpreq int
)

func init() {
	var err error
	reqCnt, err = meter.Int64Counter("request_count",
		metric.WithDescription("The number of times api is called"),
		metric.WithUnit("{count}"))
	if err != nil {
		panic(err)
	}
}

func catalogue(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "catalogue")
	defer span.End()
	httpreq = httpreq + 1

	rollValueAttr := attribute.Int("RequestCnt.value", httpreq)
	span.SetAttributes(rollValueAttr)
	reqCnt.Add(ctx, 1, metric.WithAttributes(rollValueAttr))

	resp := strconv.Itoa(httpreq) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Printf("Write failed: %v\n", err)
	}
	cart(ctx)
	time.Sleep(1*time.Second)
}

func cart(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "cart")
	defer span.End()
	
	time.Sleep(2*time.Second)

	order(ctx)
}

func order(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "order")
	defer span.End()
	
	time.Sleep(1*time.Second)
	payment(ctx)

}

func payment(ctx context.Context) {
	_, span := tracer.Start(ctx, "payment")
	defer span.End()
	
	time.Sleep(3*time.Millisecond)
}

