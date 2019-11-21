package controllers

import (
	"github.com/bitly/go-simplejson"

	"seater/models"
)

func (c *SeaterController) pagingParams() models.QueryParams {
	params := models.NewQueryParams()
	durationBegin, ok := c.getTime(models.PagingDurationBegin)
	if ok {
		params.SetTime(models.PagingDurationBegin, durationBegin)
	}
	durationEnd, ok := c.getTime(models.PagingDurationEnd)
	if ok {
		params.SetTime(models.PagingDurationEnd, durationEnd)
	}

	offset, ok := c.getInt64(models.PagingOffset)
	if ok {
		if offset < 0 {
			c.BadRequestf("offset should be greater than or equal to 0")
		}
		params.SetInt64(models.PagingOffset, offset)
	}

	limit, ok := c.getInt64(models.PagingLimit)
	if ok {
		if limit <= 0 && limit != -1 {
			c.BadRequestf("limit should be greater than 0 or equal to -1")
		}
	} else {
		limit = models.PagingDefaultLimit
	}
	params.SetInt64(models.PagingLimit, limit)

	return params
}

func (c *SeaterController) getPagingResult() (j *simplejson.Json) {
	m := c.model
	j = simplejson.New()
	var pagingResult models.QueryParams
	if len(c.pagingResult) > 0 {
		pagingResult = c.pagingResult
	} else if len(m.PagingResult) > 0 {
		pagingResult = m.PagingResult
	} else {
		return nil
	}

	limit, limitOK := pagingResult.GetInt64(models.PagingLimit)
	if limitOK {
		j.Set(models.PagingLimit, limit)
	}

	t := models.PagingParams(c.pagingParams()).PagingType()
	switch t {
	case models.PagingTypeDuration:
		durationBegin, _ := pagingResult.GetTime(models.PagingDurationBegin)
		durationEnd, _ := pagingResult.GetTime(models.PagingDurationEnd)
		j.Set(models.PagingDurationBegin, durationBegin.Format(TimestampLayout))
		j.Set(models.PagingDurationEnd, durationEnd.Format(TimestampLayout))
		if duration, ok := pagingResult.GetInt64(models.PagingDuration); ok {
			j.Set(models.PagingDuration, duration)
		}
		return
	case models.PagingTypeOffset, models.PagingTypeDefault:
		offset, _ := pagingResult.GetInt64(models.PagingOffset)
		count, _ := pagingResult.GetInt64(models.PagingCount)
		total, _ := pagingResult.GetInt64(models.PagingTotalCount)
		j.Set(models.PagingOffset, offset)
		j.Set(models.PagingCount, count)
		j.Set(models.PagingTotalCount, total)
		return
	}
	return nil
}

func (c *SeaterController) setDurationPagingResult(params models.QueryParams, container interface{}) {
	durationBegin, beginOK := params.GetTime(models.PagingDurationBegin)
	durationEnd, endOK := params.GetTime(models.PagingDurationEnd)
	if !beginOK && !endOK {
		// client wasn't passing any duration paging argument, so not set paging result
		return
	}

	paging := new(models.BasePaging)
	paging.Init(container, models.PagingTypeDuration)
	length := paging.Len()
	first, last := paging.Times()
	result := c.pagingResult

	if !beginOK && length > 0 {
		// client wasn't passing duration_begin, get it from the earliest record
		durationBegin = first
		result.SetTime(models.PagingDurationBegin, first)
		beginOK = true
	} else {
		result.SetTime(models.PagingDurationBegin, durationBegin)
	}

	if !endOK && length > 0 {
		// client wasn't passing duration_end, get it from the latest record
		durationEnd = last
		result.SetTime(models.PagingDurationEnd, last)
		endOK = true
	} else {
		result.SetTime(models.PagingDurationEnd, durationEnd)
	}

	if beginOK && endOK {
		duration := durationEnd.Sub(durationBegin)
		result.SetInt64(models.PagingDuration, int64(duration.Seconds()))
	} else {
		result.SetInt64(models.PagingDuration, int64(0))
	}
}
