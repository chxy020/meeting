package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Meeting defines a meeting
type Meeting struct {
	ID            int64     `orm:"auto;column(id)" json:"id"`
	Name          string    `json:"name"`
	Delstate      int64     `json:"delstate"`
	Roomid        int64     `json:"roomid"`
	Reciverlistid int64     `json:"reciverlistid"`
	Ruleid        int64     `json:"ruleid"`
	Ischeck       int64     `json:"ischeck"`
	Isface        int64     `json:"isface"`
	Isgrouplist   int64     `json:"isgrouplist"`
	Isstay        int64     `json:"isstay"`
	Iscard        int64     `json:"iscard"`
	Isvote        int64     `json:"isvote"`
	Meetingtime   string    `json:"meetingtime"`
	Address       string    `json:"address"`
	Memo          string    `json:"memo"`
	Modifytime    time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (meeting *Meeting) TableName() string {
	return "ljl_meeting"
}

// String returns humble representation of meeting
func (meeting *Meeting) String() string {
	return meeting.Name
}

// ListMeetings list all meetings
func (m *SeaterModel) ListMeetings(params QueryParams) (meetings []*Meeting, err error) {
	o := m.Orm()

	meetings = make([]*Meeting, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Meeting))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &meetings)
	if err == orm.ErrNoRows {
		return meetings, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetMeeting returns the Meeting
func (m *SeaterModel) GetMeeting(meetingID int64) (meeting *Meeting, err error) {
	o := m.Orm()
	meeting = new(Meeting)
	err = o.QueryTable(meeting).Filter("ID", meetingID).RelatedSel(2).One(meeting)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateMeeting creates meeting
func (m *SeaterModel) CreateMeeting(meeting *Meeting) (err error) {
	if _, err = m.O().Insert(meeting); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateMeeting updates meeting
func (m *SeaterModel) UpdateMeeting(meeting *Meeting, keys ...string) (err error) {
	_, err = m.O().Update(meeting, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteMeeting delete a meeting record
func (m *SeaterModel) DeleteMeeting(meeting *Meeting) (err error) {
	if meeting == nil {
		return
	}
	_, err = m.Orm().Delete(meeting)
	if err != nil {
		err = fmt.Errorf("failed to delete meeting %d %s: %s", meeting.ID, meeting.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Meeting))
}
