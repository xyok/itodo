package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"itodo/models"
	"strconv"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var (
	successReturn = &Resp{200, "ok", "ok"}
	err404        = &Resp{404, "404", "object not found"}
	errInputData  = &Resp{400, "数据输入错误", "客户端参数错误"}
	errToken = &Resp{401, "wrong token", ""}
)

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *UserController) Get() {
	uid, _ := strconv.Atoi(c.Ctx.Input.Param(":uid"))
	o := orm.NewOrm()
	o.Using("default")

	user := models.User{Id: uid}
	err := o.Read(&user)

	if err == orm.ErrNoRows {
		fmt.Println("查询不到")
	} else if err == orm.ErrMissPK {
		fmt.Println("miss pk")
	}

	c.Data["json"] = &user
	c.ServeJSON()
}

func (c *UserController) All() {
	o := orm.NewOrm()
	o.Using("default")
	qs := o.QueryTable("user")

	var users []*models.User

	qs.All(&users)

	c.Data["json"] = &users
	c.ServeJSON()
}
