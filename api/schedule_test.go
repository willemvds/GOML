package api

import (
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	schedule := Schedule{}
	schedule.AddOnEvent(NewTimeOfDay(10, 10))

	if len(schedule.events) != 1 {
		t.Errorf("Whut?")
	}
}

func TestScheduleOrdering(t *testing.T) {
	schedule := Schedule{}
	schedule.AddOnEvent(NewTimeOfDay(10, 10))
	schedule.AddOnEvent(NewTimeOfDay(13, 13))
	schedule.AddOnEvent(NewTimeOfDay(15, 15))
	schedule.AddOnEvent(NewTimeOfDay(5, 5))
	schedule.AddOnEvent(NewTimeOfDay(8, 8))

	prevItem := ScheduleItem{when: NewTimeOfDay(0, 0), action: ActionOn}
	for _, item := range schedule.events {
		// We can work with only the hour since we control the input values
		// and don't want to use the same comparison function as used inside the sort.
		if item.when.hour < prevItem.when.hour {
			t.Errorf("Found ScheduleItem that is out of order.")
		}
		prevItem = item
	}
}

func TestConstraints(t *testing.T) {
	timeOfDay := NewTimeOfDay(444, 222)

	if timeOfDay.hour != 23 {
		t.Errorf("Failed to keep 'hour' constrained, sitting with %d\n", timeOfDay.hour)
	}
	if timeOfDay.minute != 59 {
		t.Errorf("Failed to keep 'minute' constrained, sitting with %d\n", timeOfDay.minute)
	}
}

func TestLastActionAt(t *testing.T) {
	schedule := Schedule{}
	schedule.AddOnEvent(NewTimeOfDay(13, 13))

	// 3pm
	when, err := time.Parse(time.RFC3339, "1984-10-13T15:00:00+02:00")
	if err != nil {
		t.Errorf("Failed to construct time value for test: %s\n", err)
	}

	action := schedule.LastActionAt(when)

	if action != ActionOn {
		t.Errorf("Last Action should be ON but found %v\n", action)
	}
}

func TestLastActionAtOff(t *testing.T) {
	schedule := Schedule{}
	schedule.AddOnEvent(NewTimeOfDay(22, 00))
	schedule.AddOffEvent(NewTimeOfDay(5, 30))

	when, err := time.Parse(time.RFC3339, "1984-10-13T05:31:00+02:00")
	if err != nil {
		t.Errorf("Failed to construct time value for test: %s\n", err)
	}

	action := schedule.LastActionAt(when)

	if action != ActionOff {
		t.Errorf("Last Action should be OFF but found %v\n", action)
	}
}

func TestLastActionAtPreviousDay(t *testing.T) {
	schedule := Schedule{}
	schedule.AddOnEvent(NewTimeOfDay(22, 00))
	schedule.AddOffEvent(NewTimeOfDay(5, 30))

	when, err := time.Parse(time.RFC3339, "1984-10-13T02:00:00+02:00")
	if err != nil {
		t.Errorf("Failed to construct time value for test: %s\n", err)
	}

	action := schedule.LastActionAt(when)

	if action != ActionOn {
		t.Errorf("Last Action should be ON but found %v\n", action)
	}
}
