package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Special defines a special
type Special struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	Delstate   int64     `json:"delstate"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (special *Special) TableName() string {
	return "ljl_meetingspecial"
}

// String returns humble representation of special
func (special *Special) String() string {
	return special.Name
}

// ListSpecials list all specials
func (m *SeaterModel) ListSpecials(params QueryParams) (specials []*Special, err error) {
	o := m.Orm()

	specials = make([]*Special, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Special))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &specials)
	if err == orm.ErrNoRows {
		return specials, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetSpecial returns the special
func (m *SeaterModel) GetSpecial(specialID int64) (special *Special, err error) {
	o := m.Orm()
	special = new(Special)
	err = o.QueryTable(special).Filter("ID", specialID).RelatedSel(2).One(special)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateSpecial creates special
func (m *SeaterModel) CreateSpecial(special *Special) (err error) {
	if _, err = m.O().Insert(special); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateSpecial updates special
func (m *SeaterModel) UpdateSpecial(special *Special, keys ...string) (err error) {
	_, err = m.O().Update(special, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteSpecial delete a special record
func (m *SeaterModel) DeleteSpecial(special *Special) (err error) {
	if special == nil {
		return
	}
	_, err = m.Orm().Delete(special)
	if err != nil {
		err = fmt.Errorf("failed to delete special %d %s: %s", special.ID, special.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Special))
}
