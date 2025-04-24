import { Events } from "@wailsio/runtime";
import * as Notifications from "../bindings/github.com/wailsapp/wails/v3/pkg/services/notifications";

document.querySelector("#basic")?.addEventListener("click", async () => {
    try {
        const authorized = await Notifications.Service.CheckNotificationAuthorization();
        if (authorized) {
            await Notifications.Service.SendNotification({
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
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});
document.querySelector("#complex")?.addEventListener("click", async () => {
    try {
        const authorized = await Notifications.Service.CheckNotificationAuthorization();
        if (authorized) {
            const CategoryID = "frontend-notification-id";

            await Notifications.Service.RegisterNotificationCategory({
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

            await Notifications.Service.SendNotificationWithActions({
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
            console.warn("Notifications are not authorized.\n You can attempt to request again or let the user know in the UI.\n");
        }
    } catch (error) {
        console.error(error);
    }
});

const unlisten = Events.On("notification:action", (response) => {
    console.info(`Received a ${response.name} event`);
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