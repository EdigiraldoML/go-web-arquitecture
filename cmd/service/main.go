package main

import (
	"fmt"

	"github.com/EdigiraldoML/go-web-arquitecture/cmd/service/handler"
	"github.com/EdigiraldoML/go-web-arquitecture/internal/users"
	"github.com/EdigiraldoML/go-web-arquitecture/pkg/store"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	// controller.LoadUsersFromJSON(pathUsersJSON)

	router := gin.Default()

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
