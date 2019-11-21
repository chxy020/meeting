package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Attribute defines a attribute
type Attribute struct {
	ID            int64          `orm:"auto;column(id)" json:"id"`
	Member        *Member        `orm:"rel(fk);null;on_delete(set_null)" json:"member"`
	AttributeType *AttributeType `orm:"rel(fk);null;on_delete(set_null)" json:"attribute_type"`
	Parent        *Attribute     `orm:"rel(fk);null;on_delete(set_null)" json:"parent"`
	Name          string         `json:"name"`
	Content       string         `json:"content"`
	Label         string         `json:"label"`
	Modifytime    time.Time      `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (a *Attribute) TableName() string {
	return "ljl_attribute"
}

// String returns humble representation of attribute
func (a *Attribute) String() string {
	return a.Name + ":" + a.Content
}

// ListAttributeTypes list all attribute
func (m *SeaterModel) ListAttributes(params QueryParams) (as []*Attribute, err error) {
	o := m.Orm()

	as = make([]*Attribute, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Attribute))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &as)
	if err == orm.ErrNoRows {
		return as, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetAttribute returns the attribute
func (m *SeaterModel) GetAttribute(aID int64) (a *Attribute, err error) {
	o := m.Orm()
	a = new(Attribute)
	err = o.QueryTable(a).Filter("ID", aID).RelatedSel(2).One(a)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateAttribute creates attribute
func (m *SeaterModel) CreateAttribute(a *Attribute) (err error) {
	if _, err = m.O().Insert(a); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateAttribute updates attribute
func (m *SeaterModel) UpdateAttribute(a *Attribute, keys ...string) (err error) {
	_, err = m.O().Update(a, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteAttribute delete a attribute record
func (m *SeaterModel) DeleteAttribute(a *Attribute) (err error) {
	if a == nil {
		return
	}
	_, err = m.Orm().Delete(a)
	if err != nil {
		err = fmt.Errorf("failed to delete attribute %d %s: %s", a.ID, a.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Attribute))
}
