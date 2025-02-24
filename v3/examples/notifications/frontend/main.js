import * as Notifications from "./bindings/github.com/wailsapp/wails/v3/pkg/services/notifications/service";
import { Events, System } from "@wailsio/runtime";

const notificationsElement = document.getElementById('notifications');

window.sendNotification = async () => {
    const granted = await Notifications.RequestUserNotificationAuthorization();
    if (granted) {
        const id = System.IsWindows() ? "Wails Notification Demo" : crypto.randomUUID()
        await Notifications.SendNotification({
            id,
            title: "Title",
            body: "Body!",
            data: {
                messageId: "msg-123",
                senderId: "user-123",
                timestamp: Date.now(),
            }
        });
    }
}

window.sendComplexNotification = async () => {
    const granted = await Notifications.RequestUserNotificationAuthorization();
    if (granted) {
        await Notifications.RegisterNotificationCategory({
            id: "FRONTEND_NOTIF",
            actions: [
                { id: "VIEW_ACTION", title: "View" },
                { id: "MARK_READ_ACTION", title: "Mark as Read" },
            ],
            hasReplyField: true,
            replyButtonTitle: "Reply",
            replyPlaceholder: "Reply to frontend...",
        });
        
        const id = System.IsWindows() ? "Wails Notification Demo" : crypto.randomUUID()
        await Notifications.SendNotificationWithActions({
            id,
            title: "Complex Frontend Notification",
            subtitle: "From: Jane Doe",
            body: "Is it rainging today where you are?",
            categoryId: "FRONTEND_NOTIF",
            data: {
                messageId: "msg-456",
                senderId: "user-456",
                timestamp: Date.now(),
            }
        });
    }
}

Events.On("notificationResponse", (response) => {
    console.log(response)
    notificationsElement.innerText += JSON.stringify(response.data[0].data);
});