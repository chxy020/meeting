package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Rulecanhui defines a rulecanhui
type Rulecanhui struct {
	ID          int64     `orm:"auto;column(id)" json:"id"`
	Planid      string    `json:"planid"`
	Planname    string    `json:"planname"`
	Meetingid   int64     `json:"meetingid"`
	Groupid     string    `json:"groupid"`
	Specialid   string    `json:"specialid"`
	Rulesetupid int64     `json:"rulesetupid"`
	Delstate    int64     `json:"delstate"`
	Modifytime  time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (rulecanhui *Rulecanhui) TableName() string {
	return "ljl_rulecanhui"
}

// String returns humble representation of rulecanhui
func (rulecanhui *Rulecanhui) String() string {
	return rulecanhui.Planname
}

// ListRulecanhuis list all rulecanhuis
func (m *SeaterModel) ListRulecanhuis(params QueryParams) (rulecanhuis []*Rulecanhui, err error) {
	o := m.Orm()

	rulecanhuis = make([]*Rulecanhui, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Rulecanhui))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &rulecanhuis)
	if err == orm.ErrNoRows {
		return rulecanhuis, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetRulecanhui returns the rulecanhui
func (m *SeaterModel) GetRulecanhui(rulecanhuiID int64) (rulecanhui *Rulecanhui, err error) {
	o := m.Orm()
	rulecanhui = new(Rulecanhui)
	err = o.QueryTable(rulecanhui).Filter("ID", rulecanhuiID).RelatedSel(2).One(rulecanhui)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateRulecanhui creates rulecanhui
func (m *SeaterModel) CreateRulecanhui(rulecanhui *Rulecanhui) (err error) {
	if _, err = m.O().Insert(rulecanhui); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateRulecanhui updates rulecanhui
func (m *SeaterModel) UpdateRulecanhui(rulecanhui *Rulecanhui, keys ...string) (err error) {
	_, err = m.O().Update(rulecanhui, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteRulecanhui delete a rulecanhui record
func (m *SeaterModel) DeleteRulecanhui(rulecanhui *Rulecanhui) (err error) {
	if rulecanhui == nil {
		return
	}
	_, err = m.Orm().Delete(rulecanhui)
	if err != nil {
		err = fmt.Errorf("failed to delete rulecanhui %d %s: %s", rulecanhui.ID, rulecanhui.Planname, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Rulecanhui))
}
