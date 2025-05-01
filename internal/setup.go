package internal

import (
	"crud/docs"
	"crud/internal/controller"
	"crud/internal/repository"
	"crud/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter function to configure route and wire up dependencies
// @title					 User Template Service
// @version					 1.0
// @description				 This is sample server template
// @termsOfService			 http://swagger.io/terms/
//
// @contact.name			 API Support
// @contact.url				 http://www.swagger.io/support
// @contact.email			 support@swagger.io
//
// @license.name			 Apache 2.0
// @license.url				 http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host					 localhost:8080
// @BasePath				 /api/v1
// @schemes					 http https
//
// @externalDocs.description OpenAPI Swag Go
// @externalDocs.url         https://github.com/swaggo/swag#general-api-info
func SetupRouter(dbPool *pgxpool.Pool, app *gin.Engine) {
	v1Router := app.Group("/api/v1")
	setupV1Router(dbPool, v1Router)

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func setupV1Router(dbPool *pgxpool.Pool, router *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	userController.SetupRoutes(router)
}
