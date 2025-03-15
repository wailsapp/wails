import * as Notifications from "./bindings/github.com/wailsapp/wails/v3/pkg/services/notifications";
import { Events } from "@wailsio/runtime";

const notificationsElement = document.getElementById('notifications');

window.sendNotification = async () => {
    const granted = await Notifications.Service.RequestNotificationAuthorization();
    if (granted) {
        await Notifications.Service.SendNotification({
            id: crypto.randomUUID(),
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
    const granted = await Notifications.Service.RequestNotificationAuthorization();
    if (granted) {
        await Notifications.Service.RegisterNotificationCategory({
            id: "FRONTEND_NOTIF",
            actions: [
                { id: "VIEW_ACTION", title: "View" },
                { id: "MARK_READ_ACTION", title: "Mark as Read" },
            ],
            hasReplyField: true,
            replyButtonTitle: "Reply",
            replyPlaceholder: "Reply to frontend...",
        });

        await Notifications.Service.SendNotificationWithActions({
            id: crypto.randomUUID(),
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

Events.On("notification:response", (response) => {
    notificationsElement.innerText = JSON.stringify(response.data[0]);
});