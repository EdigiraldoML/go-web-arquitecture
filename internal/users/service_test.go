package users

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myDbFullUpdate struct {
	Users      []User
	ReadCalled bool
}

func (db *myDbFullUpdate) Read(data interface{}) (err error) {
	db.ReadCalled = true

	data2 := data.(*Users)
	data2.Users = db.Users

	return nil
}
func (db *myDbFullUpdate) Write(data interface{}) (err error) {

	data2 := data.(*Users)
	db.Users = data2.Users

	return nil
}

func TestFullUpdate(t *testing.T) {
	var idUserToUpdateLastName int64 = 1
	expectedUser := User{Id: idUserToUpdateLastName, Nombre: "user name after update", Apellido: "user last name after update", Email: "userAfterUpdate@email.com", Edad: 36, Altura: 1.59, Activo: true, FechaDeCreacion: "13/12/2021"}
	userBefore := User{Id: idUserToUpdateLastName, Nombre: "user name before update", Apellido: "user last name before update", Email: "userBeforeUpdate@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}

	db := &myDbFullUpdate{
		Users:      []User{userBefore},
		ReadCalled: false,
	}

	repo := CreateRepository(db)
	service := CreateService(repo)

	user, err := service.FullUpdate(idUserToUpdateLastName, expectedUser.Nombre, expectedUser.Apellido, expectedUser.Email, expectedUser.Edad, expectedUser.Altura)
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)

	usedRepo := (repo).(*repository)
	usedDb := (usedRepo.db).(*myDbFullUpdate)
	assert.True(t, usedDb.ReadCalled)
}

type myDbDeleteUserByID struct {
	Users []User
}

func (db *myDbDeleteUserByID) Read(data interface{}) (err error) {
	data2 := data.(*Users)
	data2.Users = db.Users

	return nil
}
func (db *myDbDeleteUserByID) Write(data interface{}) (err error) {

	data2 := data.(*Users)
	db.Users = data2.Users

	return nil
}

func TestDeleteUserByID(t *testing.T) {
	var idNonExistentUserToDelete int64 = 2
	expectedError := errors.New("el usuario no fue encontrado")

	var idExistentUserToDelete int64 = 1
	user := User{Id: idExistentUserToDelete, Nombre: "user name", Apellido: "user last name", Email: "user@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}

	db := &myDbDeleteUserByID{
		Users: []User{user},
	}

	repo := CreateRepository(db)
	service := CreateService(repo)

	// Testea el caso de eliminar un usuario inexistente
	err := service.DeleteUserByID(idNonExistentUserToDelete)
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}

	// Testea el caso de eliminar un usuario existente
	err = service.DeleteUserByID(idExistentUserToDelete)
	assert.Nil(t, err)

	usedRepo := (repo).(*repository)
	usedDb := (usedRepo.db).(*myDbDeleteUserByID)
	assert.Empty(t, usedDb.Users)
}
