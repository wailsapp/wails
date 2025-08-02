import { Events } from "@wailsio/runtime";
import { NotificationService } from "../bindings/github.com/wailsapp/wails/v3/pkg/services/notifications";

const footer = document.querySelector("#response");

document.querySelector("#request")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.RequestNotificationAuthorization();
        if (authorized) {
            if (footer) footer.innerHTML = "<p>Notifications are now authorized.</p>";
            console.info("Notifications are now authorized.");
        } else {
            if (footer) footer.innerHTML = "<p>Notifications are not authorized. You can attempt to request again or let the user know in the UI.</p>";
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});

document.querySelector("#check")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.CheckNotificationAuthorization();
        if (authorized) {
            if (footer) footer.innerHTML = "<p>Notifications are authorized.</p>";
            console.info("Notifications are authorized.");
        } else {
            if (footer) footer.innerHTML = "<p>Notifications are not authorized. You can attempt to request again or let the user know in the UI.</p>";
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});

document.querySelector("#basic")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.CheckNotificationAuthorization();
        if (authorized) {
            await NotificationService.SendNotification({
                id: crypto.randomUUID(),
                title: "Notification Title",
                subtitle: "Subtitle on macOS and Linux",
                body: "Body text of notification.",
                data: {
                    "user-id":    "user-123",
					"message-id": "msg-123",
					"timestamp":  Date.now(),
                },
            });
        } else {
            if (footer) footer.innerHTML = "<p>Notifications are not authorized. You can attempt to request again or let the user know in the UI.</p>";
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});
document.querySelector("#complex")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.CheckNotificationAuthorization();
        if (authorized) {
            const CategoryID = "frontend-notification-id";

            await NotificationService.RegisterNotificationCategory({
                id: CategoryID,
                actions: [
                    { id: "VIEW", title: "View" },
                    { id: "MARK_READ", title: "Mark as read" },
                    { id: "DELETE", title: "Delete", destructive: true },
                ],
				hasReplyField:    true,
				replyPlaceholder: "Message...",
				replyButtonTitle: "Reply",
            });

            await NotificationService.SendNotificationWithActions({
                id: crypto.randomUUID(),
                title: "Notification Title",
                subtitle: "Subtitle on macOS and Linux",
                body: "Body text of notification.",
                categoryId: CategoryID,
                data: {
                    "user-id":    "user-123",
					"message-id": "msg-123",
					"timestamp":  Date.now(),
                },
            });
        } else {
            if (footer) footer.innerHTML = "<p>Notifications are not authorized. You can attempt to request again or let the user know in the UI.</p>";
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});

const unlisten = Events.On("notification:action", (response) => {
    console.info(`Recieved a ${response.name} event`);
    const { userInfo, ...base } = response.data[0]; 
    console.info("Notification Response:");
    console.table(base);
    console.info("Notification Response Metadata:");
    console.table(userInfo);
    const table = `
        <h5>Notification Response</h5>
        <table>
            <thead>
                ${Object.keys(base).map(key => `<th>${key}</th>`).join("")}
            </thead>
            <tbody>
                ${Object.values(base).map(value => `<td>${value}</td>`).join("")}
            </tbody>
        </table>
        <h5>Notification Metadata</h5>
        <table>
            <thead>
                ${Object.keys(userInfo).map(key => `<th>${key}</th>`).join("")}
            </thead>
            <tbody>
                ${Object.values(userInfo).map(value => `<td>${value}</td>`).join("")}
            </tbody>
        </table>
    `;
    const footer = document.querySelector("#response");
    if (footer) footer.innerHTML = table;
});

window.onbeforeunload = () => unlisten();