package datamodel

import (
	"encoding/json"
	"fmt"
	"slices"
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

// 9am to 9pm
const defaultDeliveryWindowLocalTimeStart = 9 * 60 * 60
const defaultDeliveryWindowLocalTimeEnd = 21 * 60 * 60

var validInterruptionLevels = []string{"active", "critical", "passive", "timeSensitive"}

type Notification struct {
	ID         string
	Title      string
	Body       string
	ActionName string
	Sound      string

	RelevanceScore    *float64
	InterruptionLevel string

	DeliveryTime                      DeliveryTime
	DeliveryDaysOfWeek                []time.Weekday
	DeliveryWindowLocalTimeOfDayStart int
	DeliveryWindowLocalTimeOfDayEnd   int

	IdealDevlieryConditions *IdealDevlieryConditions
	CancelationEvents       *[]string
}

type IdealDevlieryConditions struct {
	Condition   Condition `json:"condition"`
	MaxWaitTime int       `json:"maxWaitTime"`
}

type EventInstanceTypeEnum int

const (
	// unknown, should not process. Could be type from future SDK
	EventInstanceTypeUnknown EventInstanceTypeEnum = iota
	// The notification's time is relative to the latest time the event occured
	EventInstanceTypeLatest
	// The notification's time is relative to the first time the event occured
	EventInstanceTypeFirst
)

type DeliveryTime struct {
	TimestampEpoch      *int64  `json:"timestamp,omitempty"`
	EventName           *string `json:"eventName,omitempty"`
	EventOffset         *int    `json:"eventOffset,omitempty"`
	EventInstanceString *string `json:"eventInstance,omitempty"`
}

func (dt *DeliveryTime) EventInstance() EventInstanceTypeEnum {
	// Default to latest for nil/empty, but not unrecognized
	if dt.EventInstanceString == nil || *dt.EventInstanceString == "" || *dt.EventInstanceString == "latest" {
		return EventInstanceTypeLatest
	}
	if dt.EventInstanceString != nil && *dt.EventInstanceString == "first" {
		return EventInstanceTypeFirst
	}
	return EventInstanceTypeUnknown
}

type jsonNotification struct {
	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	ActionName string `json:"actionName,omitempty"`
	Sound      string `json:"sound,omitempty"`

	RelevanceScore    *float64 `json:"relevanceScore,omitempty"`
	InterruptionLevel string   `json:"interruptionLevel,omitempty"`

	DeliveryTime                      DeliveryTime `json:"deliveryTime,omitempty"`
	DeliveryDaysOfWeekString          string       `json:"deliveryDaysOfWeek,omitempty"`
	DeliveryWindowLocalTimeOfDayStart *int         `json:"deliveryLocalTimeOfDayStart,omitempty"`
	DeliveryWindowLocalTimeOfDayEnd   *int         `json:"deliveryLocalTimeOfDayEnd,omitempty"`

	IdealDeliveryConditions *IdealDevlieryConditions `json:"idealDeliveryConditions,omitempty"`
	CancelationEvents       *[]string                `json:"cancelationEvents,omitempty"`
}

func (a *Notification) Validate() bool {
	return a.ValidateReturningUserReadableIssue() == ""
}

func (n *Notification) ValidateReturningUserReadableIssue() string {
	return n.ValidateReturningUserReadableIssueIgnoreID(false)
}

func (n *Notification) ValidateReturningUserReadableIssueIgnoreID(ignoreID bool) string {
	if !ignoreID && n.ID == "" {
		return "Notification must have ID"
	}
	if n.Title == "" && n.Body == "" {
		return "Notifications must have a title and/or a body."
	}
	if len(n.DeliveryDaysOfWeek) == 0 {
		return "Notifications must have at least one day of week valid for delivery."
	}
	if n.RelevanceScore != nil && (*n.RelevanceScore < 0 || *n.RelevanceScore > 1) {
		return "Relevance score must be between 0 and 1 if provided."
	}
	if StrictDatamodelParsing && n.InterruptionLevel != "" {
		if !slices.Contains(validInterruptionLevels, n.InterruptionLevel) {
			return fmt.Sprintf("Interruption level must be one of %v, got %v", validInterruptionLevels, n.InterruptionLevel)
		}
	}
	if n.CancelationEvents != nil {
		for _, event := range *n.CancelationEvents {
			if event == "" {
				return fmt.Sprintf("Notification '%v' has an blank cancelation event", n.ID)
			}
		}
	}
	if n.IdealDevlieryConditions != nil {
		if conErr := n.IdealDevlieryConditions.Condition.Validate(); conErr != nil {
			conErrUserReadable := "Unknown error"
			if uperr, ok := conErr.(*UserPresentableError); ok {
				conErrUserReadable = uperr.UserErrorString()
			}
			return "Notification has invalid ideal delivery condition: " + conErrUserReadable
		}
		if n.IdealDevlieryConditions.MaxWaitTime < -1 || n.IdealDevlieryConditions.MaxWaitTime == 0 {
			return "Notifications must have a max wait time for ideal delivery condition. Valid values are -1 (forever) or values greater than 0."
		}
	}
	if dtErr := n.DeliveryTime.ValidateReturningUserReadableIssue(); dtErr != "" {
		return "Notification has invalid delivery time: " + dtErr
	}
	return ""
}

func (d *DeliveryTime) ValidateReturningUserReadableIssue() string {
	if d.TimestampEpoch == nil && d.EventName == nil {
		return "DeliveryTime must have either a Timestamp or an EventName defined."
	}
	if d.TimestampEpoch != nil && d.EventName != nil {
		return "DeliveryTime cannot have both a Timestamp and an EventName defined."
	}
	if d.TimestampEpoch != nil && d.EventOffset != nil {
		return "DeliveryTime cannot have both a Timestamp and an EventOffset defined."
	}
	if StrictDatamodelParsing {
		if d.EventInstance() == EventInstanceTypeUnknown {
			return fmt.Sprintf("Notification event instance must be 'first' or 'latest', got '%v'", d.EventInstanceString)
		}
	}
	return ""
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
	n.IdealDevlieryConditions = jn.IdealDeliveryConditions
	n.CancelationEvents = jn.CancelationEvents
	n.DeliveryTime = jn.DeliveryTime
	n.RelevanceScore = jn.RelevanceScore
	n.InterruptionLevel = jn.InterruptionLevel

	if jn.DeliveryDaysOfWeekString != "" {
		n.DeliveryDaysOfWeek = parseDaysOfWeekString(jn.DeliveryDaysOfWeekString)
	} else {
		// Default to all days of week
		n.DeliveryDaysOfWeek = allDaysOfWeek
	}
	if jn.DeliveryWindowLocalTimeOfDayStart != nil {
		n.DeliveryWindowLocalTimeOfDayStart = *jn.DeliveryWindowLocalTimeOfDayStart
	} else {
		n.DeliveryWindowLocalTimeOfDayStart = defaultDeliveryWindowLocalTimeStart
	}
	if jn.DeliveryWindowLocalTimeOfDayEnd != nil {
		n.DeliveryWindowLocalTimeOfDayEnd = *jn.DeliveryWindowLocalTimeOfDayEnd
	} else {
		n.DeliveryWindowLocalTimeOfDayEnd = defaultDeliveryWindowLocalTimeEnd
	}

	// ignore ID which is set later from primary config
	if validationIssue := n.ValidateReturningUserReadableIssueIgnoreID(true); validationIssue != "" {
		return NewUserPresentableError(validationIssue)
	}

	return nil
}

// Parsed comman separated day of week strings, removing dupes and standardizing order
func parseDaysOfWeekString(i string) []time.Weekday {
	days := []time.Weekday{}
	components := strings.Split(i, ",")

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

// Gomobile accessor (doesn't support pointers)
func (n *Notification) GetRelevanceScore() float64 {
	return *n.RelevanceScore
}

func (n *Notification) HasRelevanceScore() bool {
	return n.RelevanceScore != nil
}
