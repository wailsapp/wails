//go:build linux

package notifications

import (
	"testing"

	"github.com/godbus/dbus/v5"
)

func TestDismissedNotificationIsAClickOnlyWithoutActions(t *testing.T) {
	notificationServiceLock.Lock()
	prev := NotificationService_
	NotificationService_ = &NotificationService{}
	notificationServiceLock.Unlock()
	t.Cleanup(func() {
		notificationServiceLock.Lock()
		NotificationService_ = prev
		notificationServiceLock.Unlock()
	})
	var got []NotificationResult
	getNotificationService().OnNotificationResponse(func(r NotificationResult) { got = append(got, r) })

	for _, actions := range []bool{true, false} {
		got = nil
		ln := &linuxNotifier{actions: actions, notifications: map[uint32]*notificationData{7: {ID: "n1"}}}
		// Reason 2: dismissed by the user.
		ln.handleNotificationClosed(&dbus.Signal{Body: []interface{}{uint32(7), uint32(2)}})
		if len(ln.notifications) != 0 {
			t.Fatalf("actions=%v: the closed notification is still tracked", actions)
		}
		if actions && len(got) != 0 {
			t.Fatalf("actions=%v: a dismissal was reported as %+v", actions, got)
		}
		if !actions && (len(got) != 1 || got[0].Response.ActionIdentifier != DefaultActionIdentifier || got[0].Response.ID != "n1") {
			t.Fatalf("actions=%v: want one default-action response for n1, got %+v", actions, got)
		}
	}
}
