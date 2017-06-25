package common

import (
	"crypto/md5"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255";not null;unique_index`
	Password string `gorm:"size:255";not null`
	Email    string `gorm:"type:varchar(100);unique_index"`
}

func Migrate() error {
	db, err := OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{})
	return nil
}
func Add(user, password, email string) error {
	db, err := OpenDb()
	if err != nil {
		return err
	}
	defer db.Close()
	db.Create(&User{
		Name:     user,
		Password: HashPasswd(MD5(password)),
		Email:    email,
	})
	return nil
}
func MD5(msg string) string {
	h := md5.New()
	io.WriteString(h, msg)
	return fmt.Sprintf("%x", h.Sum(nil))
}
func HashPasswd(pass string) string {
	return MD5(pass + "whyengineer")
}
func OpenDb() (*gorm.DB, error) {
	return gorm.Open("mysql", "root:71451085a@/db1?charset=utf8&parseTime=True&loc=Local")
}
func CheckName(username string) bool {
	db, _ := OpenDb()
	defer db.Close()
	return db.Where("name = ?", username).First(&User{}).RecordNotFound()
}
func Authenticate(username, password string) (bool, error) {
	var user User
	db, err := OpenDb()
	if err != nil {
		return false, err
	}
	defer db.Close()
	if db.Where("name = ?", username).First(&user).RecordNotFound() {
		return false, nil
	} else {
		if user.Password == HashPasswd(MD5(password)) {
			return true, nil
		} else {
			return false, nil
		}
	}
}
