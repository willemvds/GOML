package api

import (
	"cmp"
	"slices"
	"time"
)

type Action uint8

const ActionOff = 0
const ActionOn = 1

const cmpGreaterThan = +1
const cmpEqual = 0
const cmpLessThan = -1

type TimeOfDay struct {
	hour   uint
	minute uint
}

func NewTimeOfDay(hour uint, minute uint) TimeOfDay {
	if hour > 23 {
		hour = 23
	}
	if minute > 59 {
		minute = 59
	}

	return TimeOfDay{
		hour:   hour,
		minute: minute,
	}
}

func compareTimeOfDay(a TimeOfDay, b TimeOfDay) int {
	hourOrd := cmp.Compare(a.hour, b.hour)
	if hourOrd != 0 {
		return hourOrd
	}

	return cmp.Compare(a.minute, b.minute)
}

func compareScheduleItem(a ScheduleItem, b ScheduleItem) int {
	return compareTimeOfDay(a.when, b.when)
}

type ScheduleItem struct {
	when   TimeOfDay
	action Action
}

type Schedule struct {
	events []ScheduleItem
}

func (schedule *Schedule) addScheduleItem(when TimeOfDay, action Action) {
	schedule.events = append(schedule.events, newScheduleItem(when, action))
	slices.SortFunc(schedule.events, compareScheduleItem)
}

func (schedule *Schedule) AddOnEvent(when TimeOfDay) {
	schedule.addScheduleItem(when, ActionOn)
}

func (schedule *Schedule) AddOffEvent(when TimeOfDay) {
	schedule.addScheduleItem(when, ActionOff)
}

func (schedule *Schedule) LastActionAt(when time.Time) Action {
	lastAction := schedule.events[len(schedule.events)-1].action

	whenTimeOfDay := NewTimeOfDay(uint(when.Hour()), uint(when.Minute()))
	for _, ev := range schedule.events {
		if compareTimeOfDay(ev.when, whenTimeOfDay) == cmpGreaterThan {
			break
		}
		lastAction = ev.action
	}

	return lastAction
}

func newScheduleItem(when TimeOfDay, action Action) ScheduleItem {
	if action > ActionOn {
		action = ActionOn
	}

	return ScheduleItem{
		when,
		action,
	}
}
