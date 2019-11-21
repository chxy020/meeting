package models

import (
	"fmt"
	"reflect"
	"time"

	"github.com/astaxie/beego/orm"
)

// BasePaging base duration paging
type BasePaging struct {
	length int
	first  time.Time
	last   time.Time
}

// Len return the lengthof slice
func (p *BasePaging) Len() int {
	return p.length
}

// Times return the first and last time
func (p *BasePaging) Times() (first time.Time, last time.Time) {
	first = p.first
	last = p.last
	if !first.Before(last) {
		return last, first
	}
	return
}

// Init return new base paging
func (p *BasePaging) Init(slice interface{}, pagingType string) {
	s := reflect.ValueOf(slice).Elem()
	sType := s.Kind()
	if sType != reflect.Slice {
		err := fmt.Errorf("want type slice, but got %s", sType)
		panic(err)
	}

	num := s.Len()
	if num == 0 {
		return
	}
	p.length = num

	if pagingType == PagingTypeDuration {
		first := s.Index(0).Elem().FieldByName("Create")
		last := s.Index(num - 1).Elem().FieldByName("Create")
		if !(first.IsValid() && last.IsValid()) {
			return
		}
		p.first = first.Interface().(time.Time)
		p.last = last.Interface().(time.Time)
	}
}

// PagingParams store parameters used in paging
type PagingParams QueryParams

// PagingType return the type of paging
func (p PagingParams) PagingType() string {
	if len(p) == 0 {
		return NoPaging
	}

	_, beginOK := p[PagingDurationBegin]
	_, endOK := p[PagingDurationEnd]
	_, periodOK := p[PagingPeriod]
	if beginOK || endOK || periodOK {
		return PagingTypeDuration
	}

	if _, ok := p[PagingOffset]; ok {
		return PagingTypeOffset
	}
	return PagingTypeDefault
}

func (p PagingParams) setResult(result QueryParams, paging PagingInterface, totalQs orm.QuerySeter) (
	err error) {

	t := p.PagingType()
	params := QueryParams(p)

	limit, ok := params.GetInt64(PagingLimit)
	if !ok {
		limit = PagingDefaultLimit
	}
	result.SetInt64(PagingLimit, limit)

	switch t {
	case PagingTypeDuration:
		if paging == nil {
			return
		}
		p.setDurationResult(result, paging)
	case PagingTypeOffset, PagingTypeDefault:
		if paging == nil || totalQs == nil {
			return
		}
		var total int64
		total, err = totalQs.Count()
		if err != nil {
			return
		}
		p.setOffsetResult(result, int64(paging.Len()), total)
	}
	return
}

func (p PagingParams) setDurationResult(result QueryParams, paging PagingInterface) {
	params := QueryParams(p)

	durationBegin, beginOK := params.GetTime(PagingDurationBegin)
	durationEnd, endOK := params.GetTime(PagingDurationEnd)

	length := paging.Len()
	first, last := paging.Times()

	if !beginOK && !endOK {
		// client wasn't passing any duration paging argument, so not set paging result
		return
	}

	if !beginOK && length > 0 {
		// client wasn't passing duration_begin, get it from the earliest record
		durationBegin = first
		result.SetTime(PagingDurationBegin, first)
		beginOK = true
	} else {
		result.SetTime(PagingDurationBegin, durationBegin)
	}

	if !endOK && length > 0 {
		// client wasn't passing duration_end, get it from the latest record
		durationEnd = last
		result.SetTime(PagingDurationEnd, last)
		endOK = true
	} else {
		result.SetTime(PagingDurationEnd, durationEnd)
	}

	if beginOK && endOK {
		duration := durationEnd.Sub(durationBegin)
		result.SetInt64(PagingDuration, int64(duration.Seconds()))
	} else {
		result.SetInt64(PagingDuration, int64(0))
	}
}

func (p PagingParams) setOffsetResult(result QueryParams, count, total int64) {
	offset, offsetOK := QueryParams(p).GetInt64(PagingOffset)
	if !offsetOK {
		offset = 0
	} else if total < offset {
		offset = total
	}
	result.SetInt64(PagingOffset, offset)
	result.SetInt64(PagingCount, count)
	result.SetInt64(PagingTotalCount, total)
}

// PagingAll return query result and set paging
func (m *SeaterModel) PagingAll(params QueryParams, qs orm.QuerySeter, container interface{},
	pagings ...PagingInterface) (count int64, err error) {

	countQs := qs
	if params == nil {
		params = NewQueryParams()
	}
	p := params.GetPaging()
	t := p.PagingType()

	var offset int64
	limit, ok := params.GetInt64(PagingLimit)
	if !ok {
		limit = PagingUnlimit
	}

	switch t {
	case PagingTypeDuration:
		durationBegin, ok := params.GetTime(PagingDurationBegin)
		if ok {
			qs = qs.Filter("Create__gte", durationBegin)
		}
		durationEnd, ok := params.GetTime(PagingDurationEnd)
		if ok {
			qs = qs.Filter("Create__lt", durationEnd)
		}
	case PagingTypeOffset:
		offset, _ = params.GetInt64(PagingOffset)
	}
	qs = qs.Limit(limit, offset)

	count, err = qs.All(container)
	if err != nil {
		return 0, err
	}

	if t == NoPaging {
		return
	}

	var paging PagingInterface
	if len(pagings) > 0 {
		paging = pagings[0]
	} else {
		paging = new(BasePaging)
	}

	paging.Init(container, t)
	err = p.setResult(m.PagingResult, paging, countQs)
	if err != nil {
		return
	}
	return
}
