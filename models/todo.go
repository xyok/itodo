package models

import (
	"github.com/astaxie/beego/orm"
	"log"
	"reflect"
	"strings"
	"time"
)

const (
	Delete = 0
	Undo   = 1
	Doing  = 2
	Down   = 3
)

type Todo struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
	Status   int    `json:"status" orm:"column(status)"`
	UserId   int    `json:"user_id" orm:"column(user_id);size(11)"`
	//User     *User `orm:"rel(fk)"`
}

func AddTodo(m *Todo, user_id int) (t *Todo, err error) {
	o := orm.NewOrm()

	CreatedAt := time.Now().UTC().Unix()
	UpdatedAt := CreatedAt

	todo := Todo{
		Title:    m.Title,
		CreateAt: CreatedAt,
		UpdateAt: UpdatedAt,
		Status:   Undo,
		UserId:   user_id,
	}

	_, err = o.Insert(&todo)
	if err == nil {
		return &todo, err
	}

	return nil, err
}

func GetTodoById(id int64) (m *Todo, err error) {
	o := orm.NewOrm()
	m = &Todo{Id: id}
	err = o.Read(m)
	if err == nil {
		return m, err
	}
	return nil, err
}

func GetAllTodo(
	query map[string]string,
	fields []string,
	sortby []string,
	offset int64,
	limit int64,
	userId int) (ml []interface{}, totalCount int64, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(Todo))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}

	if len(sortby) != 0 {
		qs = qs.OrderBy(sortby...)
	}

	var l []Todo

	totalCount, err = qs.Filter("UserId", userId).Count()
	if _, err = qs.Filter("UserId", userId).Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, totalCount, nil
	}
	return nil, 0, err

}

func UpdateTodoById(m *Todo) (err error) {
	o := orm.NewOrm()
	v := Todo{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		m.UpdateAt = time.Now().UTC().Unix()
		if _, err = o.Update(m, "Title", "Status", "UpdateAt"); err != nil {
			log.Println(err)
		}
	}
	return
}

func DeleteTodo(id int64) (err error) {
	o := orm.NewOrm()
	v := Todo{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		if _, err = o.Delete(&Todo{Id: id}); err != nil {
			log.Println(err)
		}
	}
	return
}
