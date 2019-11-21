package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// AttributeType defines a attribute type
type AttributeType struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	IsHided    bool      `json:"is_hided"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (at *AttributeType) TableName() string {
	return "ljl_attribute_type"
}

// String returns humble representation of attribute type
func (at *AttributeType) String() string {
	return at.Name
}

// ListAttributeTypes list all attribute types
func (m *SeaterModel) ListAttributeTypes(params QueryParams) (ats []*AttributeType, err error) {
	o := m.Orm()

	ats = make([]*AttributeType, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(AttributeType))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &ats)
	if err == orm.ErrNoRows {
		return ats, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetAttributeType returns the attribute type
func (m *SeaterModel) GetAttributeType(atID int64) (at *AttributeType, err error) {
	o := m.Orm()
	at = new(AttributeType)
	err = o.QueryTable(at).Filter("ID", atID).RelatedSel(2).One(at)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateAttributeType creates attribute type
func (m *SeaterModel) CreateAttributeType(at *AttributeType) (err error) {
	if _, err = m.O().Insert(at); err != nil {
		fmt.Printf("%v\n", err)
		return errors.Trace(err)
	}
	return
}

// UpdateAttributeType updates attribute type
func (m *SeaterModel) UpdateAttributeType(at *AttributeType, keys ...string) (err error) {
	_, err = m.O().Update(at, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteAttributeType delete a attribute type record
func (m *SeaterModel) DeleteAttributeType(at *AttributeType) (err error) {
	if at == nil {
		return
	}
	_, err = m.Orm().Delete(at)
	if err != nil {
		err = fmt.Errorf("failed to delete attribute type %d %s: %s", at.ID, at.Name, err.Error())
		return
	}
	return
}

// action types
const (
	AttributeTypeGroup     = "Group"
	AttributeTypeParty     = "Party"
	AttributeTypeRole      = "Role"
	AttributeTypeSex       = "Sex"
	AttributeTypeConvenor  = "Convenor"
	AttributeTypeDifferent = "Different"
)

func init() {
	registerModel(new(AttributeType))
}
