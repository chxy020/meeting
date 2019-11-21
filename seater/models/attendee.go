package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Attendee defines a attendee
type Attendee struct {
	ID            int64      `orm:"auto;column(id)" json:"id"`
	Meeting       *Meeting   `orm:"rel(fk);null;on_delete(set_null)" json:"meeting"`
	Rulesetup     *Rulesetup `orm:"rel(fk);null;on_delete(set_null)" json:"rulesetup"`
	Name          string     `json:"name"`
	Company       string     `json:"company"`
	Duties        string     `json:"duties"`
	Phone1        string     `json:"phone1"`
	Phone2        string     `json:"phone2"`
	Contacts      string     `json:"contacts"`
	ContactsPhone string     `json:"contacts_phone"`
	CardID        string     `orm:"column(card_id)" json:"card_id"`
	Delstate      int64      `json:"delstate"`
	ImageURL      string     `orm:"column(image_url)" json:"image_url"`
	State         int64      `json:"state"`
	Compareimg1   string     `json:"compareimg1"`
	Compareimg2   string     `json:"compareimg2"`
	Compareimg3   string     `json:"compareimg3"`
	Camera        string     `json:"camera"`
	SeatID        string     `orm:"column(seat_id)" json:"seat_id"`
	SZM           string     `orm:"column(szm)" json:"szm"`
	IsLeft        int64      `json:"is_left"`
	XSOrder       int64      `orm:"column(xs_order)" json:"xs_order"`
	VIPRoom       *VIPRoom   `orm:"rel(fk);null;on_delete(set_null)" json:"vip_room"`
	Attributes    string     `json:"attributes"`
	Modifytime    time.Time  `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (a *Attendee) TableName() string {
	return "ljl_attendee"
}

// String returns humble representation of attendee
func (a *Attendee) String() string {
	return a.Name
}

// ListAttendees list all attendees
func (m *SeaterModel) ListAttendees(params QueryParams) (as []*Attendee, err error) {
	o := m.Orm()

	as = make([]*Attendee, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Attendee))
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

// GetAttendee returns the Attendee
func (m *SeaterModel) GetAttendee(aID int64) (a *Attendee, err error) {
	o := m.Orm()
	a = new(Attendee)
	err = o.QueryTable(a).Filter("ID", aID).RelatedSel(2).One(a)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateAttendee creates attendee
func (m *SeaterModel) CreateAttendee(a *Attendee) (err error) {
	if _, err = m.O().Insert(a); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateAttendee updates attendee
func (m *SeaterModel) UpdateAttendee(a *Attendee, keys ...string) (err error) {
	_, err = m.O().Update(a, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteAttendee delete a attendee record
func (m *SeaterModel) DeleteAttendee(a *Attendee) (err error) {
	if a == nil {
		return
	}
	_, err = m.Orm().Delete(a)
	if err != nil {
		err = fmt.Errorf("failed to delete attendee %d %s: %s", a.ID, a.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Attendee))
}
