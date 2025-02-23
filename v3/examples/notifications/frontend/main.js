import * as Notifications from "./bindings/github.com/wailsapp/wails/v3/pkg/services/notifications/service";
import { Events } from "@wailsio/runtime";

const timeElement = document.getElementById('time');
const notificationsElement = document.getElementById('notifications');

window.sendNotification = async () => {
    const granted = await Notifications.RequestUserNotificationAuthorization();
    if (granted) {
        await Notifications.SendNotification("some-uuid-fronted", "Frontend Notificaiton", "", "Notificaiton sent through JS!");
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
                { id: "DELETE_ACTION", title: "Delete", destructive: true },
            ],
            hasReplyField: true,
            replyButtonTitle: "Reply",
            replyPlaceholder: "Reply to frontend...",
        });
        
        await Notifications.SendNotificationWithActions({
            id: "some-uuid-complex",
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

Events.On('time', (time) => {
    timeElement.innerText = time.data;
});

Events.On("notificationResponse", (response) => {
    notificationsElement.innerText += JSON.stringify(response.data[0].data);
});