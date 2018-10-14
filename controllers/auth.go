package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"itodo/models"
	"log"
	"time"
)

type LoginController struct {
	beego.Controller
}

type AuthController struct {
	beego.Controller
}

type PostAuth struct {
	Name     string
	Password string
}

func (c *LoginController) Post() {
	//name := c.GetString("name")
	//password := c.GetString("password")

	o := orm.NewOrm()
	o.Using("default")

	var post PostAuth

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &post); err != nil {
		c.Data["json"] = errInputData
		c.ServeJSON()
		return
	}

	user := models.User{Name: post.Name}
	err := o.Read(&user, "Name")
	if err == orm.ErrNoRows {
		resp := Resp{403, "", ""}
		c.Data["json"] = &resp
	} else {
		if user.CheckPwd([]byte(post.Password)) {

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":   user.Id,
				"name": user.Name,
				"exp":  time.Now().Add(time.Hour * 24).Unix(),
			})
			secret := beego.AppConfig.String("secret")
			tokenString, _ := token.SignedString([]byte(secret))
			resp := Resp{200, "", tokenString}

			c.Ctx.SetSecureCookie(secret, "name", user.Name)

			c.Data["json"] = &resp
		} else {
			resp := Resp{405, "wrong password", ""}
			c.Data["json"] = &resp
		}
	}

	c.ServeJSON()

}

func (c *AuthController) SignUp() {
	name := c.GetString("name")
	password := c.GetString("password")

	o := orm.NewOrm()
	o.Using("default")

	u := models.User{Name: name}

	err := o.Read(&u, "Name")

	if err == orm.ErrNoRows {
		u.GenPwd([]byte(password))

		_, err := o.Insert(&u)

		if err == nil {
			c.Data["json"] = &u
		} else {
			log.Println(err)
			resp := Resp{403, "", ""}
			c.Data["json"] = resp
		}

	} else {
		resp := Resp{501, "has exists", &u}
		c.Data["json"] = resp
	}

	c.ServeJSON()

}

func (c *AuthController) ValideToken() {
	//token 验证
	tokenString := c.Ctx.Request.Header.Get("token")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(beego.AppConfig.String("secret")), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		resp := successReturn
		c.Data["json"] = resp
	} else {
		//log.Println(err)
		resp := Resp{401, "token is not valide", ""}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}
