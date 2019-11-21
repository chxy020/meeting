package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/juju/errors"
)

// Task defines a task will be running
type Task struct {
	ID        int64     `orm:"auto;column(id)" json:"id"`
	Type      string    `json:"type"`
	Data      string    `orm:"null;type(text)" json:"data"`
	Status    string    `json:"status"`
	ErrorMsg  string    `orm:"null;type(text)" json:"error_msg"`
	Create    time.Time `orm:"auto_now_add;type(datetime)" json:"create"`
	Update    time.Time `orm:"auto_now;type(datetime)" json:"update"`
	Execute   time.Time `orm:"null;type(datetime)" json:"execute"`
	Finish    time.Time `orm:"null;type(datetime)" json:"-"`
	Heartbeat time.Time `orm:"auto_now_add;type(datetime)" json:"-"`
}

func (task *Task) String() string {
	return fmt.Sprintf("Task(id=%d type=%s)", task.ID, task.Type)
}

// TableIndex defines indexes of task table
func (task *Task) TableIndex() [][]string {
	return [][]string{
		{"Type", "Status"},
	}
}

// list of task status
const (
	TaskStatusPending = "pending"
	TaskStatusRunning = "running"
	TaskStatusError   = "error"
	TaskStatusSuccess = "success"
)

// task types
const (
	// notification related tasks
	TaskTypeSendWSNotification = "SendWSNotification"
)

// NewTask create new task
func (m *SeaterModel) NewTask(t string, data interface{}) (task *Task, err error) {
	task = new(Task)
	task.Type = t
	var b []byte
	switch d := data.(type) {
	case *simplejson.Json:
		b, err = d.MarshalJSON()
		if err != nil {
			return nil, errors.Trace(err)
		}
		task.Data = string(b)
	case string:
		task.Data = d
	default:
		b, err = json.Marshal(&d)
		if err != nil {
			return nil, errors.Trace(err)
		}
		task.Data = string(b)
	}
	task.Status = TaskStatusPending
	task.ID, err = m.orm.Insert(task)
	if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateTask insert a task to db
func (m *SeaterModel) CreateTask(task *Task) (err error) {
	task.ID, err = m.orm.Insert(task)
	if err != nil {
		return errors.Trace(err)
	}
	return
}

// ListTasks returns task list by QueryParams
func (m *SeaterModel) ListTasks(params QueryParams, orders ...string) (tasks []*Task, err error) {
	qs := m.O().QueryTable(new(Task))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}
	if len(orders) > 0 {
		qs = qs.OrderBy(orders...)
	} else {
		qs = qs.OrderBy("ID")
	}
	qs = qs.RelatedSel(1)
	_, err = m.PagingAll(params, qs, &tasks)
	if err == orm.ErrNoRows {
		return tasks, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// UpdateTask updates task
func (m *SeaterModel) UpdateTask(t *Task, keys ...string) (err error) {
	_, err = m.orm.Update(t, keys...)
	if err != nil {
		err = errors.Trace(err)
	}
	return
}

func init() {
	registerModel(new(Task))
}
