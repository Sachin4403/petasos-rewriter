package main

import (
	"strconv"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var labelNames = []string{"code", "method", "host", "url"}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "server_request_count",
		Help: "total incoming HTTP requests",
	},
	labelNames,
)

var requestDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "server_request_duration_seconds",
		Help: "tracks incoming request durations in seconds",
	},
	labelNames,
)

func provideMetrics(e *echo.Echo) {
	if err := prometheus.Register(totalRequests); err != nil {
		e.Logger.Fatal(err)
	}
	if err := prometheus.Register(requestDuration); err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(getMiddleware())
	e.GET("/metrics", echoprometheus.NewHandler())
}

func getMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			err := next(c)
			elapsed := time.Since(startTime).Seconds()
			status := strconv.Itoa(c.Response().Status)

			values := make([]string, len(labelNames))
			values[0] = status
			values[1] = c.Request().Method
			values[2] = c.Request().Host
			values[3] = c.Path()

			totalRequests.WithLabelValues(values...).Inc()
			requestDuration.WithLabelValues(values...).Observe(elapsed)
			return err
		}
	}
}
