package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"
	"time"
)

// VIPRoom defines a viproom
type VIPRoom struct {
	ID           int64         `orm:"auto;column(id)" json:"id"`
	RoomID int64 `orm:"column(roomid)" json:"roomid"`
	Name      string     `json:"name"`
	Seatsnum int64 `json:"seatsnum"`
	Imgurl string `json:"imgurl"`
	Delstate int64 `json:"delstate"`
	Modifytime       time.Time     `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (romt *VIPRoom) TableName() string {
	return "ljl_viproom"
}

// String returns humble representation of viproom
func (room *VIPRoom) String() string {
	return room.Name
}

// ListVIPRooms list all viprooms
func (m *SeaterModel) ListVIPRooms(params QueryParams) (rooms []*VIPRoom, err error) {
	o := m.Orm()

	rooms = make([]*VIPRoom, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(VIPRoom))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &rooms)
	if err == orm.ErrNoRows {
		return rooms, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetVIPRoom returns the viproom
func (m *SeaterModel) GetVIPRoom(roomID int64) (room *VIPRoom, err error) {
	o := m.Orm()
	room = new(VIPRoom)
	err = o.QueryTable(room).Filter("ID", roomID).RelatedSel(2).One(room)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateVIPRoom creates viproom
func (m *SeaterModel) CreateVIPRoom(room *VIPRoom) (err error) {
	if _, err = m.O().Insert(room); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateVIPRoom updates viproom
func (m *SeaterModel) UpdateVIPRoom(room *VIPRoom, keys ...string) (err error) {
	_, err = m.O().Update(room, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteVIPRoom delete a viproom record
func (m *SeaterModel) DeleteVIPRoom(room *VIPRoom) (err error) {
	if room == nil {
		return
	}
	_, err = m.Orm().Delete(room)
	if err != nil {
		err = fmt.Errorf("failed to delete viproom %d %s: %s", room.ID, room.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(VIPRoom))
}
