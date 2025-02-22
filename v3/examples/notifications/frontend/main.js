import { GreetService } from "./bindings/notifications";
import { NotificationService } from "./bindings/github.com/wailsapp/wails/v3/pkg/services/notification";
import { Events } from "@wailsio/runtime";

const resultElement = document.getElementById('result');
const timeElement = document.getElementById('time');
const notificationsElement = document.getElementById('notifications');

const nofitications = new Map();

window.doGreet = () => {
    let name = document.getElementById('name').value;
    if (!name) {
        name = 'anonymous';
    }
    GreetService.Greet(name).then((result) => {
        resultElement.innerText = result;
    }).catch((err) => {
        console.log(err);
    });
}

window.sendNotification = async () => {
    const granted = await NotificationService.CheckNotificationAuthorization();
    if (granted) {
        await NotificationService.SendNotification("some-uuid-fronted", "Frontend Notificaiton", "", "Notificaiton sent through JS!");
    }
}

window.sendComplexNotification = async () => {
    const granted = await NotificationService.CheckNotificationAuthorization();
    if (granted) {
        await NotificationService.RegisterNotificationCategory({
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
        
        await NotificationService.SendNotificationWithActions({
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

Events.On("notification", (data) => {
    nofitications.set(data.identifier, data);
    notificationsElement.innerText = Array.from(nofitications.values()).map(notification => JSON.stringify(notification)).join(", ");
});