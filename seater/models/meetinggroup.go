package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Meetinggroup defines a group
type Meetinggroup struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	Delstate   int64     `json:"delstate"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (group *Meetinggroup) TableName() string {
	return "ljl_group"
}

// String returns humble representation of group
func (group *Meetinggroup) String() string {
	return group.Name
}

// ListMeetinggroups list all groups
func (m *SeaterModel) ListMeetinggroups(params QueryParams) (groups []*Meetinggroup, err error) {
	o := m.Orm()

	groups = make([]*Meetinggroup, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Meetinggroup))
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

// GetMeetinggroup returns the group
func (m *SeaterModel) GetMeetinggroup(aID int64) (group *Meetinggroup, err error) {
	o := m.Orm()
	group = new(Meetinggroup)
	err = o.QueryTable(group).Filter("ID", aID).RelatedSel(2).One(group)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateMeetinggroup creates group
func (m *SeaterModel) CreateMeetinggroup(group *Meetinggroup) (err error) {
	if _, err = m.O().Insert(group); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateMeetinggroup updates group
func (m *SeaterModel) UpdateMeetinggroup(group *Meetinggroup, keys ...string) (err error) {
	_, err = m.O().Update(group, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteMeetinggroup delete a group record
func (m *SeaterModel) DeleteMeetinggroup(group *Meetinggroup) (err error) {
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
	registerModel(new(Meetinggroup))
}
