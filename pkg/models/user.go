package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name		string 		`db:"name" json:"Name"`
	Posts		[]Post		`json:"posts,omitempty"`
}

func (db *DB) GetUsers() (users []*User, err error) {
	err = db.Find(&users).Error

	return
}

func (db *DB) GetUserByID(id int64) (user *User, err error) {
	err = db.First(&user, id).Error

	return
}

func (db *DB) InsertUser(user *User, tx *gorm.DB) (err error) {
	result := tx.Create(&user)

	return result.Error
}

func (db *DB) UpdateUser(user *User, tx *gorm.DB) (err error) {
	result := tx.Save(&user)

	return result.Error
}

func (db *DB) DeleteUserByID(user *User, tx *gorm.DB) (err error) {
	result := tx.Delete(&user)

	return result.Error
}