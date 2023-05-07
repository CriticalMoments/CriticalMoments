package datamodel

type Trigger struct {
	Event  Event
	Action ActionContainer
}

type jsonTrigger struct {
	EventName  string `json:"eventName"`
	ActionName string `json:"actionName"`
}

// TODO - not using validate yet

func (t jsonTrigger) Validate() bool {
	return t.ValidateReturningUserReadableIssue() == ""
}

func (t jsonTrigger) ValidateReturningUserReadableIssue() string {
	if t.EventName == "" {
		return "All triggers require an event"
	}
	if t.ActionName == "" {
		return "All triggers require an action name"
	}
	return ""
}
