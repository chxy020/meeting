package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// QueryLoadRelatedModel defines whether to load related model
const QueryLoadRelatedModel = "load_related"

// Orm query filters
const (
	QueryFilterKey    = "__filter__"
	QueryFilterAnd    = "and"
	QueryFilterAndNot = "and_not"
	QueryFilterOr     = "or"
	QueryFilterOrNot  = "or_not"
)

type queryFilter struct {
	filters []*queryFilterItem
}

type queryFilterItem struct {
	key    string
	filter string
}

// QueryParams store parameters used in database query
type QueryParams map[string]interface{}

// NewQueryParams return a newly initialized QueryParams
func NewQueryParams() QueryParams {
	return make(QueryParams)
}

// GetInt64 return parameter indexed by key and convert it to int64
func (p QueryParams) GetInt64(key string) (value int64, ok bool) {
	t, ok := p[key]
	if !ok {
		return
	}
	value, ok = t.(int64)
	return
}

// SetString store a string value into QueryParams
func (p QueryParams) SetString(key string, value string, filters ...string) QueryParams {
	p.addFilter(key, filters...)

	var t interface{} = value
	p[key] = t
	return p
}

// SetInt64 store an int64 value into QueryParams
func (p QueryParams) SetInt64(key string, value int64) QueryParams {
	var t interface{} = value
	p[key] = t
	return p
}

// SetTime store a Time value into QueryParams
func (p QueryParams) SetTime(key string, value time.Time, filters ...string) QueryParams {
	p.addFilter(key, filters...)

	var t interface{} = value
	p[key] = t
	return p
}

// GetTime return parameter indexed by key and convert it to Time
func (p QueryParams) GetTime(key string) (value time.Time, ok bool) {
	t, ok := p[key]
	if !ok {
		return
	}
	value, ok = t.(time.Time)
	return
}

// SetTimeDuration store a time.Duration value into QueryParams
func (p QueryParams) SetTimeDuration(key string, value time.Duration, filters ...string) QueryParams {
	p.addFilter(key, filters...)

	var t interface{} = value
	p[key] = t
	return p
}

// GetTimeDuration return parameter indexed by key and convert it ot time.Duration
func (p QueryParams) GetTimeDuration(key string) (value time.Duration, ok bool) {
	t, ok := p[key]
	if !ok {
		return
	}
	value, ok = t.(time.Duration)
	return
}


// Copy copy all parameters from given QueryParams
func (p QueryParams) Copy(q QueryParams) {
	for k, v := range q {
		p[k] = v
	}
}

func (p QueryParams) getQueryFilter() (qf *queryFilter, ok bool) {
	t, ok := p[QueryFilterKey]
	if !ok {
		return
	}
	qf, ok = t.(*queryFilter)
	return
}

func (p QueryParams) addFilter(key string, filters ...string) {
	if 0 == len(filters) {
		return
	}
	filter := filters[0]

	qf, ok := p.getQueryFilter()
	if !ok {
		qf = &queryFilter{make([]*queryFilterItem, 0)}
		p[QueryFilterKey] = interface{}(qf)
	}
	qf.filters = append(qf.filters, &queryFilterItem{key, filter})
}

// Condition return new condition object according to query params
func (p QueryParams) Condition() (cond *orm.Condition) {
	cond = orm.NewCondition()

	params := NewQueryParams()
	params.Copy(p)
	for _, key := range PagingReqKeys {
		delete(params, key)
	}
	delete(params, QueryLoadRelatedModel)

	qf, ok := params.getQueryFilter()
	if ok && 0 == len(qf.filters) {
		delete(params, QueryFilterKey)
		ok = false
	}
	if !ok {
		for key := range params {
			cond = cond.And(key, params.GetValue(key))
		}
		return
	}

	for _, f := range qf.filters {
		switch f.filter {
		case QueryFilterAnd:
			cond = cond.And(f.key, params.GetValue(f.key))
		case QueryFilterAndNot:
			cond = cond.AndNot(f.key, params.GetValue(f.key))
		case QueryFilterOr:
			cond = cond.Or(f.key, params.GetValue(f.key))
		case QueryFilterOrNot:
			cond = cond.OrNot(f.key, params.GetValue(f.key))
		default:
			continue
		}
		delete(params, f.key)
	}
	delete(params, QueryFilterKey)
	for key := range params {
		cond = cond.And(key, params.GetValue(key))
	}

	return
}

// GetValue return parameter indexed by key
func (p QueryParams) GetValue(key string) interface{} {
	return p[key]
}

// SetValue store an int64 value into QueryParams
func (p QueryParams) SetValue(key string, value interface{}, filters ...string) QueryParams {
	p.addFilter(key, filters...)
	p[key] = value
	return p
}

// GetPaging return paging params
func (p QueryParams) GetPaging() PagingParams {
	params := NewQueryParams()
	for _, key := range PagingReqKeys {
		if v := p.GetValue(key); v != nil {
			params[key] = v
		}
	}
	return PagingParams(params)
}
