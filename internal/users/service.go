package users

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetAll() (users *Users, err error)
	Store(id int64, nombre string, apellido string, email string, edad int64, altura float64, activo bool, fecha_de_creacion string) (user User, err error)
	FilterByUrlParams(c *gin.Context) (filteredUsers Users, err error)
	GetUserByID(id int64) (user User, err error)
	NewUser(c *gin.Context) (user User, err error)
	FullUpdate(id int64, nombre string, apellido string, email string, edad int64, altura float64) (user User, err error)
	DeleteUserByID(id int64) (err error)
	UpdateUserLastName(id int64, apellido string) (user User, err error)
	UpdateUserAge(id int64, edad int64) (user User, err error)
}

type service struct {
	repository Repository
}

func CreateService(r Repository) Service {
	newService := &service{
		repository: r,
	}

	return newService
}

func (s *service) GetAll() (users *Users, err error) {
	return s.repository.GetAll()
}

func (s *service) Store(id int64, nombre string, apellido string, email string, edad int64, altura float64, activo bool, fecha_de_creacion string) (user User, err error) {
	user, err = s.repository.Store(id, nombre, apellido, email, edad, altura, activo, fecha_de_creacion)

	return user, err
}

func (s *service) FilterByUrlParams(c *gin.Context) (filteredUsers Users, err error) {

	availableParams := CheckAvailableParamsFromGinContext(c)
	searchedUser, err := CreateSearchedUser(availableParams, c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(availableParams)
	users, err := s.repository.GetAll()
	if err != nil {
		return filteredUsers, err
	}

	filteredUsers = GetUsersWithGivenParams(availableParams, *users, searchedUser)
	return filteredUsers, err
}

func (s *service) GetUserByID(id int64) (user User, err error) {

	users, err := s.repository.GetAll()
	if err != nil {
		return user, err
	}

	for _, user := range users.Users {
		if user.Id == id {
			return user, err
		}
	}

	return user, err
}

func (s *service) NewUser(c *gin.Context) (user User, err error) {
	err = c.ShouldBindJSON(&user)
	if err != nil {
		return
	}

	usersInDatabase, err := s.repository.GetAll()
	if err != nil {
		return user, err
	}

	lastRegisteredUserindex := len(usersInDatabase.Users) - 1
	lastRegisteredUser := usersInDatabase.Users[lastRegisteredUserindex]

	// Assign a new id
	user.Id = lastRegisteredUser.Id + 1

	user, err = s.repository.Store(user.Id, user.Nombre, user.Apellido, user.Email, user.Edad, user.Altura, user.Activo, user.FechaDeCreacion)

	return user, err
}

func (s *service) FullUpdate(id int64, nombre string, apellido string, email string, edad int64, altura float64) (user User, err error) {
	user, err = s.repository.FullUpdate(id, nombre, apellido, email, edad, altura)

	return user, err
}

func (s *service) DeleteUserByID(id int64) (err error) {
	err = s.repository.DeleteUserByID(id)

	return err
}

func (s *service) UpdateUserLastName(id int64, apellido string) (user User, err error) {
	user, err = s.repository.UpdateUserLastName(id, apellido)

	return user, err
}

func (s *service) UpdateUserAge(id int64, edad int64) (user User, err error) {
	user, err = s.repository.UpdateUserAge(id, edad)

	return user, err
}

func CreateSearchedUser(availableParams []string, c *gin.Context) (searchedUser User, err error) {
	for _, param := range availableParams {
		switch param {
		case "id":
			id, err := strconv.ParseInt(c.Query("id"), 10, 64)
			if err != nil {
				return searchedUser, err
			}
			searchedUser.Id = id
		case "nombre":
			searchedUser.Nombre = c.Query("nombre")
		case "apellido":
			searchedUser.Apellido = c.Query("apellido")
		case "email":
			searchedUser.Email = c.Query("email")
		case "edad":
			edad, err := strconv.ParseInt(c.Query("edad"), 10, 64)
			if err != nil {
				return searchedUser, err
			}
			searchedUser.Edad = edad
		case "altura":
			altura, err := strconv.ParseFloat(c.Query("altura"), 64)
			if err != nil {
				return searchedUser, err
			}
			searchedUser.Altura = altura
		case "activo":
			activo, err := strconv.ParseBool(c.Query("activo"))
			if err != nil {
				return searchedUser, err
			}
			searchedUser.Activo = activo

		case "fecha_de_creacion":
			searchedUser.FechaDeCreacion = c.Query("fecha_de_creacion")
		}
	}

	return searchedUser, nil

}

func CheckAvailableParamsFromGinContext(c *gin.Context) (availableParams []string) {
	allParams := []string{"id", "nombre", "apellido", "email", "edad", "altura", "activo", "fecha_de_creacion"}
	for _, param := range allParams {
		if c.Query(param) != "" {
			availableParams = append(availableParams, param)
		}
	}

	return availableParams
}

func GetUsersWithGivenParams(availableParams []string, users Users, searchedUser User) (filteredUsers Users) {
	for _, user := range users.Users {
		areCompatible := true
		for _, attribute := range availableParams {
			switch attribute {
			case "id":
				if user.Id != searchedUser.Id {
					areCompatible = false
				}
			case "nombre":
				if user.Nombre != searchedUser.Nombre {
					areCompatible = false
				}
			case "apellido":
				if user.Apellido != searchedUser.Apellido {
					areCompatible = false
				}
			case "email":
				if user.Email != searchedUser.Email {
					areCompatible = false
				}
			case "edad":
				if user.Edad != searchedUser.Edad {
					areCompatible = false
				}
			case "altura":
				if user.Altura != searchedUser.Altura {
					areCompatible = false
				}
			case "activo":
				if user.Activo != searchedUser.Activo {
					areCompatible = false
				}
			case "fecha_de_creacion":
				if user.FechaDeCreacion != searchedUser.FechaDeCreacion {
					areCompatible = false
				}
			}
			if !areCompatible {
				break
			}
		}
		if areCompatible {
			filteredUsers.Users = append(filteredUsers.Users, user)
		}
	}

	return filteredUsers
}
