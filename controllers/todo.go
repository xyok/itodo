package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"itodo/models"
	"strconv"
	"strings"
)

type BaseController struct {
	UserId int
	beego.Controller
}

func (c *BaseController) auth() {
	tokenString := c.Ctx.Request.Header.Get("token")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(beego.AppConfig.String("secret")), nil
	})

	if err!=nil{
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = errToken
		c.ServeJSON()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		c.UserId = int(claims["id"].(float64))
	} else {
		c.Ctx.ResponseWriter.WriteHeader(401)
		c.Data["json"] = errToken
		c.ServeJSON()
		return
	}

}

func (c *BaseController) Prepare() {
	c.auth()
}

//  TodoController operations for Todo
type TodoController struct {
	BaseController
}

// URLMapping ...
func (c *TodoController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Todo
// @Param	body		body 	models.Todo	true		"body for Todo content"
// @Success 201 {int} models.Todo
// @Failure 403 body is empty
// @router / [post]
func (c *TodoController) Post() {
	var v models.Todo
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	user_id := c.UserId
	if t, err := models.AddTodo(&v, user_id); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = &t
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Todo by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Todo
// @Failure 403 :id is empty
// @router /:id [get]
func (c *TodoController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetTodoById(id)
	if err != nil {
		c.Data["json"] = err404
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Todo
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Todo
// @Failure 403
// @router / [get]
func (c *TodoController) GetAll() {
	var fields []string
	var sortby []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}

	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, total, err := models.GetAllTodo(query, fields, sortby, offset, limit, c.UserId)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = &Resp{
			Data: map[string]interface{}{
				"data":  l,
				"count": total,
			},
		}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Todo
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Todo	true		"body for Todo content"
// @Success 200 {object} models.Todo
// @Failure 403 :id is not int
// @router /:id [put]
func (c *TodoController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Todo{Id: int64(id)}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateTodoById(&v); err == nil {
		c.Data["json"] = successReturn
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Todo
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *TodoController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteTodo(id); err == nil {
		c.Data["json"] = successReturn
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
