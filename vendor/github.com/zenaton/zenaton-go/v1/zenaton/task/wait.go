package task

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"reflect"

	"github.com/zenaton/zenaton-go/v1/zenaton/internal/engine"
	"github.com/zenaton/zenaton-go/v1/zenaton/service/serializer"
)

const (
	modeAt        = "AT"
	modeWeekDay   = "WEEK_DAY"
	modeMonthDay  = "MONTH_DAY"
	modeTimestamp = "TIMESTAMP"
)

// A WaitTask is a special form of task. This holds all the relevant information to wait for a duration or until a timestamp.
type WaitTask struct {
	task      *Instance
	eventName string
	buffer    []duration
	mode      string
	timezone  *time.Location
}

// Wait returns a *WaitTask. This is a special type of task that can be used to wait for a given duration, or until a
// given time. For example: task.Wait().Weeks(2).Hours(2).Minutes(15).Seconds(23).Execute() will wait for exactly 14
// days, 2 hours, 15 minutes and 23 seconds!
func Wait() *WaitTask {
	return &WaitTask{
		task: waitTask,
	}
}

var waitTask = NewCustom("_Wait", &wait{}).New()

type wait struct{}

func (wt *wait) Handle() (interface{}, error) { return nil, nil }

// ForEvent lets you wait for an event, with an optional timeout.
// If you need to wait until an event occurs, it's as easy as:
//
// 		var event Event
// 		task.Wait().ForEvent("UserActivatedEvent").Execute().Output(&event)
//
// In this example, the workflow will stall up to receiving a UserActivatedEvent. Then event variable will contain the
// received event object.
//
// Usually, you want to add a timeout, in order to ensure that the workflow will finish:
//
// 		var event Event
// 		task.Wait().ForEvent("UserActivatedEvent").Days(5).Execute().Output(&event)
//
// In this example, after 5 days, the execution of the workflow will be released, but the event will be null. You can
// use the same methods as above to define the duration or the date of this timeout.
//
// If a workflow instance receives an event before the instructions to wait for it, then this event will only be handled
// by the OnEvent method. If you want to have a persistent event, you can do something like:
//
// 		func (w WorkflowType) Handle() (interface{}, error) {
// 			...
// 			if w.Event != nil {
// 				task.Wait().ForEvent("UserActivatedEvent").Days(5).Execute(&w.Event)
// 			}
// 			...
// 		}
//
// 		func (w WorkflowType) OnEvent(name string, data interface{}) {
// 			if name == "UserActivatedEvent" {
// 				w.Event = data
// 			}
// 		}
func (w *WaitTask) ForEvent(eventName string) *WaitTask {
	w.eventName = eventName
	return w
}

// Event simply returns the event name that the WaitTask is waiting for
func (w *WaitTask) Event() string {
	return w.eventName
}

// ************************************************************************
// ************************************************************************
// DURATION METHODS
// ************************************************************************
// ************************************************************************
type duration struct {
	method string
	value  interface{}
}

// Seconds waits for the provided number of seconds
func (w *WaitTask) Seconds(value int64) *WaitTask {
	w.push(duration{
		method: "seconds",
		value:  value,
	})

	return w
}

// Minutes waits for the provided number of minutes
func (w *WaitTask) Minutes(value int64) *WaitTask {
	w.push(duration{
		method: "minutes",
		value:  value,
	})

	return w
}

// Hours waits for the provided number of hours
func (w *WaitTask) Hours(value int64) *WaitTask {
	w.push(duration{
		method: "hours",
		value:  value,
	})

	return w
}

// Days waits for the provided number of days
func (w *WaitTask) Days(value int64) *WaitTask {
	w.push(duration{
		method: "days",
		value:  value,
	})

	return w
}

// Weeks waits for the provided number of weeks
func (w *WaitTask) Weeks(value int64) *WaitTask {
	w.push(duration{
		method: "weeks",
		value:  value,
	})

	return w
}

// Months waits for the provided number of weeks
func (w *WaitTask) Months(value int64) *WaitTask {
	w.push(duration{
		method: "months",
		value:  value,
	})

	return w
}

// Years waits for the provided number of years
func (w *WaitTask) Years(value int64) *WaitTask {
	w.push(duration{
		method: "years",
		value:  value,
	})

	return w
}

// ************************************************************************
// ************************************************************************
// TIMESTAMP METHODS
// ************************************************************************
// ************************************************************************

// Timezone allows you to define a different timezone than your local one. Use this before launching a wait task.
// For example:
// 		task.Wait().Timezone("Europe/Paris").Execute()
func (w *WaitTask) Timezone(timezone string) error {
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	w.timezone = tz
	return nil
}

// Timestamp waits until the given timestamp
func (w *WaitTask) Timestamp(value int64) *WaitTask {
	w.push(duration{"timestamp", value})
	return w
}

// At waits until the time given as a string.
// For example: task.Wait().At("15:10:23").Execute() waits until 3:10PM and 23 seconds.
func (w *WaitTask) At(value string) *WaitTask {
	w.push(duration{"at", value})
	return w
}

// DayOfMonth waits until the next given day of the month.
// For example: .DayOfMonth(12) waits to next 12th day (same hour)
func (w *WaitTask) DayOfMonth(value int) *WaitTask {
	w.push(duration{"dayOfMonth", value})
	return w
}

// Monday waits until the next monday (same hour)
func (w *WaitTask) Monday(value int) *WaitTask {
	w.push(duration{"Monday", value})
	return w
}

// Tuesday waits until the next tuesday (same hour)
func (w *WaitTask) Tuesday(value int) *WaitTask {
	w.push(duration{"Tuesday", value})
	return w
}

// Wednesday waits until the next wednesday (same hour)
func (w *WaitTask) Wednesday(value int) *WaitTask {
	w.push(duration{"Wednesday", value})
	return w
}

// Thursday waits until the next thursday (same hour)
func (w *WaitTask) Thursday(value int) *WaitTask {
	w.push(duration{"Thursday", value})
	return w
}

// Friday waits until the next friday (same hour)
func (w *WaitTask) Friday(value int) *WaitTask {
	w.push(duration{"Friday", value})
	return w
}

// Saturday waits until the next saturday (same hour)
func (w *WaitTask) Saturday(value int) *WaitTask {
	w.push(duration{"Saturday", value})
	return w
}

// Sunday waits until the next sunday (same hour)
func (w *WaitTask) Sunday(value int) *WaitTask {
	w.push(duration{"Sunday", value})
	return w
}

func (w *WaitTask) push(data duration) {
	w.buffer = append(w.buffer, data)
}

func (w *WaitTask) initNowThen() (time.Time, time.Time) {
	// get set or current time zone

	if w.timezone == nil {
		w.timezone = time.Local
	}
	n := time.Now()
	var now = time.Date(n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second(), n.Nanosecond(), w.timezone)
	var then = now
	return now, then
}

// GetTimestampOrDuration will return either a timestamp or a duration (but not both). If it returns a duration,
// the WaitTask will wait for that length of time (in seconds). If it returns a timestamp, it will wait until that timestamp.
func (w *WaitTask) GetTimestampOrDuration() (int64, int64, error) {

	now, then := w.initNowThen()

	w.mode = ""

	var err error
	for _, duration := range w.buffer {
		then, err = w.apply(duration.method, duration.value, now, then)
		if err != nil {
			return 0, 0, err
		}
	}

	isTimestamp := w.mode != ""

	if isTimestamp {
		//todo: these shouldn't be 0, right? what if the time until then is actually 0?
		return then.Unix(), 0, nil
	}

	return 0, then.Unix() - now.Unix(), nil
}

func (w *WaitTask) apply(method string, value interface{}, now, then time.Time) (time.Time, error) {
	switch method {
	case "timestamp":
		return w._timestamp(value.(int64))
	case "at":
		return w._at(value.(string), now, then)
	case "dayOfMonth":
		return w._dayOfMonth(value.(int), now, then)
	case "monday":
		return w._weekDay(value.(int), 1, then)
	case "tuesday":
		return w._weekDay(value.(int), 2, then)
	case "wednesday":
		return w._weekDay(value.(int), 3, then)
	case "thursday":
		return w._weekDay(value.(int), 4, then)
	case "friday":
		return w._weekDay(value.(int), 5, then)
	case "saturday":
		return w._weekDay(value.(int), 6, then)
	case "sunday":
		return w._weekDay(value.(int), 7, then)
	default:
		return w._applyDuration(method, value.(int64), then)
	}
}

func (w *WaitTask) _timestamp(timestamp int64) (time.Time, error) {
	err := w._setMode(modeTimestamp)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timestamp, 0), nil
}

func (w *WaitTask) _at(t string, now, then time.Time) (time.Time, error) {
	err := w._setMode(modeAt)
	if err != nil {
		return time.Time{}, err
	}

	segments := strings.Split(t, ":")
	h, err := strconv.Atoi(segments[0])
	if err != nil {
		return time.Time{}, errors.New("time formatted incorrectly")
	}
	var m int
	if len(segments) > 1 {
		m, err = strconv.Atoi(segments[1])
		if err != nil {
			return time.Time{}, errors.New("time formatted incorrectly")
		}
	}
	var s int
	if len(segments) > 2 {
		s, err = strconv.Atoi(segments[2])
		if err != nil {
			return time.Time{}, errors.New("time formatted incorrectly")
		}
	}

	then = time.Date(now.Year(), now.Month(), now.Day(), h, m, s, 0, w.timezone)

	if now.After(then) {
		switch w.mode {
		case modeAt:
			then = then.AddDate(0, 0, 1)
			break
		case modeWeekDay:
			then = then.AddDate(0, 0, 7)
			break
		case modeMonthDay:
			then = then.AddDate(0, 1, 0)
			break
		default:
			return time.Time{}, errors.New("Unknown mode: " + w.mode)
		}
	}

	return then, nil
}

func (w *WaitTask) _dayOfMonth(day int, now, then time.Time) (time.Time, error) {

	err := w._setMode(modeMonthDay)
	if err != nil {
		return time.Time{}, err
	}

	then = time.Date(now.Year(), now.Month(), day, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), w.timezone)

	if now.After(then) {
		then = then.AddDate(0, 1, 0)
	}

	return then, nil
}

func (w *WaitTask) _weekDay(n int, day int, then time.Time) (time.Time, error) {
	err := w._setMode(modeWeekDay)
	if err != nil {
		return time.Time{}, err
	}

	d := int(then.Weekday())
	then = then.AddDate(0, 0, day-d)

	if d > day {
		then.AddDate(0, 0, n*7)
	} else {
		then.AddDate(0, 0, (n-1)*7)
	}

	return then, nil
}

func (w *WaitTask) _setMode(mode string) error {
	// can not apply twice the same method
	if mode == w.mode {
		return errors.New("incompatible definition in WaitTask methods")
	}
	// timestamp can only be used alone
	if (w.mode != "" && mode == modeTimestamp) || w.mode == modeTimestamp {
		return errors.New("incompatible definition in WaitTask methods")
	}

	// other mode takes precedence to modeAt
	if w.mode == "" || modeAt == w.mode {
		w.mode = mode
	}

	return nil
}

func (w *WaitTask) _applyDuration(method string, value int64, then time.Time) (time.Time, error) {
	switch method {
	case "seconds":
		return then.Add(time.Duration(value) * time.Second), nil
	case "minutes":
		return then.Add(time.Duration(value) * time.Minute), nil
	case "hours":
		return then.Add(time.Duration(value) * time.Hour), nil
	case "days":
		return then.AddDate(0, 0, int(value)), nil
	case "weeks":
		return then.AddDate(0, 0, int(value)*7), nil
	case "months":
		return then.AddDate(0, int(value), 0), nil
	case "years":
		return then.AddDate(int(value), 0, 0), nil
	default:
		return time.Time{}, errors.New("Unknown method " + method)
	}
}

// Handle is used by the engine, and shouldn't be called directly.
func (w *WaitTask) Handle() (interface{}, error) {
	return w.task.Handle()
}

// LaunchInfo is used by the engine to determine the relevant information for dispatching Jobs (tasks or workflows)
func (w *WaitTask) LaunchInfo() engine.LaunchInfo {
	return w.task.LaunchInfo()
}

// GetName will return the name of the task. Here it will always be _Wait.
func (w *WaitTask) GetName() string {
	return w.task.GetName()
}

// GetData will return the underlying handler of the task.
func (w *WaitTask) GetData() engine.Handler {
	return waitTask.Handler
}

// WaitExecution is the result of calling Execute on a WaitTask. This will hold the serialized event value (if there is one).
type WaitExecution struct {
	SerializedEventValue string
}

// Execute actually starts the WaitTask. It returns a WaitExecution that can be used to retrieve event data (if the WaitTask
// was waiting for an event)
func (w *WaitTask) Execute() WaitExecution {
	_, serializedEvents, _ := engine.NewEngine().Execute([]engine.Job{w})

	var waitExecution WaitExecution
	if len(serializedEvents) == 0 {
		waitExecution.SerializedEventValue = ""
	} else {
		waitExecution.SerializedEventValue = serializedEvents[0]
	}

	return waitExecution
}

// EventReceived returns true if the Event was received. This is useful in the case that you have a timeout on the
// wait event. For example:
// 		task.Wait().ForEvent("UserActivatedEvent").Days(5).Execute().EventReceived()
// This will wait until an event called "UserActivatedEvent", but will stop waiting after 5 days. EventReceived will
// return true if the event was received and false if the 5 days are up before the event being received.
func (we WaitExecution) EventReceived() bool {
	return we.SerializedEventValue != ""
}

// Output will give you the data passed in the event. You must pass a pointer to Output, and the event data will be
// unmarshaled from json into your pointer. For example:
// 		var event Event
// 		task.Wait().ForEvent("UserActivatedEvent").Execute().Output(&event)
func (we WaitExecution) Output(value interface{}) {

	if we.SerializedEventValue == "" {
		return
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic(fmt.Sprint("must pass a non-nil pointer to WaitExecution.Output"))
	}

	var eventMap map[string]json.RawMessage
	err := serializer.Decode(we.SerializedEventValue, &eventMap)
	if err != nil {
		panic(err)
	}

	eventInput := string(eventMap["event_input"])
	// todo: what if the user wants to send a string "{}". should fix this
	if eventInput == `"{}"` {
		return
	}

	err = serializer.Decode(eventInput, &value)
	if err != nil {
		panic(err)
	}
}
