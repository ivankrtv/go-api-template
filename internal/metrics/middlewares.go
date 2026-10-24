package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
)

func NewMetricsMiddleware(serviceName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		startTime := time.Now()
		_ = c.Next()
		statusCode := c.Response().StatusCode()
		method := c.Request().URI().String()

		requestsCountMetric.With(prometheus.Labels{
			"service": serviceName,
			"method":  method,
			"code":    strconv.Itoa(statusCode),
		}).Add(1)

		requestsLatencyMetric.With(prometheus.Labels{
			"service": serviceName,
			"method":  method,
			"code":    strconv.Itoa(statusCode),
		}).Observe(time.Since(startTime).Seconds())

		return nil
	}
}
