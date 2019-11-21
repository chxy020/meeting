package controllers

import (
	"fmt"
	"seater/models"

	"github.com/juju/errors"
)

// CheckGetMeeting checks and gets meeting
func (c *SeaterController) CheckGetMeeting(meetingID int64, isNotFound ...bool) (meeting *models.Meeting) {
	meeting, err := c.model.GetMeeting(meetingID)
	if err != nil {
		c.TraceServerError(errors.Annotatef(err, "failed to get record of meeting %d", meetingID))
	} else if meeting == nil {
		msg := fmt.Sprintf("meeting not found: %v", meetingID)
		if len(isNotFound) > 0 && isNotFound[0] {
			c.NotFoundf(msg)
		} else {
			c.BadRequestf(msg)
		}
	}
	return meeting
}

// CheckGetRulesetup checks and gets rulesetup
func (c *SeaterController) CheckGetRulesetup(rulesetupID int64, isNotFound ...bool) (rulsetup *models.Rulesetup) {
	rulesetup, err := c.model.GetRulesetup(rulesetupID)
	if err != nil {
		c.TraceServerError(errors.Annotatef(err, "failed to get record of rulesetup %d", rulesetupID))
	} else if rulesetup == nil {
		msg := fmt.Sprintf("rulesetup not found: %v", rulesetupID)
		if len(isNotFound) > 0 && isNotFound[0] {
			c.NotFoundf(msg)
		} else {
			c.BadRequestf(msg)
		}
	}
	return rulesetup
}
