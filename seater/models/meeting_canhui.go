package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// MeetingCanhui defines a meetingcanhui
type MeetingCanhui struct {
	ID            int64     `orm:"auto;column(id)" json:"id"`
	Meetingid     int64     `json:"meetingid"`
	Rulesetupid   int64     `json:"rulesetupid"`
	Name          string    `json:"name"`
	Sexid         int64     `json:"sexid"`
	Company       string    `json:"company"`
	Duties        string    `json:"duties"`
	Phone1        string    `json:"phone1"`
	Phone2        string    `json:"phone2"`
	Contacts      string    `json:"contacts"`
	Contactsphone string    `json:"contactsphone"`
	Cardid        string    `json:"cardid"`
	Groupid       int64     `json:"groupid"`
	Partyid       int64     `json:"partyid"`
	Specialid     int64     `json:"specialid"`
	Specialorder  int64     `json:"specialorder"`
	Differentid   int64     `json:"differentid"`
	Isconvenor    int64     `json:"isconvenor"`
	convenornum   int64     `json:"convenornum"`
	Delstate      int64     `json:"delstate"`
	ImageURL      string    `orm:"column(imgurl)" json:"image_url"`
	State         int64     `json:"state"`
	Compareimg1   string    `json:"compareimg1"`
	Compareimg2   string    `json:"compareimg2"`
	Compareimg3   string    `json:"compareimg3"`
	Camera        string    `json:"camera"`
	SeatID        string    `orm:"column(seatid)" json:"seat_id"`
	SZM           string    `orm:"column(szm)" json:"szm"`
	Isleave       int64     `json:"isleave"`
	XSOrder       int64     `orm:"column(xsorder)" json:"xs_order"`
	Viproom       int64     `json:"viproom"`
	Modifytime    time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (mc *MeetingCanhui) TableName() string {
	return "ljl_meetingcanhui"
}

// String returns humble representation of Meetingcanhui
func (mc *MeetingCanhui) String() string {
	return mc.Name
}

// ListMeetingCanhuis list all Meetingcanhuis
func (m *SeaterModel) ListMeetingCanhuis(params QueryParams) (mcs []*MeetingCanhui, err error) {
	o := m.Orm()

	mcs = make([]*MeetingCanhui, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(MeetingCanhui))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &mcs)
	if err == orm.ErrNoRows {
		return mcs, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// ListSortedMeetingCanhuis list all Sorted meetingcanhuis
func (m *SeaterModel) ListSortedMeetingCanhuis(params QueryParams) (mcs []*MeetingCanhui, err error) {
	o := m.Orm()

	mcs = make([]*MeetingCanhui, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(MeetingCanhui))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("Groupid", "Specialid", "ID")
	_, err = m.PagingAll(params, qs, &mcs)
	if err == orm.ErrNoRows {
		return mcs, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetMeetingCanhui returns the MeetingCanhui
func (m *SeaterModel) GetMeetingCanhui(aID int64) (a *MeetingCanhui, err error) {
	o := m.Orm()
	a = new(MeetingCanhui)
	err = o.QueryTable(a).Filter("ID", aID).RelatedSel(2).One(a)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateMeetingCanhui creates Meetingcanhui
func (m *SeaterModel) CreateMeetingCanhui(a *MeetingCanhui) (err error) {
	if _, err = m.O().Insert(a); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateMeetingCanhui updates Meetingcanhui
func (m *SeaterModel) UpdateMeetingCanhui(a *MeetingCanhui, keys ...string) (err error) {
	_, err = m.O().Update(a, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteMeetingCanhui delete a Meetingcanhui record
func (m *SeaterModel) DeleteMeetingCanhui(a *MeetingCanhui) (err error) {
	if a == nil {
		return
	}
	_, err = m.Orm().Delete(a)
	if err != nil {
		err = fmt.Errorf("failed to delete Meetingcanhui %d %s: %s", a.ID, a.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(MeetingCanhui))
}
