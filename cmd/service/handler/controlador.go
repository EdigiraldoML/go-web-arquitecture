package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/EdigiraldoML/go-web-arquitecture/internal/users"
	"github.com/EdigiraldoML/go-web-arquitecture/pkg/web"

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
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		users, err := u.service.GetAll()
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, users, ""))
	}
}

func (u *User) FilterByUrlParams() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		err = CheckQueryParams(c)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		filteredUsers, err := u.service.FilterByUrlParams(c)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, filteredUsers, ""))
	}
}

func (u *User) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		filteredUsers, err := u.service.GetUserByID(id)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, filteredUsers, ""))
	}
}

func (u *User) NewUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		err = c.Bind(&users.User{})
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		newUser, err := u.service.NewUser(c)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, newUser, ""))
	}
}

func (u *User) FullUpdate() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		var user users.User

		err = c.Bind(&user)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}
		if user.Nombre == "" {
			errMsg := "el nuevo nombre del usuario es requerido"
			c.JSON(400, web.NewResponse(400, nil, errMsg))
			return
		}
		if user.Apellido == "el nuevo apellido del usuario es requerido" {
			errMsg := ""
			c.JSON(400, web.NewResponse(400, nil, errMsg))
			return
		}
		if user.Email == "" {
			errMsg := "el nuevo email del usuario es requerido"
			c.JSON(400, web.NewResponse(400, nil, errMsg))
			return
		}
		if user.Edad == 0 {
			errMsg := "la nueva edad del usuario es requerida"
			c.JSON(400, web.NewResponse(400, nil, errMsg))
			return
		}
		if user.Altura == 0.0 {
			errMsg := "la nueva altura del usuario es requerida"
			c.JSON(400, web.NewResponse(400, nil, errMsg))
			return
		}

		user, err = u.service.FullUpdate(id, user.Nombre, user.Apellido, user.Email, user.Edad, user.Altura)
		if err != nil {
			statusCode := 400
			if err.Error() == "el usuario no fue encontrado" {
				statusCode = 404
			}

			c.JSON(statusCode, web.NewResponse(statusCode, nil, err.Error()))
			return
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user, ""))
	}
}

func (u *User) DeleteUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := CheckAccessToken(c)
		if err != nil {
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		err = u.service.DeleteUserByID(id)
		if err != nil {
			statusCode := 400
			if err.Error() == "el usuario no fue encontrado" {
				statusCode = 404
			}

			c.JSON(statusCode, web.NewResponse(statusCode, nil, err.Error()))
			return
		}

		data := "el usuario fue eliminado satisfactoriamente"
		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, data, ""))
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
			c.JSON(403, web.NewResponse(403, nil, err.Error()))
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		err = c.ShouldBindJSON(&newPartialUser)
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
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
				c.JSON(statusCode, web.NewResponse(statusCode, nil, err.Error()))
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

				c.JSON(statusCode, web.NewResponse(statusCode, nil, err.Error()))
				return
			}
		}

		c.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user, ""))
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

func CheckQueryParams(c *gin.Context) (err error) {
	queryParams := c.Request.URL.Query()
	for key, val := range queryParams {
		value := val[0]
		if value == "" { // nombre, apellido, email, fecha_de_creacion != ""
			err = fmt.Errorf("el valor del parametro %s no puede ser nulo", key)
			return
		}

		// "id", "nombre", "apellido", "email", "edad", "altura", "fecha_de_creacion"
		switch key {
		case "id":
			id, err := strconv.ParseInt(value, 10, 64)
			if err != nil || id <= 0 {
				err = fmt.Errorf("id debe ser un entero mayor a cero(recibido: %s)", value)
				return err
			}
		case "edad":
			edad, err := strconv.ParseInt(value, 10, 64)
			if err != nil || edad <= 0 {
				err = fmt.Errorf("edad debe ser un entero mayor a cero(recibido: %s)", value)
				return err
			}
		case "altura":
			altura, err := strconv.ParseFloat(value, 64)
			if err != nil || altura <= 0.0 {
				err = fmt.Errorf("altura debe ser un float mayor a cero(recibido: %s)", value)
				return err
			}
		}
	}

	return nil

}
