package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/juju/errors"
)

// Party defines a party
type Party struct {
	ID         int64     `orm:"auto;column(id)" json:"id"`
	Name       string    `json:"name"`
	Delstate   int64     `json:"delstate"`
	Modifytime time.Time `orm:"auto_now;type(datetime)" json:"modifytime"`
}

// TableName set table name
func (party *Party) TableName() string {
	return "ljl_party"
}

// String returns humble representation of party
func (party *Party) String() string {
	return party.Name
}

// ListPartys list all partys
func (m *SeaterModel) ListPartys(params QueryParams) (partys []*Party, err error) {
	o := m.Orm()

	partys = make([]*Party, 0, PagingDefaultLimit)

	qs := o.QueryTable(new(Party))
	if params != nil {
		qs = qs.SetCond(params.Condition())
	}

	qs = qs.OrderBy("-ID")
	_, err = m.PagingAll(params, qs, &partys)
	if err == orm.ErrNoRows {
		return partys, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}

	return
}

// GetParty returns the party
func (m *SeaterModel) GetParty(partyID int64) (party *Party, err error) {
	o := m.Orm()
	party = new(Party)
	err = o.QueryTable(party).Filter("ID", partyID).RelatedSel(2).One(party)
	if err == orm.ErrNoRows {
		return nil, nil
	} else if err != nil {
		err = errors.Trace(err)
		return
	}
	return
}

// CreateParty creates party
func (m *SeaterModel) CreateParty(party *Party) (err error) {
	if _, err = m.O().Insert(party); err != nil {
		return errors.Trace(err)
	}
	return
}

// UpdateParty updates party
func (m *SeaterModel) UpdateParty(party *Party, keys ...string) (err error) {
	_, err = m.O().Update(party, keys...)
	if err != nil {
		return errors.Trace(err)
	}

	return
}

// DeleteParty delete a party record
func (m *SeaterModel) DeleteParty(party *Party) (err error) {
	if party == nil {
		return
	}
	_, err = m.Orm().Delete(party)
	if err != nil {
		err = fmt.Errorf("failed to delete party %d %s: %s", party.ID, party.Name, err.Error())
		return
	}
	return
}

func init() {
	registerModel(new(Party))
}
