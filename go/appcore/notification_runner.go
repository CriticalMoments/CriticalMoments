package appcore

import (
	"time"

	datamodel "github.com/CriticalMoments/CriticalMoments/go/cmcore/data_model"
)

type NotificationPlan struct {
	unscheduledNotifications    []*datamodel.Notification
	scheduledNotifications      []*ScheduledNotification
	RequiresBackgroundExecution bool
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

func (ac *Appcore) NotificationPlan() (NotificationPlan, error) {
	plan := NotificationPlan{
		unscheduledNotifications: make([]*datamodel.Notification, 0),
		scheduledNotifications:   make([]*ScheduledNotification, 0),
	}
	for _, notification := range ac.config.Notifications {
		if deliveryTimestamp := ac.deliveryTimeForNotification(notification); deliveryTimestamp != nil {
			sn := ScheduledNotification{
				Notification: notification,
				scheduledAt:  *deliveryTimestamp,
			}
			plan.scheduledNotifications = append(plan.scheduledNotifications, &sn)
		} else {
			plan.unscheduledNotifications = append(plan.unscheduledNotifications, notification)
		}
	}
	// TODO_P0: filter those already delivered
	// TODO_P0: set BG time needed somewhere

	return plan, nil
}

func (ac *Appcore) deliveryTimeForNotification(notification *datamodel.Notification) *time.Time {
	// TODO_P0: if ideal, move to end of delviery window

	// Statically scheduled
	if staticTimestamp := notification.DeliveryTime.Timestamp(); staticTimestamp != nil {
		if time.Now().After(*staticTimestamp) {
			// Time has passed, we should not deliver this notification
			return nil
		}
		return staticTimestamp
	}

	// Event based scheduling
	if eventName := notification.DeliveryTime.EventName; eventName != nil {
		// TODO_P0: this is always the latest. Need modes here I think.
		lastEventTime, err := ac.db.LatestEventTimeByName(*eventName)
		if lastEventTime == nil || err != nil {
			// Contunue with the next notification
			return nil
		}
		targetTime := lastEventTime
		if notification.DeliveryTime.EventOffset != nil {
			offsetTime := lastEventTime.Add(time.Duration(*notification.DeliveryTime.EventOffset) * time.Second)
			targetTime = &offsetTime
		}
		return targetTime
	}

	return nil
}
