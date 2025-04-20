package main

import (
	"crud/controller"
	"crud/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title					 User Template Service
// @version					 1.0
// @description				 This is sample server template
// @termsOfService			 http://swagger.io/terms/

// @contact.name			 API Support
// @contact.url				 http://www.swagger.io/support
// @contact.email			 support@swagger.io

// @license.name			 Apache 2.0
// @license.url				 http://www.apache.org/licenses/LICENSE-2.0.html

// @host					 localhost:8080
// @BasePath				 /api/v1
// @schemes					 http https

// @externalDocs.description OpenAPI Swag Go
// @externalDocs.url         https://github.com/swaggo/swag#general-api-info
func main() {
	app := gin.Default()
	router := app.Group("/api/v1")
	controller.UserRoutes(router)

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := app.Run(":8080")
	if err != nil {
		fmt.Println("Something went wrong")
		return
	}
}
