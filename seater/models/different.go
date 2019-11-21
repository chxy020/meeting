package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Different defines a different
type Different struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	Delstate   int64     `json:"delstate"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (different *Different) TableName() string {
	return "ljl_meetingdifferent"
}

// String returns humble representation of different
func (different *Different) String() string {
	return different.Name
}

// ListDifferents list all differents
func (m *SeaterModel) ListDifferents(params QueryParams) (differents []*Different, err error) {
	o := m.Orm()

	differents = make([]*Different, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Different))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &differents)
	if err == orm.ErrNoRows {
		return differents, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetDifferent returns the different
func (m *SeaterModel) GetDifferent(differentID int64) (different *Different, err error) {
	o := m.Orm()
	different = new(Different)
	err = o.QueryTable(different).Filter("ID", differentID).RelatedSel(2).One(different)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateDifferent creates different
func (m *SeaterModel) CreateDifferent(different *Different) (err error) {
	if _, err = m.O().Insert(different); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateDifferent updates different
func (m *SeaterModel) UpdateDifferent(different *Different, keys ...string) (err error) {
	_, err = m.O().Update(different, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteDifferent delete a different record
func (m *SeaterModel) DeleteDifferent(different *Different) (err error) {
	if different == nil {
		return
	}
	_, err = m.Orm().Delete(different)
	if err != nil {
		err = fmt.Errorf("failed to delete different %d %s: %s", different.ID, different.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Different))
}
