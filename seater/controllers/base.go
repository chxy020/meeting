package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
	"github.com/juju/errors"

	"seater/models"
	"seater/schema"
)

// Define controller constants
const (
	TimestampLayout = time.RFC3339

	PeriodDay  = "day"
	PeriodHour = "hour"
	PeriodWeek = "week"
)

type msgBody struct {
	Msg string `json:"msg"`
}

// SeaterController defines the base controller
type SeaterController struct {
	beego.Controller
	model        *models.SeaterModel
	orm          orm.Ormer
	deferrers    []deferrer
	errs         []error
	pagingResult models.QueryParams
}

// Prepare prepares controller context
func (c *SeaterController) Prepare() {
	model, err := models.NewModel()
	if err != nil {
		c.TraceServerError(errors.Annotatef(err, "failed to init model"))
	}
	if err = model.Begin(); err != nil {
		c.TraceServerError(errors.Annotatef(err, "failed to begin database transaction"))
	}
	c.model = model
	c.orm = model.Orm()
	c.pagingResult = models.NewQueryParams()
}

// Finish ends transaction
func (c *SeaterController) Finish() {
	defer c.execDeferrers()

	err := c.endTransaction()
	if err != nil {
		c.TraceServerError(errors.Annotatef(err, "failed to end transaction"))
	}
}

// M returns the model object
func (c *SeaterController) M() *models.SeaterModel {
	return c.model
}

type deferrer func() error

func (c *SeaterController) deferExec(f deferrer) {
	c.deferrers = append(c.deferrers, f)
}

// Code sets the response status
func (c *SeaterController) Code(code int) {
	c.Ctx.Output.SetStatus(code)
}

func (c *SeaterController) execDeferrers() {
	var err error
	for i := len(c.deferrers) - 1; i >= 0; i-- {
		err = c.deferrers[i]()
		if err != nil {
			c.errs = append(c.errs, err)
		}
	}
}

func (c *SeaterController) traceJSONAbort(err error, code int, args ...string) {
	c.jsonAbort(code, args...)
}

// jsonAbort trace and abort error
func (c *SeaterController) jsonAbort(code int, args ...string) {
	defer c.execDeferrers()

	c.Header("Content-Type", "application/json; charset=utf-8")
	var msg string
	if len(args) == 0 || args[0] == "" {
		switch code {
		case 400:
			msg = "Bad Request"
		case 401:
			msg = "Unauthorized"
		case 404:
			msg = "Resource Not Found"
		case 409:
			msg = "Conflict"
		case 500:
			msg = "Server Error"
		default:
			msg = ""
		}
	} else {
		msg = args[0]
	}
	c.addError(fmt.Errorf(msg))
	err := c.endTransaction()
	if err != nil {
		code = 500
		msg = "Server Error"
	}

	body, err := json.Marshal(msgBody{Msg: msg})
	if err != nil {
		c.CustomAbort(500, `{"msg": "Unknown Error"}`)
	}
	c.CustomAbort(code, string(body))
}

// BadRequestf returns bad request response with formatted message
func (c *SeaterController) BadRequestf(format string, args ...interface{}) {
	c.TraceBadRequestf(nil, format, args...)
}

// TraceBadRequestf traces error and returns bad request response with formatted message
func (c *SeaterController) TraceBadRequestf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.traceJSONAbort(nil, 400, msg)
}

// TraceServerError traces error and returns server error
func (c *SeaterController) TraceServerError(err error) {
	c.traceJSONAbort(err, 500)
}

// Forbiddenf returns forbidden response with formatted message
func (c *SeaterController) Forbiddenf(format string, args ...interface{}) {
	c.TraceForbiddenf(nil, format, args...)
}

// TraceForbiddenf traces error and returns forbidden response with formatted message
func (c *SeaterController) TraceForbiddenf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.traceJSONAbort(err, 403, msg)
}

// NotFoundf returns not found response with formatted message
func (c *SeaterController) NotFoundf(format string, args ...interface{}) {
	c.TraceNotFoundf(nil, format, args...)
}

// TraceNotFoundf traces error and returns not found response with formatted message
func (c *SeaterController) TraceNotFoundf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.traceJSONAbort(err, 404, msg)
}

// Conflictf returns conflict response with formatted message
func (c *SeaterController) Conflictf(format string, args ...interface{}) {
	c.TraceConflictf(nil, format, args...)
}

// TraceConflictf traces error and returns conflict response with formatted message
func (c *SeaterController) TraceConflictf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.traceJSONAbort(err, 409, msg)
}

// Unauthorizedf returns authorized response with formatted message
func (c *SeaterController) Unauthorizedf(format string, args ...interface{}) {
	c.TraceUnauthorizedf(nil, format, args...)
}

// TraceUnauthorizedf traces error and returns authorized reponse with formatted message
func (c *SeaterController) TraceUnauthorizedf(err error, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.traceJSONAbort(err, 401, msg)
}

func (c *SeaterController) addError(err error) {
	c.errs = append(c.errs, err)
}

// jsonResp serves json response
func (c *SeaterController) jsonResp(data interface{}) {
	if obj, ok := data.(*simplejson.Json); ok {
		data = obj.Interface()
	}

	paging := c.getPagingResult()

	if paging != nil {
		bytes, err := json.Marshal(data)
		if err != nil {
			err = errors.Annotatef(err, "failed to marshal resp interface")
			c.TraceServerError(err)
		}

		j, err := simplejson.NewJson(bytes)
		if err != nil {
			err = errors.Annotatef(err, "failed to unmarshal resp bytes")
			c.TraceServerError(err)
		}
		j.Set("paging", paging)
		data = j.Interface()
	}
	c.Data["json"] = data

	c.ServeJSON()
}

// OK response 200 OK with json data
func (c *SeaterController) OK(data interface{}) {
	c.Code(200)
	c.jsonResp(data)
}

// Accepted response an asynchronous resource
func (c *SeaterController) Accepted(data interface{}) {
	c.Code(202)
	c.jsonResp(data)
}

// Created response an asynchronous resource
func (c *SeaterController) Created(data interface{}) {
	c.Code(201)
	c.jsonResp(data)
}

// NoContent responses with code 204
func (c *SeaterController) NoContent(code ...int) {
	if len(code) > 0 {
		c.Code(code[0])
	} else {
		c.Code(204)
	}
	c.Ctx.Output.Body([]byte(""))
}

// Validate validates with json schema
func (c *SeaterController) Validate(sche string, document ...string) {
	var doc string
	if len(document) > 0 {
		doc = document[0]
	} else {
		doc = string(c.Ctx.Input.RequestBody)
		if len(doc) == 0 {
			c.BadRequestf("request body is empty")
		}
	}
	_, err := simplejson.NewJson([]byte(doc))
	if err != nil {
		c.BadRequestf("invalid json format")
	}
	result, err := schema.Validate(sche, doc)
	if err != nil {
		c.TraceServerError(errors.Annotatef(err, "invalid schema"))
	}
	if !result.Valid() {
		s := "invalid parameters:\n"
		var e interface{}
		for _, err := range result.Errors() {
			s += fmt.Sprintf("%s\n", err)
			e = err
		}
		c.BadRequestf("%s", e)
	}
}

func (c *SeaterController) getInt64(key string, defs ...int64) (v int64, ok bool) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		val, err := strconv.ParseInt(strv, 10, 64)
		if err != nil {
			c.BadRequestf("invalid int64 argument %s: %s", key, strv)
		}
		return val, true
	}
	return
}

func (c *SeaterController) getString(key string, defs ...string) (v string, ok bool) {
	if v = c.Ctx.Input.Query(key); v != "" {
		return v, true
	}
	if len(defs) > 0 {
		return defs[0], false
	}
	return "", false
}

// getTime return input as an time and the existence of the input
func (c *SeaterController) getTime(key string, defs ...time.Time) (v time.Time, ok bool) {
	if strv := c.Ctx.Input.Query(key); strv != "" {
		val, err := time.Parse(TimestampLayout, strv)
		if err != nil {
			c.BadRequestf("invalid time argument %s: %s", key, strv)
		}
		return val, true
	} else if len(defs) > 0 {
		v = defs[0]
		return
	}
	return
}

// Header get or set a header if value is provided
func (c *SeaterController) Header(key string, value ...interface{}) string {
	if len(value) == 0 {
		return c.Ctx.Input.Header(key)
	}
	retval := fmt.Sprintf("%v", value[0])
	c.Ctx.Output.Header(key, retval)
	return retval
}

func (c *SeaterController) endTransaction() (err error) {
	if c.model == nil {
		return
	}
	rollback := false
	if len(c.errs) > 0 {
		rollback = true
	}
	if rollback {
		err = c.model.Rollback()
		if err != nil {
			panic(fmt.Sprintf("failed to rollback transaction: %v", err))
		}
	} else {
		err = c.model.Commit()
		if err != nil {
			panic(fmt.Sprintf("failed to commit transaction: %v", err))
		}
	}
	return
}

func (c *SeaterController) parseJSONBody(keys ...string) (v *simplejson.Json) {
	v, err := simplejson.NewJson(c.Ctx.Input.RequestBody)
	if err != nil {
		c.BadRequestf("invalid json format")
	}
	if len(keys) > 0 {
		for _, k := range keys {
			_, ok := v.CheckGet(k)
			if !ok {
				c.BadRequestf("Bad Request")
			} else {
				v = v.Get(k)
			}
		}
	}
	return
}

// UnmarshalJSONBody unmarshal request json body
func (c *SeaterController) UnmarshalJSONBody(v interface{}, keys ...string) {
	var bytes []byte
	var err error

	if len(keys) > 0 {
		j := c.parseJSONBody(keys...)
		bytes, err = j.MarshalJSON()
		if err != nil {
			err = errors.Annotate(err, "failed to unmarshal json")
			c.TraceServerError(err)
		}
	} else {
		bytes = c.Ctx.Input.RequestBody
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		c.BadRequestf("invalid request body")
	}
}

// UserInfo defines session value
type UserInfo struct {
	UserID     int64  `json:"user_id"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

func (c *SeaterController) getURLParam(key string) string {
	return c.Ctx.Input.Param(key)
}

func (c *SeaterController) getURLID(name string) int64 {
	id, err := strconv.ParseInt(c.getURLParam(name), 10, 64)
	if err != nil {
		c.BadRequestf("invalid id")
	}
	return id
}

// CreateTask create task
func (c *SeaterController) CreateTask(t string, data *simplejson.Json) (task *models.Task, err error) {
	if task, err = c.model.NewTask(t, data); err != nil {
		err = errors.Annotatef(err, "failed to create task %s", t)
		return
	}
	return
}
