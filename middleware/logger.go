package middleware

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tigerbig1242/evacuation-planning/logger"

	// "github.com/tigerbig1242/evacuation-planning/models"
	// "github.com/tigerbig1242/evacuation-planning/services"
	"go.uber.org/zap"
)

// Log Middleware is a middleware function that logs the request and response details
func LoggingMiddleware(c *fiber.Ctx) error {
	// Start time for logging
	start := time.Now()

	var requestBody []byte
	var requestBodyString string

	panicValue := recover()
	if panicValue != nil {
		logger.Log.Error("Panic in LoggingMiddleware",
			zap.Any("panic", panicValue),
			zap.String("stack", string(debug.Stack())),
		)

		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	requestBody = c.Body()
	requestBodyString = string(requestBody)

	if requestBody != nil {
		c.Context().SetBody(requestBody)
	}

	err := c.Next()

	if err != nil {
		logger.Log.Error("Error processing request",
			zap.Error(err),
			zap.String("url", c.OriginalURL()),
			zap.String("method", c.Method()),
			zap.String("request_body", requestBodyString),
		)

		return err
	}

	// requestBody := c.Body()

	// c.Context().SetBody(requestBody)

	// Process the request for logging
	// err := c.Next()

	// Calculate the processing time duration
	duration := time.Since(start)

	// Log the request information details
	fmt.Printf("Request Middleware URL: %s, - Method: %s, - Duration: %s\n", c.OriginalURL(), c.Method(), duration)

	// func() {
	// func() {
	panicValue = recover()
	if panicValue != nil {
		fmt.Printf("Panic in logger: %v\n", panicValue)
	}
	// }()

	if logger.Log == nil {
		fmt.Println("Warning: Logger is not initialized")
		return nil
	}

	childLogger := logger.Log.With(
		zap.String("Services", "Evacuation Planning"),
		zap.String("Version", "1.0.0"),
	)
	childLogger.Info("evacuation plans success",
		// zap.String("username", "tigerbig"),
		// zap.Float64("age", 25),
		// zap.Float64("height", 1.75),
		zap.String("Request URL", c.OriginalURL()),
		zap.String("Method", c.Method()),
		zap.Duration("Duration", duration),
		zap.Int("Status Code", c.Response().StatusCode()),
		zap.String("Request Body", requestBodyString),
		zap.String("Response Body", string(c.Response().Body())),
		zap.String("Date time logger", time.Now().Format("02-01-2006 15:04:05")),
	)

	// childLogger.Error("Failed to update evacuation plan",
	// 	zap.String("username", "tigerbig"),
	// 	zap.String("Request URL", c.OriginalURL()),
	// 	zap.String("Method", c.Method()),
	// 	zap.Duration("Duration", duration),
	// 	zap.Int("Status Code", c.Response().StatusCode()),
	// 	zap.String("Request Body", requestBodyString),
	// 	zap.String("Response Body", string(c.Response().Body())),
	// )
	// }()

	// logger.InitLogger()
	// defer logger.Log.Sync()
	// childLogger := logger.Log.With(
	// 	zap.String("Services", "Evacuation Planning"),
	// 	zap.String("Version", "1.0.0"),
	// 	zap.String("Request URL", c.OriginalURL()),
	// 	zap.String("Request Method", c.Method()),
	// 	zap.String("Duration", duration.String()),
	// )
	// childLogger.Info("Hello from Zap logger my TigerBig again!",
	// 	zap.String("username", "tigerbig"),
	// 	zap.Float64("age", 25),
	// 	zap.Float64("height", 1.75),
	// )

	return nil
}

// ยังไม่แน่ใจว่าทำงานได้ตามที่ต้องการมั๊ย ยังต้องทดสอบปรับแก้อีกเยอะ***
