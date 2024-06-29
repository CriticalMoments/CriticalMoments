package datamodel

import (
	"encoding/json"
	"fmt"
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

type Notification struct {
	ID         string
	Title      string
	Body       string
	ActionName string

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

type DeliveryTime struct {
	TimestampEpoch *int64  `json:"timestamp,omitempty"`
	EventName      *string `json:"eventName,omitempty"`
	EventOffset    *int    `json:"eventOffset,omitempty"`
}

type jsonNotification struct {
	Title      string `json:"title,omitempty"`
	Body       string `json:"body,omitempty"`
	ActionName string `json:"actionName,omitempty"`

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
	if n.Title == "" {
		return "Notifications must have a title."
	}
	if len(n.DeliveryDaysOfWeek) == 0 {
		return "Notifications must have at least one day of week valid for delivery."
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
	return ""
}

func (n *Notification) UnmarshalJSON(data []byte) error {
	var jn jsonNotification
	err := json.Unmarshal(data, &jn)
	if err != nil {
		return NewUserPresentableErrorWSource("Unable to parse the json of an notification with type=notification. Check the format, variable names, and types (eg float vs int).", err)
	}

	n.Title = jn.Title
	n.Body = jn.Body
	n.ActionName = jn.ActionName
	n.IdealDevlieryConditions = jn.IdealDeliveryConditions
	n.CancelationEvents = jn.CancelationEvents
	n.DeliveryTime = jn.DeliveryTime

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
