package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Member defines a member
type Member struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Username   string    `json:"username"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (member *Member) TableName() string {
	return "ljl_member"
}

// String returns humble representation of member
func (member *Member) String() string {
	return member.Username
}

// ListMembers list all members
func (m *SeaterModel) ListMembers(params QueryParams) (members []*Member, err error) {
	o := m.Orm()

	members = make([]*Member, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Member))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &members)
	if err == orm.ErrNoRows {
		return members, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetMember returns the Member
func (m *SeaterModel) GetMember(memberID int64) (member *Member, err error) {
	o := m.Orm()
	member = new(Member)
	err = o.QueryTable(member).Filter("ID", memberID).RelatedSel(2).One(member)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateMember creates member
func (m *SeaterModel) CreateMember(member *Member) (err error) {
	if _, err = m.O().Insert(member); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateMember updates member
func (m *SeaterModel) UpdateMember(member *Member, keys ...string) (err error) {
	_, err = m.O().Update(member, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteMember delete a member record
func (m *SeaterModel) DeleteMember(member *Member) (err error) {
	if member == nil {
		return
	}
	_, err = m.Orm().Delete(member)
	if err != nil {
		err = fmt.Errorf("failed to delete member %d %s: %s", member.ID, member.Username, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Member))
}
