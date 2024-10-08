package datamodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

var allDaysOfWeek = []time.Weekday{
	time.Sunday,
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
}

// Allow any time of day by default (0:00 to 23:59)
const maxWindowLocalTimeEnd = 23*60 + 59
const defaultDeliveryWindowLocalTimeStart = 0
const defaultDeliveryWindowLocalTimeEnd = maxWindowLocalTimeEnd

var validInterruptionLevels = []string{"active", "critical", "passive", "timeSensitive"}

const NotificationMaxIdealWaitTimeForever = -1

type Notification struct {
	ID         string
	Title      string
	Body       string
	BadgeCount int
	ActionName string
	Sound      string

	LaunchImageName string

	RelevanceScore    *float64
	InterruptionLevel string

	ScheduleCondition *Condition

	DeliveryTime DeliveryTime

	// Delivery Window
	DeliveryDaysOfWeek            []time.Weekday
	DeliveryWindowTODStartMinutes int
	DeliveryWindowTODEndMinutes   int

	IdealDeliveryConditions *IdealDeliveryConditions
	CancelationEvents       *[]string
}

type IdealDeliveryConditions struct {
	Condition          Condition `json:"condition"`
	MaxWaitTimeSeconds int       `json:"maxWaitTimeSeconds"`
}

func (i *IdealDeliveryConditions) MaxWaitTime() time.Duration {
	if i.WaitForever() {
		// 200 years in case the caller skips the check WaitForever()
		return time.Hour * 24 * 365 * 200
	}
	return time.Second * time.Duration(i.MaxWaitTimeSeconds)
}

func (i *IdealDeliveryConditions) WaitForever() bool {
	return i.MaxWaitTimeSeconds == NotificationMaxIdealWaitTimeForever
}

type EventInstanceTypeEnum int

const (
	// unknown, should not process. Could be type from future SDK
	EventInstanceTypeUnknown EventInstanceTypeEnum = iota
	// The notification's time is relative to the latest time the event occurred, but once it occurs it won't fire again
	EventInstanceTypeLatestOnce
	// The notification's time is relative to the latest time the event occurred
	EventInstanceTypeLatest
	// The notification's time is relative to the first time the event occurred
	EventInstanceTypeFirst
)

type DeliveryTime struct {
	// Exact point in time
	TimestampEpoch *int64 `json:"timestamp,omitempty"`

	// Event based
	EventName           *string `json:"eventName,omitempty"`
	EventOffsetSeconds  *int    `json:"eventOffsetSeconds,omitempty"`
	EventInstanceString *string `json:"eventInstance,omitempty"`
}

func (dt *DeliveryTime) EventInstance() EventInstanceTypeEnum {
	// Default to latest-once for nil/empty, but not unrecognized
	if dt.EventInstanceString == nil || *dt.EventInstanceString == "" || *dt.EventInstanceString == "latest-once" {
		return EventInstanceTypeLatestOnce
	} else if *dt.EventInstanceString == "latest" {
		return EventInstanceTypeLatest
	}
	if *dt.EventInstanceString == "first" {
		return EventInstanceTypeFirst
	}
	return EventInstanceTypeUnknown
}

func (dt *DeliveryTime) EventOffsetDuration() time.Duration {
	if dt.EventOffsetSeconds == nil {
		return 0
	}
	return time.Duration(*dt.EventOffsetSeconds) * time.Second
}

type jsonNotification struct {
	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	BadgeCount *int   `json:"badgeCount,omitempty"`
	ActionName string `json:"tapActionName,omitempty"`
	Sound      string `json:"sound,omitempty"`

	LaunchImageName string `json:"launchImageName,omitempty"`

	RelevanceScore    *float64 `json:"relevanceScore,omitempty"`
	InterruptionLevel string   `json:"interruptionLevel,omitempty"`

	ScheduleCondition *Condition `json:"scheduleCondition"`

	DeliveryTime             DeliveryTime `json:"deliveryTime,omitempty"`
	DeliveryDaysOfWeekString string       `json:"deliveryDaysOfWeek,omitempty"`
	// HH:MM format
	DeliveryWindowTODStart string `json:"deliveryTimeOfDayStart,omitempty"`
	DeliveryWindowTODEnd   string `json:"deliveryTimeOfDayEnd,omitempty"`

	IdealDeliveryConditions *IdealDeliveryConditions `json:"idealDeliveryConditions,omitempty"`
	CancelationEvents       *[]string                `json:"cancelationEvents,omitempty"`
}

func (a *Notification) Valid() bool {
	return a.Check() == nil
}

func (n *Notification) Check() UserPresentableErrorInterface {
	return n.CheckIgnoreID(false)
}

func (n *Notification) CheckIgnoreID(ignoreID bool) UserPresentableErrorInterface {
	// ID is set later from primary config map ID
	if !ignoreID && n.ID == "" {
		return NewUserPresentableError("Notification must have ID")
	}
	if n.Title == "" &&
		n.Body == "" &&
		n.BadgeCount < 0 {
		return NewUserPresentableError("Notifications must have one or more of: title, body, or badgeCount.")
	}
	if n.DeliveryWindowTODStartMinutes < 0 || n.DeliveryWindowTODStartMinutes > maxWindowLocalTimeEnd {
		return NewUserPresentableError("Notifications must have a deliveryTimeOfDayStart between 0 and 23:59 mins.")
	}
	if n.DeliveryWindowTODEndMinutes < 0 || n.DeliveryWindowTODEndMinutes > maxWindowLocalTimeEnd {
		return NewUserPresentableError("Notifications must have a deliveryTimeOfDayEnd between 0 and 23:59 mins.")
	}
	if n.DeliveryWindowTODStartMinutes > n.DeliveryWindowTODEndMinutes {
		return NewUserPresentableError("Notifications must have a deliveryTimeOfDayStart before deliveryTimeOfDayEnd.")
	}
	if len(n.DeliveryDaysOfWeek) == 0 {
		return NewUserPresentableError("Notifications must have at least one day of week valid for delivery in deliveryDaysOfWeek.")
	}
	if n.RelevanceScore != nil && (*n.RelevanceScore < 0 || *n.RelevanceScore > 1) {
		return NewUserPresentableError("Notification relevanceScore must be between 0 and 1 if provided.")
	}
	if StrictDatamodelParsing && n.InterruptionLevel != "" {
		if !slices.Contains(validInterruptionLevels, n.InterruptionLevel) {
			return NewUserPresentableError(fmt.Sprintf("Notification interruptionLevel must be one of %v, got %v", validInterruptionLevels, n.InterruptionLevel))
		}
	}
	if n.ScheduleCondition != nil {
		if err := n.ScheduleCondition.Validate(); err != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Notification has invalid scheduleCondition [%v]", n.ScheduleCondition.conditionString), err)
		}
	}
	if n.CancelationEvents != nil {
		for _, event := range *n.CancelationEvents {
			if event == "" {
				return NewUserPresentableError("Notification has an empty string in cancelationEvents. Remove it, or specify a non-empty string.")
			}
		}
	}
	if n.IdealDeliveryConditions != nil {
		if conErr := n.IdealDeliveryConditions.Condition.Validate(); conErr != nil {
			return NewUserPresentableErrorWSource(fmt.Sprintf("Notification idealDeliveryCondition invalid for notification with id '%v'", n.ID), conErr)
		}
		if n.IdealDeliveryConditions.MaxWaitTimeSeconds != NotificationMaxIdealWaitTimeForever &&
			n.IdealDeliveryConditions.MaxWaitTimeSeconds < 1 {
			return NewUserPresentableError("Notifications must have a maxWaitTimeSeconds for ideal delivery condition. Valid values are -1 (forever) or values greater than 0.")
		}
	}
	if dtErr := n.DeliveryTime.Check(); dtErr != nil {
		return NewUserPresentableErrorWSource("Notification has invalid delivery time", dtErr)
	}
	return nil
}

func (d *DeliveryTime) Check() UserPresentableErrorInterface {
	if d.TimestampEpoch == nil && d.EventName == nil {
		return NewUserPresentableError("Notification DeliveryTime must have either a Timestamp or an EventName defined.")
	}
	if d.TimestampEpoch != nil && d.EventName != nil {
		return NewUserPresentableError("Notification DeliveryTime cannot have both a Timestamp and an EventName defined.")
	}
	if d.TimestampEpoch != nil && d.EventOffsetSeconds != nil {
		return NewUserPresentableError("Notification DeliveryTime cannot have both a Timestamp and an EventOffset defined.")
	}
	if StrictDatamodelParsing {
		if d.EventInstance() == EventInstanceTypeUnknown {
			return NewUserPresentableError(fmt.Sprintf("Notification eventInstance must be 'first', 'latest', or 'latest-once' (default), got '%v'", d.EventInstanceString))
		}
	}
	return nil
}

func (d *DeliveryTime) Timestamp() *time.Time {
	if d.TimestampEpoch == nil {
		return nil
	}
	time := time.Unix(*d.TimestampEpoch, 0)
	return &time
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	var jn jsonNotification
	err := json.Unmarshal(data, &jn)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an notification with type=notification. Check the format, variable names, and types (eg float vs int).", err)
	}

	n.Title = jn.Title
	n.Body = jn.Body
	n.Sound = jn.Sound
	n.ActionName = jn.ActionName
	n.IdealDeliveryConditions = jn.IdealDeliveryConditions
	n.CancelationEvents = jn.CancelationEvents
	n.DeliveryTime = jn.DeliveryTime
	n.RelevanceScore = jn.RelevanceScore
	n.InterruptionLevel = jn.InterruptionLevel
	n.LaunchImageName = jn.LaunchImageName
	n.ScheduleCondition = jn.ScheduleCondition

	if jn.BadgeCount != nil {
		n.BadgeCount = *jn.BadgeCount
		if StrictDatamodelParsing && n.BadgeCount < 0 {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError("Notification badgeCount must be greater than or equal to 0"))
		}
	} else {
		n.BadgeCount = -1 // default to -1 for unset
	}
	if jn.DeliveryDaysOfWeekString != "" {
		n.DeliveryDaysOfWeek = parseDaysOfWeekString(jn.DeliveryDaysOfWeekString)
	} else {
		// Default to all days of week
		n.DeliveryDaysOfWeek = allDaysOfWeek
	}

	// Defaults could change over time, so either all custom, or all default or config could be invalid
	if (jn.DeliveryWindowTODStart == "" && jn.DeliveryWindowTODEnd != "") ||
		(jn.DeliveryWindowTODStart != "" && jn.DeliveryWindowTODEnd == "") {
		return NewUserErrorForJsonIssue(data, NewUserPresentableError("DeliveryTime must have both deliveryTimeOfDayStart and deliveryTimeOfDayEnd defined if either is defined."))
	}
	if jn.DeliveryWindowTODStart != "" {
		deliveryStart, err := parseMinutesFromHHMMString(jn.DeliveryWindowTODStart)
		if err != nil && StrictDatamodelParsing {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError("Invalid deliveryTimeOfDayStart. Expect HH:MM format. Was: "+jn.DeliveryWindowTODStart))
		} else if err != nil {
			fmt.Printf("CriticalMoments: invalid deliveryTimeOfDayStart [%v]. Using default: %v\n", jn.DeliveryWindowTODStart, defaultDeliveryWindowLocalTimeStart)
			n.DeliveryWindowTODStartMinutes = defaultDeliveryWindowLocalTimeStart
		} else {
			n.DeliveryWindowTODStartMinutes = deliveryStart
		}
	} else {
		n.DeliveryWindowTODStartMinutes = defaultDeliveryWindowLocalTimeStart
	}
	if jn.DeliveryWindowTODEnd != "" {
		deliveryEnd, err := parseMinutesFromHHMMString(jn.DeliveryWindowTODEnd)
		if err != nil && StrictDatamodelParsing {
			return NewUserErrorForJsonIssue(data, NewUserPresentableError("Invalid notification 'deliveryTimeOfDayEnd'. Expect HH:MM format. Was: "+jn.DeliveryWindowTODEnd))
		} else if err != nil {
			fmt.Printf("CriticalMoments: invalid notification 'deliveryTimeOfDayEnd' [%v]. Using default: %v\n", jn.DeliveryWindowTODEnd, defaultDeliveryWindowLocalTimeEnd)
			n.DeliveryWindowTODEndMinutes = defaultDeliveryWindowLocalTimeEnd
		} else {
			n.DeliveryWindowTODEndMinutes = deliveryEnd
		}
	} else {
		n.DeliveryWindowTODEndMinutes = defaultDeliveryWindowLocalTimeEnd
	}

	// ignore ID which is set later from primary config
	if err := n.CheckIgnoreID(true); err != nil {
		return NewUserErrorForJsonIssue(data, err)
	}

	return nil
}

func parseMinutesFromHHMMString(i string) (int, error) {
	if i == "" {
		return 0, errors.New("empty string for HH:MM time of day")
	}
	values := strings.Split(i, ":")
	if len(values) != 2 {
		return 0, errors.New("invalid HH:MM time of day")
	}
	hours, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(values[1])
	if err != nil {
		return 0, err
	}
	if hours < 0 || hours > 23 || minutes < 0 || minutes > 59 {
		return 0, errors.New("invalid HH:MM time of day")
	}
	return hours*60 + minutes, nil
}

// Parsed comman separated day of week strings, removing dupes and standardizing order
func parseDaysOfWeekString(i string) []time.Weekday {
	days := []time.Weekday{}
	components := strings.Split(i, ",")

	// This format orders and dedupes them, for consistency
	for dow := range allDaysOfWeek {
		dowWD := time.Weekday(dow)
		dowName := dowWD.String()
		dayFound := false
		for _, dayStringComponent := range components {
			if dowName == dayStringComponent {
				dayFound = true
			}
		}
		if dayFound {
			days = append(days, dowWD)
		}
	}

	return days
}

const NotificationUniqueIDPrefix = "io.criticalmoments.notifications."

func (n *Notification) UniqueID() string {
	return fmt.Sprintf("%v%v", NotificationUniqueIDPrefix, n.ID)
}

func (n *Notification) DeliveredEventName() string {
	return fmt.Sprintf("notifications:delivered:%v", n.UniqueID())
}

// Gomobile accessors. Gomobiledoesn't support pointers
func (n *Notification) GetRelevanceScore() float64 {
	return *n.RelevanceScore
}
func (n *Notification) HasRelevanceScore() bool {
	return n.RelevanceScore != nil
}
