package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   	string   `db:"title" json:"title"`
}

func (db *DB) GetPosts() (posts []Post, err error) {
	//db.Limit(10).Offset(5).Find(&users)
	err = db.Find(&posts).Error

	return
}

func (db *DB) GetPostByID(id int64) (post *Post, err error) {
	err = db.First(&post, id).Error

	return
}

func (db *DB) InsertPost(post *Post, tx *gorm.DB) (err error) {
	result := tx.Create(&post)

	return result.Error
}

func (db *DB) UpdatePost(post *Post, tx *gorm.DB) (err error) {
	result := tx.Save(&post)

	return result.Error
}

func (db *DB) DeletePostByID(post *Post, tx *gorm.DB) (err error) {
	result := tx.Delete(&post)

	return result.Error
}