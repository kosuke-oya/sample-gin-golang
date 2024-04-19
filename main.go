package main

import (
	"httpserver/docs"
	handlers "httpserver/handlers"
	"httpserver/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	r.Use(middleware.Logger(logger))
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
