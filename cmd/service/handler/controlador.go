package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/EdigiraldoML/go-web-arquitecture/internal/users"

	"github.com/gin-gonic/gin"
)

type User struct {
	service users.Service
}

func CreateUser(u users.Service) *User {
	newUser := &User{
		service: u,
	}

	return newUser
}

func (u *User) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		users, err := u.service.GetAll()
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

func (u *User) FilterByUrlParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		fmt.Println("Method FilterByUrlParams called.")

		filteredUsers, err := u.service.FilterByUrlParams(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, filteredUsers)
	}
}

func (u *User) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		filteredUsers, err := u.service.GetUserByID(id)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, filteredUsers)
	}
}

func (u *User) NewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		newUser, err := u.service.NewUser(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, newUser)
	}
}

func (u *User) FullUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var user users.User

		err = c.Bind(&user)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if user.Nombre == "" {
			c.JSON(400, gin.H{"error": "el nuevo nombre del usuario es requerido"})
			return
		}
		if user.Apellido == "" {
			c.JSON(400, gin.H{"error": "el nuevo apellido del usuario es requerido"})
			return
		}
		if user.Email == "" {
			c.JSON(400, gin.H{"error": "el nuevo email del usuario es requerido"})
			return
		}
		if user.Edad == 0 {
			c.JSON(400, gin.H{"error": "la nueva edad del usuario es requerida"})
			return
		}
		if user.Altura == 0.0 {
			c.JSON(400, gin.H{"error": "la nueva altura del usuario es requerida"})
			return
		}

		user, err = u.service.FullUpdate(id, user.Nombre, user.Apellido, user.Email, user.Edad, user.Altura)
		if err != nil {
			statusCode := 400
			if err.Error() == "el usuario no fue encontrado" {
				statusCode = 404
			}

			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, user)
	}
}

func (u *User) DeleteUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = u.service.DeleteUserByID(id)
		if err != nil {
			statusCode := 400
			if err.Error() == "el usuario no fue encontrado" {
				statusCode = 404
			}

			c.JSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"data": "el usuario fue eliminado satisfactoriamente"})
	}
}

func (u *User) PartialUpdateToUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		type partialUser struct {
			Apellido string `json:"apellido"`
			Edad     int64  `json:"edad"`
		}

		var newPartialUser partialUser

		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		err = c.ShouldBindJSON(&newPartialUser)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		apellido := newPartialUser.Apellido
		edad := newPartialUser.Edad
		var user users.User

		if apellido != "" {
			user, err = u.service.UpdateUserLastName(id, apellido)
			if err != nil {
				statusCode := 400
				if err.Error() == "el usuario no fue encontrado" {
					statusCode = 404
				}

				c.JSON(statusCode, gin.H{"error": err.Error()})
				return
			}
		}

		if edad != 0 {
			user, err = u.service.UpdateUserAge(id, edad)
			if err != nil {
				statusCode := 400
				if err.Error() == "el usuario no fue encontrado" {
					statusCode = 404
				}

				c.JSON(statusCode, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(200, user)
	}
}

func CheckAccessToken(c *gin.Context) (err error) {
	token := c.GetHeader("token")
	acceptedToken := os.Getenv("TOKEN")
	if token == "" {
		err = errors.New("el token de acceso no fue proporcionado")
	} else if token != acceptedToken {
		err = errors.New("el token enviado no es correcto")
	}

	return err
}
