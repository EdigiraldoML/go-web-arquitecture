package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myDbGetAll struct {
}

func (db *myDbGetAll) Read(data interface{}) (err error) {
	user1 := User{Id: 1, Nombre: "user1 name", Apellido: "user1 last name", Email: "user1@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}
	user2 := User{Id: 2, Nombre: "user2 name", Apellido: "user2 last name", Email: "user2@email.com", Edad: 35, Altura: 1.82, Activo: true, FechaDeCreacion: "13/12/2021"}

	data2 := data.(*Users)

	data2.Users = []User{user1, user2}

	return nil
}
func (db *myDbGetAll) Write(data interface{}) (err error) {
	return nil
}

func TestGetAll(t *testing.T) {
	user1 := User{Id: 1, Nombre: "user1 name", Apellido: "user1 last name", Email: "user1@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}
	user2 := User{Id: 2, Nombre: "user2 name", Apellido: "user2 last name", Email: "user2@email.com", Edad: 35, Altura: 1.82, Activo: true, FechaDeCreacion: "13/12/2021"}
	expectedUsers := &Users{Users: []User{user1, user2}}
	db := &myDbGetAll{}

	repo := &repository{
		db: db,
	}

	users, err := repo.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, expectedUsers, users)

}

type myDbUpdateLastName struct {
	Users      []User
	ReadCalled bool
}

func (db *myDbUpdateLastName) Read(data interface{}) (err error) {
	db.ReadCalled = true

	data2 := data.(*Users)
	data2.Users = db.Users

	return nil
}
func (db *myDbUpdateLastName) Write(data interface{}) (err error) {

	data2 := data.(*Users)
	db.Users = data2.Users

	return nil
}

func TestUpdateUserLastName(t *testing.T) {
	var idUserToUpdateLastName int64 = 1
	expectedUser := User{Id: idUserToUpdateLastName, Nombre: "user name", Apellido: "user last name after update", Email: "user@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}
	user1 := User{Id: idUserToUpdateLastName, Nombre: "user name", Apellido: "user last name before update", Email: "user@email.com", Edad: 35, Altura: 1.58, Activo: true, FechaDeCreacion: "13/12/2021"}

	db := &myDbUpdateLastName{
		Users:      []User{user1},
		ReadCalled: false,
	}

	repo := repository{
		db: db,
	}

	user, err := repo.UpdateUserLastName(idUserToUpdateLastName, "user last name after update")
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)

	usedDb := repo.db.(*myDbUpdateLastName)
	assert.True(t, usedDb.ReadCalled)
}
