package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/juju/errors"
	// define the default db
	_ "github.com/lib/pq"

	"seater/database"
)

// Resource Status
const (
	StatusWarning = "warning"
	StatusError   = "error"

	StatusStop  = "stop"
	StatusStart = "start"
	StatusDoing = "doing"

	StatusDeleting = "deleting"
	StatusOffline  = "offline"
	StatusSuccess  = "success"
	StatusFailure  = "failure"
	StatusFinished = "finished"
)

// registerModel registers models
func registerModel(vals ...interface{}) {
	orm.RegisterModel(vals...)
}

// SeaterModel defines model context
type SeaterModel struct {
	orm          orm.Ormer
	PagingResult QueryParams
}

func newModel(args ...interface{}) (m *SeaterModel, err error) {
	var o orm.Ormer
	for _, arg := range args {
		if vo, ok := arg.(orm.Ormer); ok {
			o = vo
		}
	}
	m = new(SeaterModel)
	if o != nil {
		m.orm = o
	} else {
		m.orm = database.NewOrm()
	}

	m.PagingResult = NewQueryParams()
	return
}

// CountFunc function count resources
type CountFunc func(QueryParams) (int64, error)

// ModelFunc defines func type that return a new SeaterModel
type ModelFunc func(...interface{}) (*SeaterModel, error)

// NewModel is the function return a new SeaterModel
var NewModel ModelFunc = newModel

// SetOrm sets orm
func (m *SeaterModel) SetOrm(orm orm.Ormer) {
	m.orm = orm
}

// Orm gets orm of the model
func (m *SeaterModel) Orm() orm.Ormer {
	return m.orm
}

// O gets orm of the model
func (m *SeaterModel) O() orm.Ormer {
	return m.Orm()
}

// Argument names of paging query
const (
	PagingDefaultLimit  int64 = 10
	PagingUnlimit       int64 = -1
	PagingLimit               = "limit"
	PagingOffset              = "offset"
	PagingCount               = "count"
	PagingTotalCount          = "total_count"
	PagingDurationBegin       = "duration_begin"
	PagingDurationEnd         = "duration_end"
	PagingDuration            = "duration"
	PagingPeriod              = "period"
)

// PagingReqKeys arguments of paging request
var PagingReqKeys = []string{PagingDurationBegin, PagingDurationEnd, PagingOffset, PagingLimit, PagingPeriod}

// Paging types
const (
	PagingTypeOffset   = "offset"
	PagingTypeDuration = "duration"
	PagingTypeDefault  = "default"
	NoPaging           = "none"
)

// Paging defines paging object
type Paging struct {
	Offset        int64     `json:"offset"`
	Limit         int64     `json:"limit"`
	Count         int64     `json:"count"`
	TotalCount    int64     `json:"total_count"`
	DurationBegin time.Time `json:"duration_begin"`
	DurationEnd   time.Time `json:"duration_end"`
	Duration      int64     `json:"duration"`
	Period        string    `json:"period"`
}

// PagingInterface defines methods used in pagination
type PagingInterface interface {
	Init(interface{}, string)
	Len() int
	Times() (first, last time.Time)
}

// Rollback rollback the transaction
func (m *SeaterModel) Rollback() error {
	err := m.orm.Rollback()
	if err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Commit commits the transaction
func (m *SeaterModel) Commit() error {
	err := m.orm.Commit()
	if err != nil {
		return errors.Annotate(err, "commit transaction")
	}
	return nil
}

// Begin starts the transaction for the model
func (m *SeaterModel) Begin() error {
	err := m.orm.Begin()
	if err != nil {
		return errors.Annotate(err, "begin transaction")
	}
	return nil
}
