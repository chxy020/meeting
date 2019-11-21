package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Rulesetup defines a rulesetup
type Rulesetup struct {
	ID             int64     `orm:"auto;column(id)" json:"id"`
	Rulename       string    `json:"rulename"`
	Priorityindex  int64     `json:"priorityindex"`
	Modifierid     string    `json:"modifierid"`
	Delstate       int64     `json:"delstate"`
	Stauts         int64     `json:"stauts"`
	Rulezone       string    `json:"rulezone"`
	Bgcolor        string    `json:"bgcolor"`
	Ruletemplateid int64     `json:"ruletemplateid"`
	Seatsnum       string    `json:"seatsnum"`
	Ruleid         string    `json:"ruleid"`
	Roomtemplateid string    `json:"roomtemplateid"`
	Groupid        string    `json:"groupid"`
	Specialid      string    `json:"specialid"`
	Isincrement    int64     `json:"isincrement"`
	Linkurl        string    `json:"linkurl"`
	Modifytime     time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (rulesetup *Rulesetup) TableName() string {
	return "ljl_rulesetup"
}

// String returns humble representation of rulesetup
func (rulsetup *Rulesetup) String() string {
	return rulsetup.Rulename
}

// ListRulesetups list all rulesetups
func (m *SeaterModel) ListRulesetups(params QueryParams) (rulesetups []*Rulesetup, err error) {
	o := m.Orm()

	rulesetups = make([]*Rulesetup, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Rulesetup))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("Priorityindex")
	_, err = m.PagingAll(params, qs, &rulesetups)
	if err == orm.ErrNoRows {
		return rulesetups, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetRulesetup returns the Rulesetup
func (m *SeaterModel) GetRulesetup(rulesetupID int64) (rulesetup *Rulesetup, err error) {
	o := m.Orm()
	rulesetup = new(Rulesetup)
	err = o.QueryTable(rulesetup).Filter("ID", rulesetupID).RelatedSel(2).One(rulesetup)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateRulesetup creates rulesetup
func (m *SeaterModel) CreateRulesetup(rulesetup *Rulesetup) (err error) {
	if _, err = m.O().Insert(rulesetup); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateRulesetup updates rulesetup
func (m *SeaterModel) UpdateRulesetup(rulesetup *Rulesetup, keys ...string) (err error) {
	_, err = m.O().Update(rulesetup, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteRulesetup delete a rulesetup record
func (m *SeaterModel) DeleteRulesetup(rulesetup *Rulesetup) (err error) {
	if rulesetup == nil {
		return
	}
	_, err = m.Orm().Delete(rulesetup)
	if err != nil {
		err = fmt.Errorf("failed to delete rulesetup %d %s: %s", rulesetup.ID, rulesetup.Rulename, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Rulesetup))
}
