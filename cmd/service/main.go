package main

import (
	"fmt"
	"os"

	"github.com/EdigiraldoML/go-web-arquitecture/cmd/service/handler"
	"github.com/EdigiraldoML/go-web-arquitecture/docs"
	"github.com/EdigiraldoML/go-web-arquitecture/internal/users"
	"github.com/EdigiraldoML/go-web-arquitecture/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title MeLi Bootcamp API
// @version 1.0
// @description API built to manage users.
// @termsOfService https://developers.mercadolibre.com.co/es_ar/terminos-y-condiciones

// @contact.name API Support
// @contact.url https://developers.mercadolibre.com.ar/support

// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0
func main() {

	pathUsersJSON := "users.json"

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	db := store.NewStorage(store.FileType, pathUsersJSON)
	repository := users.CreateRepository(db)
	service := users.CreateService(repository)
	controller := handler.CreateUser(service)

	router := gin.Default()

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	usrs := router.Group("/users")
	usrs.GET("/", controller.FilterByUrlParams())
	usrs.GET("/GetAll", controller.GetAll())
	usrs.GET("/:id", controller.GetUserByID())
	usrs.POST("/", controller.NewUser())
	usrs.PUT("/:id", controller.FullUpdate())
	usrs.DELETE("/:id", controller.DeleteUserByID())
	usrs.PATCH("/:id", controller.PartialUpdateToUser())

	router.Run()
}
