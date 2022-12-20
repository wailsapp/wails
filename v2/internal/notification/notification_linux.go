//go:build linux

package notification

import (
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/wailsapp/wails/v2/internal/frontend"
	"github.com/wailsapp/wails/v2/internal/logger"
)

const (
	dbusObjectPath             = "/org/freedesktop/Notifications"
	dbusNotificationsInterface = "org.freedesktop.Notifications"
	signalNotificationClosed   = "org.freedesktop.Notifications.NotificationClosed"
	signalActionInvoked        = "org.freedesktop.Notifications.ActionInvoked"
	callGetCapabilities        = "org.freedesktop.Notifications.GetCapabilities"
	callCloseNotification      = "org.freedesktop.Notifications.CloseNotification"

	MethodNotifySend = "notify-send"
	MethodDbus       = "dbus"
	MethodKdialog    = "kdialog"

	notifyChannelBufferSize = 25
)

type closedReason uint32

func (r closedReason) string() string {
	switch r {
	case 1:
		return "expired"
	case 2:
		return "dismissed-by-user"
	case 3:
		return "closed-by-call"
	case 4:
		return "unknown"
	case 5:
		return "activated-by-user"
	default:
		return "other"
	}
}

type Notifier struct {
	*sync.Mutex
	logger       *logger.Logger
	method       string
	dbusConn     *dbus.Conn
	sendPath     string
	intitialized bool
}

// NewNotifier returns a new Notifier
func NewNotifier(myLogger *logger.Logger) *Notifier {
	return &Notifier{
		Mutex:  &sync.Mutex{},
		logger: myLogger,
	}
}

func (n *Notifier) init() error {
	var err error

	checkDbus := func() (*dbus.Conn, error) {
		conn, err := dbus.SessionBusPrivate()
		if err != nil {
			return conn, err
		}

		if err = conn.Auth(nil); err != nil {
			return conn, err
		}

		if err = conn.Hello(); err != nil {
			return conn, err
		}

		obj := conn.Object(dbusNotificationsInterface, dbusObjectPath)
		call := obj.Call(callGetCapabilities, 0)
		if call.Err != nil {
			return conn, call.Err
		}

		var ret []string
		err = call.Store(&ret)
		if err != nil {
			return conn, err
		}

		// add a listener (matcher) in dbus for signals to Notification interface.
		err = conn.AddMatchSignal(
			dbus.WithMatchObjectPath(dbusObjectPath),
			dbus.WithMatchInterface(dbusNotificationsInterface),
		)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	n.dbusConn, err = checkDbus()
	if err == nil {
		n.method = MethodDbus
		return nil
	}
	n.dbusConn.Close()
	n.dbusConn = nil

	send, err := exec.LookPath("notify-send")
	if err == nil {
		n.sendPath = send
		n.method = MethodNotifySend
		return nil
	}

	send, err = exec.LookPath("sw-notify-send")
	if err == nil {
		n.sendPath = send
		n.method = MethodNotifySend
		return nil
	}

	n.method = "none"
	n.sendPath = ""

	return err
}

var initOnce sync.Once

// SendNotification sends notifications
func (n *Notifier) SendNotification(options frontend.NotificationOptions) error {
	n.Lock()
	defer n.Unlock()

	var (
		ID  uint32
		err error
	)

	initOnce.Do(func() {
		err = n.init()
	})

	if err != nil {
		return errors.New("notification: could not initialize notifications")
	}

	switch n.method {
	case MethodDbus:
		ID, err = n.sendViaDbus(options)
	case MethodNotifySend:
		ID, err = n.sendViaNotifySend(options)
	case MethodKdialog:
		ID, err = n.sendViaKnotify(options)
	default:
		err = errors.New("no notification method is available")
	}

	if err == nil && options.LinuxOptions.OnShow != nil {
		options.LinuxOptions.OnShow(ID)
	}

	return err
}

func (n *Notifier) sendViaDbus(options frontend.NotificationOptions) (result uint32, err error) {

	var actions []string
	if options.LinuxOptions.Actions != nil {
		for i := range options.LinuxOptions.Actions {
			actions = append(actions, options.LinuxOptions.Actions[i].Key, options.LinuxOptions.Actions[i].Label)
		}
	}

	timeout := int32(-1)
	if options.Timeout > 0 {
		timeout = int32(options.Timeout.Milliseconds())
	}

	hints := map[string]dbus.Variant{}

	hints["urgency"] = dbus.MakeVariant(frontend.LinuxNotificationUrgency(options.LinuxOptions.Urgency).Uint())

	if options.LinuxOptions.Sound != nil {
		if options.LinuxOptions.Sound.File != nil {
			s, err := SoundPath(options.LinuxOptions.Sound.File)
			if err != nil {
				n.logger.Error("Notification sound error: %v")
			} else {
				hints["sound-file"] = dbus.MakeVariant(s)
			}
		} else if options.LinuxOptions.Sound.Name != "" {
			hints["sound-name"] = dbus.MakeVariant(options.LinuxOptions.Sound.Name)
		}
		if options.LinuxOptions.Sound.Suppress {
			hints["suppress-sound"] = dbus.MakeVariant(true)
		}
	}

	appIcon := ""
	if options.AppIcon != nil {
		appIcon, err = AppIconPath(options.AppIcon)
		if err != nil {
			n.logger.Error("Notification app icon err: %v", err)
		}
	}

	obj := n.dbusConn.Object(dbusNotificationsInterface, dbusObjectPath)
	dbusArgs := []interface{}{
		options.AppID,
		options.LinuxOptions.ReplacesID,
		appIcon,
		options.Title,
		options.Message,
		actions,
		hints,
		timeout,
	}

	call := obj.Call("org.freedesktop.Notifications.Notify", 0, dbusArgs...)
	if call.Err != nil {
		return 0, fmt.Errorf("runtime dbus notification err: %v", err)
	}

	err = call.Store(&result)
	if err != nil {
		return 0, err
	}

	if result != options.LinuxOptions.ReplacesID && (len(actions) > 0 || options.LinuxOptions.OnClose != nil) {
		go dbusListener(n.dbusConn, n.logger, result, options)
	}
	return
}

func (n *Notifier) sendViaNotifySend(options frontend.NotificationOptions) (uint32, error) {

	args := []string{
		options.Title,
		options.Message,
	}

	if options.AppIcon != nil {
		if appIcon, err := AppIconPath(options.AppIcon); err == nil {
			args = append(args, fmt.Sprintf("--icon=%s", appIcon))
		}
	}

	timeout := int32(-1)
	if options.Timeout > 0 {
		timeout = int32(options.Timeout.Milliseconds())
		args = append(args, fmt.Sprintf("--expire-time=%d", timeout))
	}

	args = append(args, "--urgency=%s", frontend.LinuxNotificationUrgency(options.LinuxOptions.Urgency).String())

	if options.LinuxOptions.Sound != nil {
		if options.LinuxOptions.Sound.File != nil {
			s, err := SoundPath(options.LinuxOptions.Sound.File)
			if err != nil {
				n.logger.Error("Notification sound error: %v")
			} else {
				args = append(args, fmt.Sprintf("--hint=string:sound-file:%s", s))
			}
		} else if options.LinuxOptions.Sound.Name != "" {
			args = append(args, fmt.Sprintf("--hint=string:sound-name:%s", options.LinuxOptions.Sound.Name))
		}
		if options.LinuxOptions.Sound.Suppress {
			args = append(args, fmt.Sprintf("--hint=string:suppress-sound:%v", options.LinuxOptions.Sound.Suppress))
		}
	}

	c := exec.Command(n.sendPath, args...)
	err := c.Run()
	if err != nil {
		return 0, fmt.Errorf("runtime notify-send notification err: %v", err)
	}
	return 0, nil
}

func (n *Notifier) sendViaKnotify(options frontend.NotificationOptions) (uint32, error) {

	c := exec.Command(n.sendPath, "--title", options.Title, "--passivepopup", options.Message, "10", "--icon", "")
	err := c.Run()
	if err != nil {
		return 0, fmt.Errorf("runtime knotify notification err: %v", err)
	}
	return 0, nil
}

func dbusListener(conn *dbus.Conn, logger *logger.Logger, notificationID uint32, options frontend.NotificationOptions) {
	// register in dbus for signal delivery
	signal := make(chan *dbus.Signal, notifyChannelBufferSize)
	conn.Signal(signal)

	var (
		signalID        uint32
		signalActionKey string
	)

	for {
		select {
		case c := <-options.Close:
			if c {
				obj := conn.Object(dbusNotificationsInterface, dbusObjectPath)
				call := obj.Call(callCloseNotification, 0, notificationID)
				if call.Err != nil {
					logger.Error("Notification cancel error %v", call.Err)
				}
			}
		case s := <-signal:
			if s == nil {
				logger.Error("Notification signal error: empty signal received")
				return
			}
			if len(s.Body) < 2 {
				logger.Error("Notification signal error: incomplete signal received")
			}
			switch s.Name {
			case signalNotificationClosed:
				signalID = s.Body[0].(uint32)
				if options.LinuxOptions.OnClose != nil && notificationID == signalID {
					options.LinuxOptions.OnClose(signalID, closedReason(s.Body[1].(uint32)).string())
					return
				}
			case signalActionInvoked:
				signalID = s.Body[0].(uint32)
				signalActionKey = s.Body[1].(string)
				if notificationID != signalID {
					continue
				}
				for i := 0; i < len(options.LinuxOptions.Actions); i++ {
					if options.LinuxOptions.Actions[i].OnAction != nil && options.LinuxOptions.Actions[i].Key == signalActionKey {
						options.LinuxOptions.Actions[i].OnAction(signalID)
						if options.LinuxOptions.OnClose != nil {
							options.LinuxOptions.OnClose(signalID, closedReason(5).string())
						}
						return
					}
				}
			}
		}
	}
}
