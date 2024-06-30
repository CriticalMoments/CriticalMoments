package appcore

import (
	"errors"
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

func (ac *Appcore) SendNotficationPlanToLib() error {
	plan, err := ac.NotificationPlan()
	if err != nil {
		return err
	}
	err = ac.libBindings.UpdateNotificationPlan(&plan)
	if err != nil {
		return err
	}

	return nil
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
	if staticTimestamp := notification.DeliveryTime.Timestamp(); staticTimestamp != nil {
		// Statically scheduled
		// If time has passed, we should not deliver static time notification
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

		// TODO_P0: in past do we still schedule? I think so but confirm
		return deliveryTime
	}

	return nil
}

func eventTimeForDeliveryTime(ac *Appcore, dt *datamodel.DeliveryTime) (*time.Time, error) {
	if dt.EventInstance() == datamodel.EventInstanceTypeLatest {
		return ac.db.LatestEventTimeByName(*dt.EventName)
	} else if dt.EventInstance() == datamodel.EventInstanceTypeFirst {
		return ac.db.FirstEventTimeByName(*dt.EventName)
	}
	return nil, errors.New("invalid event instance type")
}
