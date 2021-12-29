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

// GetUser godoc
// @Summary List all users in database
// @Tags Users
// @Description List all users that are recorder in database
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Router /users/GetAll [get]
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

// FilterByUrlParams godoc
// @Summary List users based on received url params
// @Tags Users
// @Description List users satisfying received url params
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id query int false "user id"
// @Param nombre query string false "user name"
// @Param apellido query string false "user last name"
// @Param email query string false "user email"
// @Param edad query int false "user age"
// @Param altura query number false "user height"
// @Param fecha_de_creacion query string false "user sign up date"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Router /users/ [get]
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

// GetUserByID godoc
// @Summary List user given the id
// @Tags Users
// @Description List user given the id as a param in url
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id path int true "user id"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Router /users/{id} [get]
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

// NewUser godoc
// @Summary Creates a new user
// @Tags Users
// @Description Creates a new user given params in body
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id body int true "user id, ignored param."
// @Param nombre body string true "user name"
// @Param apellido body string true "user last name"
// @Param email body string true "user email"
// @Param edad body int true "user edad"
// @Param altura body number true "user height"
// @Param activo body bool true "ignored, always true"
// @Param fecha_de_creacion body string true "user sign up date, use todays date. 'dd/mm/yyyy'"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Router /users/ [post]
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

// FullUpdate godoc
// @Summary Full update to an existing user
// @Tags Users
// @Description Full update to an existing user with body params
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id body int true "user id, ignored param."
// @Param nombre body string true "user name"
// @Param apellido body string true "user last name"
// @Param email body string true "user email"
// @Param edad body int true "user edad"
// @Param altura body number true "user height"
// @Param activo body bool true "ignored, always true"
// @Param fecha_de_creacion body string true "user sign up date, use todays date. 'dd/mm/yyyy'"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Router /users/{id} [put]
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

// DeleteUserByID godoc
// @Summary Delete an existing user
// @Tags Users
// @Description Delete an existing user given the id as an url param
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param id path int true "user id"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Failure 404 {object} web.Response
// @Router /users/{id} [delete]
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

// PartialUpdateToUser godoc
// @Summary Partial update to an existing user
// @Tags Users
// @Description Partial update to an existing user with body params
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param apellido body string false "user last name"
// @Param edad body int false "user edad"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 403 {object} web.Response
// @Failure 404 {object} web.Response
// @Router /users/{id} [patch]
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
