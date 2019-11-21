package controllers

import (
	"fmt"
	"seater/models"
	"sort"
	"strconv"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/juju/errors"
)

// MeetingController meeting management
type MeetingController struct {
	SeaterController
}

type MeetingShowReq struct {
	Append    int64 `form:"append"`
	MeetingID int64 `form:"meeting_id"`
}

type Canhui struct {
	ID        int64  `json:"id"`
	Attender  string `json:"attender"`
	Seatid    string `json:"seatid"`
	Groupid   int64  `json:"groupid"`
	Specialid int64  `json:"specialid"`
	Bgcolor   string `json:"bgcolor"`
}

type MeetingShowResp struct {
	Canhuis []*Canhui `json:"attendees"`
}

type Seat struct {
	Row    int  `json:"row"`
	Column int  `json:"column"`
	IsUsed bool `json:"is_used"`
	Score  int  `json:"score"`
}

// MeetingShow API
// @Title MeetingSHow
// @Description sort and show meeting seats
// @Param body body controllers.MeetingShowReq true "meeting info"
// @Success 200 {object} controllers.MeetingResp OK
// @Failures 404, 500
// @router /sort [Post]
func (c *MeetingController) MeetingShow() {
	req := new(MeetingShowReq)
	if err := c.ParseForm(req); err != nil {
		c.TraceServerError(errors.Trace(err))
	}

	m := c.M()
	meetingID := req.MeetingID
	//meeting := c.CheckGetMeeting(meetingID, true)

	params := models.NewQueryParams()
	params.SetInt64("Meetingid", meetingID)
	canhuiList, _ := m.ListSortedMeetingCanhuis(params)
	canhuis := make([]*Canhui, 0, len(canhuiList))

	params = models.NewQueryParams()
	params.SetInt64("Isused", 1)
	canhui2rules, _ := m.ListCanhui2rules(params)
	ruletemplateid := canhui2rules[0].Ruletemplateid
	planid := canhui2rules[0].ID

	params = models.NewQueryParams()
	params.SetInt64("Ruletemplateid", ruletemplateid)
	rulesetups, _ := m.ListRulesetups(params)

	ruleSeatMap := make(map[int64][]*Seat)
	//ruleSeatNumMap := make(map[int64]int64)
	num := 0
	for _, rulesetup := range rulesetups {
		params = models.NewQueryParams()
		params.SetInt64("Planid", planid)
		params.SetInt64("rulesetupid", rulesetup.ID)
		rulecanhuis, _ := m.ListRulecanhuis(params)
		rulecanhui := rulecanhuis[0]
		fmt.Printf("%v\n", rulecanhui)

		rulezoneJ, _ := simplejson.NewJson([]byte(rulesetup.Rulezone))
		rulezoneInterface := rulezoneJ.MustArray()

		seatNum := len(rulezoneInterface)
		seatList := make([]*Seat, 0, seatNum)

		indexMap := [1000][1000]*Seat{}

		columns, rows := 0, 0
		count := 0
		beginRow, _ := strconv.Atoi(strings.Split(rulezoneInterface[0].(map[string]interface{})["seatid"].(string), "-")[0])

		for i, rulezoneI := range rulezoneInterface {
			rz := rulezoneI.(map[string]interface{})
			seatID := rz["seatid"].(string)
			seatString := strings.Split(seatID, "-")
			seat := new(Seat)
			seat.Row, _ = strconv.Atoi(seatString[0])
			seat.Column, _ = strconv.Atoi(seatString[1])
			seat.IsUsed = false

			if beginRow != seat.Row {
				beginRow = seat.Row
				rows = rows + 1
				columns = i - count
				count = i
			}

			indexMap[rows][i-count] = seat
		}
		rows = rows + 1

		score := 0
		if rulesetup.Ruleid == "R1" {
			for i := 0; i < rows; i++ {
				lmid, rmid := columns/2, columns/2
				cur := lmid
				for j := 0; j < columns; j++ {
					if j == 0 {
						cur = rmid
					} else if j%2 == 0 {
						rmid = rmid + 1
						cur = rmid
					} else {
						lmid = lmid - 1
						cur = lmid
					}
					seat := indexMap[i][cur]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R2" {
			for i := 0; i < rows; i++ {
				lmid, rmid := columns/2-1, columns/2-1
				cur := lmid
				for j := 0; j < columns; j++ {
					if j == 0 {
						cur = lmid
					} else if j%2 == 0 {
						lmid = lmid - 1
						cur = lmid
					} else {
						rmid = rmid + 1
						cur = rmid
					}
					seat := indexMap[i][cur]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R3" {
			for i := 0; i < rows; i++ {
				for j := 0; j < columns; j++ {
					seat := indexMap[i][j]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R4" {
			for i := 0; i < rows; i++ {
				for j := columns - 1; j >= 0; j-- {
					seat := indexMap[i][j]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R5" {
			for j := 0; j < columns; j++ {
				for i := 0; i < rows; i++ {
					seat := indexMap[i][j]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R6" {
			for j := columns - 1; j < columns; j-- {
				for i := 0; i < rows; i++ {
					seat := indexMap[i][j]
					seat.Score = score
					seatList = append(seatList, seat)
					score = score + 1
				}
			}
		} else if rulesetup.Ruleid == "R7" {
			for i := 0; i < rows; i++ {
				if i%2 == 0 {
					for j := 0; j < columns; j++ {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				} else {
					for j := columns - 1; j >= 0; j-- {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				}
			}
		} else if rulesetup.Ruleid == "R8" {
			for i := rows; i < rows; i++ {
				if i%2 == 0 {
					for j := columns - 1; j >= 0; j-- {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				} else {
					for j := 0; j < columns; j++ {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				}
			}
		} else if rulesetup.Ruleid == "R9" {
			for j := 0; j < columns; j++ {
				if j%2 == 0 {
					for i := 0; i < rows; i++ {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				} else {
					for i := rows - 1; i >= 0; i-- {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				}
			}
		} else if rulesetup.Ruleid == "R10" {
			for j := 0; j < columns; j++ {
				if j%2 == 0 {
					for i := rows - 1; i >= 0; i-- {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				} else {
					for i := 0; i < rows; i++ {
						seat := indexMap[i][j]
						seat.Score = score
						seatList = append(seatList, seat)
						score = score + 1
					}
				}
			}
		}

		sort.Slice(seatList, func(i, j int) bool {
			return seatList[i].Score < seatList[j].Score
		})

		ruleSeatMap[rulesetup.ID] = seatList
		for _, seat := range seatList {
			canhui := new(Canhui)
			canhui.Seatid = fmt.Sprintf("%d-%d", seat.Row, seat.Column)
			canhui.Specialid = canhuiList[num].Specialid
			canhui.ID = canhuiList[num].ID
			canhui.Attender = canhuiList[num].Name
			canhui.Groupid = canhuiList[num].Groupid
			canhui.Bgcolor = rulesetup.Bgcolor
			canhuis = append(canhuis, canhui)
			num = num + 1
		}
	}
	c.OK(&MeetingShowResp{Canhuis: canhuis})
}
