package models

import (
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Id       int
	Name     string `orm:"unique"`
	Password string `json:"-"`
	//Todo        []*Todo `orm:"reverse(many)"` // 设置一对多的反向关系
}

func (this *User) GenPwd(password []byte) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	this.Password = string(hash)
}

func (this *User) CheckPwd(password []byte) bool {
	byteHash := []byte(this.Password)
	err := bcrypt.CompareHashAndPassword(byteHash, password)
	if err != nil {
		return false
	}
	return true
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(User), new(Todo))
}
