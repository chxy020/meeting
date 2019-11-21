package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Group defines a group
type Group struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	Delstate   int64     `json:"delstate"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (group *Group) TableName() string {
	return "ljl_meetinggroup"
}

// String returns humble representation of group
func (group *Group) String() string {
	return group.Name
}

// ListGroups list all groups
func (m *SeaterModel) ListGroups(params QueryParams) (groups []*Group, err error) {
	o := m.Orm()

	groups = make([]*Group, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Group))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &groups)
	if err == orm.ErrNoRows {
		return groups, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetGroup returns the group
func (m *SeaterModel) GetGroup(groupID int64) (group *Group, err error) {
	o := m.Orm()
	group = new(Group)
	err = o.QueryTable(group).Filter("ID", groupID).RelatedSel(2).One(group)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateGroup creates group
func (m *SeaterModel) CreateGroup(group *Group) (err error) {
	if _, err = m.O().Insert(group); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateGroup updates group
func (m *SeaterModel) UpdateGroup(group *Group, keys ...string) (err error) {
	_, err = m.O().Update(group, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteGroup delete a group record
func (m *SeaterModel) DeleteGroup(group *Group) (err error) {
	if group == nil {
		return
	}
	_, err = m.Orm().Delete(group)
	if err != nil {
		err = fmt.Errorf("failed to delete group %d %s: %s", group.ID, group.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Group))
}
