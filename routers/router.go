package routers

import (
	"github.com/astaxie/beego"
	"itodo/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/signup", &controllers.AuthController{}, "post:SignUp")
	beego.Router("/valide", &controllers.AuthController{}, "post:ValideToken")
	beego.Router("/user", &controllers.UserController{}, "get:All")
	beego.Router("/user/?:uid:int", &controllers.UserController{})
	beego.Router("/todo", &controllers.TodoController{}, "get:GetAll;post:Post")
	beego.Router("/todo/:id:int", &controllers.TodoController{}, "get:GetOne;delete:Delete;put:Put")
}
