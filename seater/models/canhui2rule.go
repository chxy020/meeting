package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Canhui2rule defines a canhui2rule
type Canhui2rule struct {
	ID             int64     `orm:"auto;column(id)" json:"id"`
	Name           string    `json:"name"`
	Meetingid      int64     `json:"meetingid"`
	Ruletemplateid int64     `json:"ruletemplateid"`
	Isused         int64     `json:"isused"`
	Delstate       int64     `json:"delstate"`
	Modifytime     time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (canhui2rule *Canhui2rule) TableName() string {
	return "ljl_canhui2rule"
}

// String returns humble representation of canhui2rule
func (canhui2rule *Canhui2rule) String() string {
	return canhui2rule.Name
}

// ListCanhui2rules list all canhui2rules
func (m *SeaterModel) ListCanhui2rules(params QueryParams) (canhui2rules []*Canhui2rule, err error) {
	o := m.Orm()

	canhui2rules = make([]*Canhui2rule, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Canhui2rule))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &canhui2rules)
	if err == orm.ErrNoRows {
		return canhui2rules, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetCanhui2rule returns the canhui2rule
func (m *SeaterModel) GetCanhui2rule(canhui2ruleID int64) (canhui2rule *Canhui2rule, err error) {
	o := m.Orm()
	canhui2rule = new(Canhui2rule)
	err = o.QueryTable(canhui2rule).Filter("ID", canhui2ruleID).RelatedSel(2).One(canhui2rule)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateCanhui2rule creates canhui2rule
func (m *SeaterModel) CreateCanhui2rule(canhui2rule *Canhui2rule) (err error) {
	if _, err = m.O().Insert(canhui2rule); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateCanhui2rule updates canhui2rule
func (m *SeaterModel) UpdateCanhui2rule(canhui2rule *Canhui2rule, keys ...string) (err error) {
	_, err = m.O().Update(canhui2rule, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteCanhui2rule delete a canhui2rule record
func (m *SeaterModel) DeleteCanhui2rule(canhui2rule *Canhui2rule) (err error) {
	if canhui2rule == nil {
		return
	}
	_, err = m.Orm().Delete(canhui2rule)
	if err != nil {
		err = fmt.Errorf("failed to delete canhui2rule %d %s: %s", canhui2rule.ID, canhui2rule.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Canhui2rule))
}
