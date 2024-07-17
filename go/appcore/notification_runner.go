package appcore

import (
	"errors"
	"slices"
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type NotificationPlan struct {
	unscheduledNotifications []*datamodel.Notification
	scheduledNotifications   []*ScheduledNotification
}

// Gomobile needs count and AtIndex accessors for UnscheduledNotifications and scheduledNotifications
func (plan *NotificationPlan) UnscheduledNotificationCount() int {
	return len(plan.unscheduledNotifications)
}
func (plan *NotificationPlan) UnscheduledNotificationAtIndex(index int) *datamodel.Notification {
	if index >= len(plan.unscheduledNotifications) {
		return nil
	}
	return plan.unscheduledNotifications[index]
}
func (plan *NotificationPlan) ScheduledNotificationCount() int {
	return len(plan.scheduledNotifications)
}
func (plan *NotificationPlan) ScheduledNotificationAtIndex(index int) *ScheduledNotification {
	if index >= len(plan.scheduledNotifications) {
		return nil
	}
	return plan.scheduledNotifications[index]
}

type ScheduledNotification struct {
	Notification *datamodel.Notification
	scheduledAt  time.Time
}

// for gomobile, as time.Time is not supported
func (sn *ScheduledNotification) ScheduledAtEpochMilliseconds() int64 {
	return sn.scheduledAt.UnixMilli()
}

func (ac *Appcore) initializeNotificationPlan() error {
	if ac.notificationPlan == nil {
		return ac.ForceUpdateNotificationPlan()
	}
	return nil
}

func (ac *Appcore) ForceUpdateNotificationPlan() error {
	plan, err := ac.generateNotificationPlan()
	if err != nil {
		return err
	}
	ac.notificationPlan = &plan
	err = ac.libBindings.UpdateNotificationPlan(&plan)
	if err != nil {
		return err
	}

	return nil
}

func (ac *Appcore) generateNotificationPlan() (NotificationPlan, error) {
	plan := NotificationPlan{
		unscheduledNotifications: make([]*datamodel.Notification, 0),
		scheduledNotifications:   make([]*ScheduledNotification, 0),
	}
	for _, notification := range ac.config.Notifications {
		if deliveryTimestamp := ac.notificationDeliveryTime(notification); deliveryTimestamp != nil {
			sn := ScheduledNotification{
				Notification: notification,
				scheduledAt:  *deliveryTimestamp,
			}
			plan.scheduledNotifications = append(plan.scheduledNotifications, &sn)
		} else {
			plan.unscheduledNotifications = append(plan.unscheduledNotifications, notification)
		}
	}

	return plan, nil
}

// Get the delivery time of a notification
// 1) First check when it should be delivered (static time, event based)
// 2) Then consider ideal delivery window, delivering sooner or later if we have special targeting in mind
// 3) Then consider the allowed time of day, and days of week for delivery
func (ac *Appcore) notificationDeliveryTime(notification *datamodel.Notification) *time.Time {
	nonIdealDeliveryTime := ac.baseDeliveryTimeForNotification(notification)
	idealDeliveryTime := ac.shiftDeliveryTimeForIdealWindow(notification, nonIdealDeliveryTime)
	shiftedDeliveryTime := shiftDeliveryTimeForFilters(notification, idealDeliveryTime)
	return shiftedDeliveryTime
}

// timeNow is a mockable version of time.Now for testing
var timeNow = time.Now

// Checks if this notification has an ideal delivery window and now is currently in the time-range of that window
func notificationInIdealDeliveryWindow(notification *datamodel.Notification, nonIdealDeliveryTime *time.Time) bool {
	if notification == nil ||
		notification.IdealDevlieryConditions == nil ||
		nonIdealDeliveryTime == nil {
		return false
	}

	now := timeNow()

	// Check current time is in the ideal delivery time.
	// Must be after deliveryTime, but before dt+offset.
	if nonIdealDeliveryTime.After(now) {
		return false
	}
	if now.Sub(*nonIdealDeliveryTime) > notification.IdealDevlieryConditions.MaxWaitTime() {
		return false
	}

	// Check TOD and DOW filters work for current time
	if !timeMeetsFilterConditions(notification, &now) {
		return false
	}

	// All checks passed, delivery time is in ideal window
	return true
}

func timeMeetsFilterConditions(notification *datamodel.Notification, t *time.Time) bool {
	if !slices.Contains(notification.DeliveryDaysOfWeek, t.Weekday()) {
		return false
	}
	mintueOfDay := t.Hour()*60 + t.Minute()
	if mintueOfDay < notification.DeliveryWindowTODStartMinutes ||
		mintueOfDay > notification.DeliveryWindowTODEndMinutes {
		return false
	}

	return true
}

// If we're in the ideal delivery window, and condition passes: now is new ideal delivery time
// If we're in the ideal delivery window, and condition fails: delay delivery until end of ideal window
func (ac *Appcore) shiftDeliveryTimeForIdealWindow(notification *datamodel.Notification, nonIdealDeliveryTime *time.Time) *time.Time {
	if nonIdealDeliveryTime == nil ||
		notification == nil {
		return nil
	}

	// Check if now is in ideal delivery window, and if the condition passes
	inIdealDeliveryWindow := notificationInIdealDeliveryWindow(notification, nonIdealDeliveryTime)
	if inIdealDeliveryWindow {
		idealConditionResult, err := ac.propertyRegistry.evaluateCondition(&notification.IdealDevlieryConditions.Condition)
		if idealConditionResult && err == nil {
			now := timeNow()
			return &now
		}
	}

	// Shift delivery time back to end of offset (or nil it out for offset=forever)
	if notification.IdealDevlieryConditions != nil {
		if notification.IdealDevlieryConditions.WaitForever() {
			return nil
		} else {
			endOfIdealDeliveryWindow := nonIdealDeliveryTime.Add(notification.IdealDevlieryConditions.MaxWaitTime())
			return &endOfIdealDeliveryWindow
		}
	}

	// No ideal time, so return non ideal time
	return nonIdealDeliveryTime
}

// Base delivery time for notification based on static delivery time and event time, ignoring ideal time and delivery window filters
func (ac *Appcore) baseDeliveryTimeForNotification(notification *datamodel.Notification) *time.Time {
	if canceled := ac.isNotificationCanceled(notification); canceled {
		return nil
	}
	if notification.ScheduleCondition != nil {
		condResult, condErr := ac.propertyRegistry.evaluateCondition(notification.ScheduleCondition)
		if !condResult || condErr != nil {
			return nil
		}
	}

	if staticTimestamp := notification.DeliveryTime.Timestamp(); staticTimestamp != nil {
		// Statically scheduled
		// If time has passed, we should not schedule static time notification
		if time.Now().After(*staticTimestamp) {
			return nil
		}
		return staticTimestamp
	} else if eventName := notification.DeliveryTime.EventName; eventName != nil {
		// Event based scheduling
		eventTime, err := eventTimeForDeliveryTime(ac, &notification.DeliveryTime)
		if eventTime == nil || err != nil {
			return nil
		}
		deliveryTime := eventTime
		if notification.DeliveryTime.EventOffset != nil {
			offsetTime := eventTime.Add(time.Duration(*notification.DeliveryTime.EventOffset) * time.Second)
			deliveryTime = &offsetTime
		}

		return deliveryTime
	}

	return nil
}

// Move the time forward until it is in the delivery window filters (time of day, day of week)
func shiftDeliveryTimeForFilters(notification *datamodel.Notification, deliveryTime *time.Time) *time.Time {
	if deliveryTime == nil || notification == nil {
		return nil
	}
	if timeMeetsFilterConditions(notification, deliveryTime) {
		return deliveryTime
	}

	newTime := *deliveryTime

	// Shift hours first, if needed
	deliveryMinuteOfDay := newTime.Minute() + newTime.Hour()*60
	startHour := notification.DeliveryWindowTODStartMinutes / 60
	startMinute := notification.DeliveryWindowTODStartMinutes % 60
	if deliveryMinuteOfDay < notification.DeliveryWindowTODStartMinutes {
		// Shift to start time on same day
		newTime = time.Date(newTime.Year(), newTime.Month(), newTime.Day(), startHour, startMinute, 0, 0, newTime.Location())
	} else if deliveryMinuteOfDay > notification.DeliveryWindowTODEndMinutes {
		// Shift to next day at start time (we never shift backwards)
		newTime = time.Date(newTime.Year(), newTime.Month(), newTime.Day()+1, startHour, startMinute, 0, 0, newTime.Location())
	}

	// Shift day of week forward 1d until it's in the window
	for int := 0; int < 7; int++ {
		if slices.Contains(notification.DeliveryDaysOfWeek, newTime.Weekday()) {
			break
		}
		newTime = newTime.AddDate(0, 0, 1)
	}

	return &newTime
}

func (ac *Appcore) isNotificationCanceled(notification *datamodel.Notification) bool {
	if notification.CancelationEvents == nil {
		return false
	}
	uncachedCancelationEvents := make([]string, 0)
	// Check cache first.
	// If any canceled=true are found, we can return true without checking DB
	for _, cancelEventName := range *notification.CancelationEvents {
		canceled := ac.seenCancelationEvents[cancelEventName]
		if canceled != nil && *canceled {
			return true
		} else if canceled == nil {
			uncachedCancelationEvents = append(uncachedCancelationEvents, cancelEventName)
		}
	}
	// Check DB for uncached cancelation events
	for _, cancelEventName := range uncachedCancelationEvents {
		cancelEventCount, err := ac.db.EventCountByNameWithLimit(cancelEventName, 1)
		if err != nil {
			// shouldn't fail, but better to be safe
			return true
		}
		canceled := cancelEventCount > 0
		ac.seenCancelationEvents[cancelEventName] = &canceled
		if canceled {
			return true
		}
	}
	return false
}

func eventTimeForDeliveryTime(ac *Appcore, dt *datamodel.DeliveryTime) (*time.Time, error) {
	if dt.EventInstance() == datamodel.EventInstanceTypeLatest {
		return ac.db.LatestEventTimeByName(*dt.EventName)
	} else if dt.EventInstance() == datamodel.EventInstanceTypeFirst {
		return ac.db.FirstEventTimeByName(*dt.EventName)
	}
	return nil, errors.New("invalid event instance type")
}

func (ac *Appcore) notificationRunnerProcessEvent(event *datamodel.Event) error {
	ac.updateCancelationEventCache(event)

	needsUpdate, err := ac.notificationsNeedUpdateForEvent(event)
	if err != nil {
		return err
	}
	if needsUpdate {
		err = ac.ForceUpdateNotificationPlan()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ac *Appcore) updateCancelationEventCache(event *datamodel.Event) {
	// Check if already cached
	cacheValue := ac.seenCancelationEvents[event.Name]
	if cacheValue != nil && *cacheValue {
		return
	}

	// Check if this is a cancelation event
	for _, notif := range ac.config.Notifications {
		if notif.CancelationEvents != nil {
			for _, cancelEventName := range *notif.CancelationEvents {
				if event.Name == cancelEventName {
					canceled := true
					ac.seenCancelationEvents[event.Name] = &canceled
					return
				}
			}
		}
	}
}

func (ac *Appcore) notificationsNeedUpdateForEvent(event *datamodel.Event) (bool, error) {
	// Don't update before we are initialized to prevent several runs on startup
	if ac.notificationPlan == nil {
		return false, nil
	}

	// Check if this event cancels existing scheduled notification
	for _, sns := range ac.notificationPlan.scheduledNotifications {
		sn := sns.Notification
		if sn.CancelationEvents != nil {
			for _, cancelEventName := range *sn.CancelationEvents {
				if event.Name == cancelEventName {
					return true, nil
				}
			}
		}
	}

	// Need update if a notification is triggered by this event
	for _, notif := range ac.config.Notifications {
		if notif.DeliveryTime.EventName == nil || *notif.DeliveryTime.EventName != event.Name {
			// Not a trigger for this event
			continue
		}

		// Latest case: always update test plan
		if notif.DeliveryTime.EventInstance() == datamodel.EventInstanceTypeLatest {
			return true, nil
		}
		// First case: only update if this is the first event, and not already scheduled
		if notif.DeliveryTime.EventInstance() == datamodel.EventInstanceTypeFirst {
			alreadyScheduled := false
			for _, sns := range ac.notificationPlan.scheduledNotifications {
				if sns.Notification.ID == notif.ID {
					alreadyScheduled = true
					break
				}
			}
			if !alreadyScheduled {
				return true, nil
			}
		}
	}

	return false, nil
}
