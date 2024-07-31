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

	// The earliest time to run background work for notifications
	earliestBgCheckTimeEpochSeconds int64
}

// Expalainin the achitecture a bit here for notifications. It's a bit tricky due to restructions of iOS APIs.
// Main thing to know: iOS doesn't store delivered notifications. It might have them, but it might not if they are cleared once the user dismisses them. So the pattern is "if they are scheduled, we assume they are delivered".
// General idea: golang delivers truth from DB. It consistently can reproduce the same plan with identical delivery times (for a given now() time). iOS layer takes care of the rest.
// CM iOS layer: deals with the golang plan. This including knowing if it already delviered a notification for a given timestamp (some assumtions of scheduled get delivered).
// To handle lack of concrete "delivered", and async messages, we use a trick. iOS code pushes out delivery 1s; this gives us time to propagate cancelations or reschedules before it is delivered to the user. However go may say "that should have already fired and shouldn't fire twice" and keep the original time (in which case that's also okay).
// iOS also makes the assumtion that notifications with old delivery times were already delivered. Since go is calling iOS to schedule as soon as the DB entries are set, this should be true.

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
	now := time.Now()
	return ac.generateNotificationPlanForTime(now)
}

func (ac *Appcore) generateNotificationPlanForTime(now time.Time) (NotificationPlan, error) {
	plan := NotificationPlan{
		unscheduledNotifications: make([]*datamodel.Notification, 0),
		scheduledNotifications:   make([]*ScheduledNotification, 0),
	}

	var earliestBgCheckTime *time.Time

	for _, notification := range ac.config.Notifications {
		deliveryTimestamp, bgCheckTime := ac.notificationDeliveryTime(notification, now)
		if deliveryTimestamp != nil {
			sn := ScheduledNotification{
				Notification: notification,
				scheduledAt:  *deliveryTimestamp,
			}
			plan.scheduledNotifications = append(plan.scheduledNotifications, &sn)
		} else {
			plan.unscheduledNotifications = append(plan.unscheduledNotifications, notification)
		}

		if bgCheckTime != nil {
			if earliestBgCheckTime == nil || earliestBgCheckTime.After(*bgCheckTime) {
				earliestBgCheckTime = bgCheckTime
			}
		}
	}

	if earliestBgCheckTime != nil {
		plan.earliestBgCheckTimeEpochSeconds = earliestBgCheckTime.Unix()
	}

	return plan, nil
}

// Get the delivery time of a notification
// 1) First check when it should be delivered (static time, event based)
// 2) Then consider ideal delivery window, delivering sooner or later if we have special targeting in mind
// 3) Then consider the allowed time of day, and days of week for delivery
func (ac *Appcore) notificationDeliveryTime(notification *datamodel.Notification, now time.Time) (deliveryTime *time.Time, bgCheckTime *time.Time) {
	nonIdealDeliveryTime := ac.baseDeliveryTimeForNotification(notification, now)
	idealDeliveryTime, bgCheckTime := ac.shiftDeliveryTimeForIdealWindow(notification, nonIdealDeliveryTime, now)
	shiftedDeliveryTime := shiftDeliveryTimeForFilters(notification, idealDeliveryTime)
	return shiftedDeliveryTime, bgCheckTime
}

// Checks if this notification has an ideal delivery window and now is currently in the time-range of that window
func notificationInIdealDeliveryWindow(notification *datamodel.Notification, nonIdealDeliveryTime *time.Time, now time.Time) bool {
	if notification == nil ||
		notification.IdealDevlieryConditions == nil ||
		nonIdealDeliveryTime == nil {
		return false
	}

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
// Also: check what time we should schedule background checks for this notification's ideal delivery window, which meet the ideal delivery time (offset and filters)
func (ac *Appcore) shiftDeliveryTimeForIdealWindow(notification *datamodel.Notification, nonIdealDeliveryTime *time.Time, now time.Time) (shiftedTime *time.Time, checkTime *time.Time) {
	if nonIdealDeliveryTime == nil ||
		notification == nil {
		return nil, nil
	}

	// No ideal time window, so return non ideal time, and nil checkTime
	if notification.IdealDevlieryConditions == nil {
		return nonIdealDeliveryTime, nil
	}

	// Check if now is in ideal delivery window, and if the condition passes
	inIdealDeliveryWindow := notificationInIdealDeliveryWindow(notification, nonIdealDeliveryTime, now)
	if inIdealDeliveryWindow {
		idealConditionResult, err := ac.propertyRegistry.evaluateCondition(&notification.IdealDevlieryConditions.Condition)
		if idealConditionResult && err == nil {
			// No need for checkTime, since the condition is currently met
			return &now, nil
		}
	}

	// Shift delivery time back to end of offset, or nil it out for offset=forever
	var shiftedDeliveryTime *time.Time

	if notification.IdealDevlieryConditions.WaitForever() {
		shiftedDeliveryTime = nil
	} else {
		endOfIdealDeliveryWindow := nonIdealDeliveryTime.Add(notification.IdealDevlieryConditions.MaxWaitTime())
		shiftedDeliveryTime = &endOfIdealDeliveryWindow
	}

	// Build a checkTime: the time to run background check for this notification's ideal delivery window
	bgCheckTime := bgCheckTimeForIdealDeliveryWindow(notification, now, shiftedDeliveryTime)

	return shiftedDeliveryTime, bgCheckTime
}

const checkTimeDelay = 15 * time.Minute
const filterTimeBuffer = 2 * time.Minute

// Build a bgCheckTime: the time to run background check for this notification's ideal delivery window
func bgCheckTimeForIdealDeliveryWindow(notification *datamodel.Notification, now time.Time, shiftDeliveryTime *time.Time) *time.Time {
	if notification == nil {
		return nil
	}

	// Run at earliest 15 mins from now (too often uses quota), and first time meeting filters after that
	var checkTime = now.Add(checkTimeDelay)
	filterShiftedCheckTime := shiftDeliveryTimeForFilters(notification, &checkTime)
	if filterShiftedCheckTime != nil && filterShiftedCheckTime.After(checkTime) {
		// Add small buffer time after the first possible time passing the filter. Too close, and it might fire seconds before the filter time, garunteeing a failure.
		checkTime = filterShiftedCheckTime.Add(filterTimeBuffer)
	}

	// We don't want to run background after the shiftDeliveryTime, as we'll have already delivered by then
	if shiftDeliveryTime != nil && checkTime.After(*shiftDeliveryTime) {
		return nil
	}

	return &checkTime
}

// Base delivery time for notification based on static delivery time and event time, ignoring ideal time and delivery window filters
func (ac *Appcore) baseDeliveryTimeForNotification(notification *datamodel.Notification, now time.Time) *time.Time {
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
		if now.After(*staticTimestamp) {
			return nil
		}
		return staticTimestamp
	} else if eventName := notification.DeliveryTime.EventName; eventName != nil {
		// Event based scheduling
		deliveryTime, err := deliveryTimeFromDB(ac, &notification.DeliveryTime)
		if deliveryTime == nil || err != nil {
			return nil
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

func deliveryTimeFromDB(ac *Appcore, dt *datamodel.DeliveryTime) (*time.Time, error) {
	var t *time.Time
	var err error

	offset := dt.EventOffsetDuration()

	// Latest Once becomes first when there's zero offset. Use first logic since it's more efficient
	isFirst := dt.EventInstance() == datamodel.EventInstanceTypeFirst || (offset == 0 && dt.EventInstance() == datamodel.EventInstanceTypeLatestOnce)

	if dt.EventInstance() == datamodel.EventInstanceTypeLatest {
		t, err = ac.db.LatestEventTimeByName(*dt.EventName)
		if err != nil {
			return nil, err
		}
	} else if isFirst {
		t, err = ac.db.FirstEventTimeByName(*dt.EventName)
		if err != nil {
			return nil, err
		}
	} else if dt.EventInstance() == datamodel.EventInstanceTypeLatestOnce {
		t, err = latestOnceEventTimeFromDB(ac, dt)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("invalid event instance type")
	}

	// no event is valid, not an error case
	if t == nil {
		return nil, nil
	}

	// Apply offset
	offsetTime := t.Add(offset)
	return &offsetTime, nil
}

// "event time" and not "delivery time", the caller must apply the offset
func latestOnceEventTimeFromDB(ac *Appcore, dt *datamodel.DeliveryTime) (*time.Time, error) {
	eventTimes, err := ac.db.AllEventTimesByName(*dt.EventName)
	if err != nil {
		return nil, err
	}

	return latestOnceEventTimeFromEventList(dt, eventTimes)
}

func latestOnceEventTimeFromEventList(dt *datamodel.DeliveryTime, eventTimes []time.Time) (*time.Time, error) {
	// Latest once strategy: iterate through the events. The delivery time is the first event that is followed by a gap of at least the offset (+offset).
	// This will consistently return the same delivery time from DB state without additional DB state.
	if len(eventTimes) == 0 {
		return nil, nil
	}

	offset := dt.EventOffsetDuration()

	lastTime := eventTimes[0]
	for i, eventTime := range eventTimes {
		if i == 0 {
			continue
		}
		lastScheduledTime := lastTime.Add(offset)
		if eventTime.After(lastScheduledTime) {
			// The notification should have been delivered at lastScheduledTime
			break
		}
		lastTime = eventTime
	}

	return &lastTime, nil
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

		// Latest case and LatestOnce case: always update test plan
		if notif.DeliveryTime.EventInstance() == datamodel.EventInstanceTypeLatest ||
			notif.DeliveryTime.EventInstance() == datamodel.EventInstanceTypeLatestOnce {
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

func (ac *Appcore) performBackgroundWorkForNotifications() error {
	// TODO_P0: optimize this? Can check if any notifications are in ideal window and not update if not needed.
	return ac.ForceUpdateNotificationPlan()
}
