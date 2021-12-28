package users

import (
	"errors"

	"github.com/EdigiraldoML/go-web-arquitecture/pkg/store"
)

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	Id              int64   `json:"id"`
	Nombre          string  `json:"nombre" binding:"required"`
	Apellido        string  `json:"apellido" binding:"required"`
	Email           string  `json:"email" binding:"required"`
	Edad            int64   `json:"edad" binding:"required"`
	Altura          float64 `json:"altura" binding:"required"`
	Activo          bool    `json:"activo" binding:"required"`
	FechaDeCreacion string  `json:"fecha_de_creacion" binding:"required"`
}

type Repository interface {
	GetAll() (users *Users, err error)
	Store(id int64, nombre string, apellido string, email string, edad int64, altura float64, activo bool, fecha_de_creacion string) (user User, err error)
	FullUpdate(id int64, nombre string, apellido string, email string, edad int64, altura float64) (user User, err error)
	DeleteUserByID(id int64) (err error)
	UpdateUserLastName(id int64, apellido string) (user User, err error)
	UpdateUserAge(id int64, edad int64) (user User, err error)
}

type repository struct {
	db store.Store
}

func CreateRepository(db store.Store) Repository {
	newRepository := &repository{
		db: db,
	}

	return newRepository
}

func (r *repository) GetAll() (users *Users, err error) {

	users = &Users{}

	err = r.db.Read(users)

	return users, err
}

func (r *repository) Store(id int64, nombre string, apellido string, email string, edad int64, altura float64, activo bool, fecha_de_creacion string) (user User, err error) {

	var usuarios Users
	err = r.db.Read(&usuarios)
	if err != nil {
		return user, err
	}

	user = User{id, nombre, apellido, email, edad, altura, activo, fecha_de_creacion}

	usuarios.Users = append(usuarios.Users, user)

	err = r.db.Write(&usuarios)

	return user, err
}

func (r *repository) FullUpdate(id int64, nombre string, apellido string, email string, edad int64, altura float64) (user User, err error) {
	users, err := r.GetAll()
	if err != nil {
		return user, err
	}

	userToUpdate, err := GetUserById(id, users)
	if err != nil {
		return user, err
	}

	userToUpdate.Nombre = nombre
	userToUpdate.Apellido = apellido
	userToUpdate.Email = email
	userToUpdate.Edad = edad
	userToUpdate.Altura = altura

	user = *userToUpdate

	err = r.db.Write(&users)

	return user, err
}

func (r *repository) DeleteUserByID(id int64) (err error) {
	users, err := r.GetAll()
	if err != nil {
		return err
	}

	userIndexToDelete, err := GetUserIndexById(id, users)
	if err != nil {
		return err
	}

	users.Users[userIndexToDelete].Activo = false
	users.Users = append(users.Users[:userIndexToDelete], users.Users[userIndexToDelete+1:]...)

	err = r.db.Write(&users)

	return err
}

func (r *repository) UpdateUserLastName(id int64, apellido string) (user User, err error) {
	users, err := r.GetAll()
	if err != nil {
		return user, err
	}

	ptrUser, err := GetUserById(id, users)
	if err != nil {
		return user, err
	}

	ptrUser.Apellido = apellido
	user = *ptrUser

	err = r.db.Write(&users)

	return user, err
}

func (r *repository) UpdateUserAge(id int64, edad int64) (user User, err error) {
	users, err := r.GetAll()
	if err != nil {
		return user, err
	}

	ptrUser, err := GetUserById(id, users)
	if err != nil {
		return user, err
	}

	ptrUser.Edad = edad
	user = *ptrUser

	err = r.db.Write(&users)

	return user, err
}

func GetUserById(id int64, users *Users) (user *User, err error) {
	for idx, user := range users.Users {
		if user.Id == id {
			return &users.Users[idx], err
		}
	}

	err = errors.New("el usuario no fue encontrado")

	return user, err
}

func GetUserIndexById(id int64, users *Users) (index int64, err error) {
	for idx, user := range users.Users {
		if user.Id == id {
			return int64(idx), err
		}
	}

	err = errors.New("el usuario no fue encontrado")

	return index, err
}
