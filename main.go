package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"httpserver/docs"
	handlers "httpserver/handlers"
	"io"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a httpserver
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      your-url.com
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// get env value
	ENV_VALUE := os.Getenv("ENV_KEY")

	// set zap logger as default logger
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowAllOrigins: true,
		AllowWebSockets: true,
		MaxAge:          12 * time.Hour,
	}))

	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zap.Field {
			start := time.Now()
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// CloudLoggingに送信用のメッセージを定義
			heaserBytes, err := json.Marshal(c.Request.Header)
			if err != nil {
				fmt.Println("failed to marshal header")
			}
			headerStr := string(heaserBytes)
			return []zap.Field{
				zap.String("uuid", xid.New().String()),
				zap.Int("status", c.Writer.Status()),
				zap.Int64("content_length", c.Request.ContentLength),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("query", c.Request.URL.RawQuery),
				zap.String("request_body", string(bodyBytes)),
				zap.String("ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
				zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
				zap.Duration("elapsed", time.Since(start)),
				zap.String("header", headerStr),
				zap.Int("response_size(bytes)", c.Writer.Size()),
			}
		}),
	}))

	docs.SwaggerInfo.Title = "API Docs"
	docs.SwaggerInfo.Description = "This is a http server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// swagger
	// use ginSwagger middleware to serve the API docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// GET /sample
	r.GET("/sample", func(c *gin.Context) {
		handlers.SampleHandler(c, ENV_VALUE)
	})
	r.Run(":8080")

}
