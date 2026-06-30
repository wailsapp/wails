import { Events } from "@wailsio/runtime";
import {
    NotificationService,
} from "../bindings/github.com/wailsapp/wails/v3/pkg/services/notifications";
import * as DemoAssetsService from "../bindings/notifications/demoassets";

// Inline the option type so we don't have to fight tsconfig path resolution
// for the generated d.ts files; the generated bindings already accept this
// shape via $models.NotificationOptions.
type Sound = { silent?: boolean; name?: string };
type Attachment = { id?: string; path: string; type?: string };
type Schedule = { delaySeconds?: number; at?: number };
type Notif = {
    id: string;
    title: string;
    subtitle?: string;
    body?: string;
    categoryId?: string;
    data?: Record<string, unknown>;
    sound?: Sound | null;
    attachments?: Attachment[];
    threadId?: string;
    interruptionLevel?: string;
    schedule?: Schedule | null;
};

const footer = document.querySelector("#response") as HTMLElement | null;
const statusEl = document.querySelector("#status") as HTMLElement | null;

const TEST_CATEGORY_ID = "demo-actions";
let lastSentID: string | null = null;
let lastScheduledID: string | null = null;

// SampleImagePath rewrites the bundled sample image to its TempDir path on
// every call — macOS UNNotificationAttachment moves the source file out
// after each delivery, so caching the path between sends would leave a
// stale reference that fails the os.Stat check on the Go side.
function getSampleImagePath(): Promise<string> {
    return DemoAssetsService.SampleImagePath();
}

async function ensureCategory(): Promise<void> {
    await NotificationService.RegisterNotificationCategory({
        id: TEST_CATEGORY_ID,
        actions: [
            { id: "VIEW", title: "View" },
            { id: "MARK_READ", title: "Mark as read" },
            { id: "DELETE", title: "Delete", destructive: true },
        ],
        hasReplyField: true,
        replyPlaceholder: "Message...",
        replyButtonTitle: "Reply",
    });
}

async function send(notif: Notif): Promise<void> {
    lastSentID = notif.id;
    if (notif.schedule) {
        lastScheduledID = notif.id;
    }
    if (notif.categoryId) {
        await NotificationService.SendNotificationWithActions(notif as any);
    } else {
        await NotificationService.SendNotification(notif as any);
    }
}

async function update(notif: Notif): Promise<void> {
    lastSentID = notif.id;
    await NotificationService.UpdateNotification(notif as any);
}

function status(message: string): void {
    if (statusEl) statusEl.innerHTML = `<p>${message}</p>`;
    console.info(message);
}

async function ensureAuthorized(): Promise<boolean> {
    const authorized = await NotificationService.CheckNotificationAuthorization();
    if (!authorized) {
        status(
            "Notifications are not authorized. Click <em>Request Authorization</em> first (macOS only — Windows and Linux always return true).",
        );
    }
    return authorized;
}

// --- Authorization buttons ----------------------------------------------------

document.querySelector("#request")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.RequestNotificationAuthorization();
        status(authorized ? "Notifications are now authorized." : "Authorization denied.");
    } catch (error) {
        console.error(error);
        status(`Authorization request failed: ${error}`);
    }
});

document.querySelector("#check")?.addEventListener("click", async () => {
    try {
        const authorized = await NotificationService.CheckNotificationAuthorization();
        status(authorized ? "Notifications are authorized." : "Notifications are NOT authorized.");
    } catch (error) {
        console.error(error);
        status(`Authorization check failed: ${error}`);
    }
});

// --- Quick test buttons -------------------------------------------------------

document.querySelector("#basic")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    await send({
        id: crypto.randomUUID(),
        title: "Basic notification",
        subtitle: "Subtitle on macOS and Linux",
        body: "Plain body text. No actions, default sound.",
        data: { source: "basic" },
    });
    status("Basic notification sent.");
});

document.querySelector("#complex")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    await ensureCategory();
    await send({
        id: crypto.randomUUID(),
        title: "Complex notification",
        subtitle: "With actions + reply",
        body: "Click an action button or type a reply.",
        categoryId: TEST_CATEGORY_ID,
        data: { source: "complex" },
    });
    status("Complex notification sent. Try the action buttons or reply field.");
});

document.querySelector("#thread")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    const threadId = "demo-thread-" + Math.floor(Math.random() * 1000);
    await send({
        id: crypto.randomUUID(),
        title: "Threaded #1",
        body: `First message in thread ${threadId}.`,
        threadId,
        data: { source: "thread", thread: threadId },
    });
    await new Promise((r) => setTimeout(r, 500));
    await send({
        id: crypto.randomUUID(),
        title: "Threaded #2",
        body: `Second message in thread ${threadId}. macOS groups these in Notification Center.`,
        threadId,
        data: { source: "thread", thread: threadId },
    });
    status(`Two notifications sent with threadId="${threadId}".`);
});

document.querySelector("#schedule")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    const id = crypto.randomUUID();
    await send({
        id,
        title: "Scheduled in 5s",
        body: "If you see this, the schedule path works on this platform.",
        schedule: { delaySeconds: 5 },
        data: { source: "schedule" },
    });
    status(
        `Scheduled notification id=${id}. macOS persists across app restart; Windows/Linux use an in-process timer.`,
    );
});

document.querySelector("#update")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    const id = crypto.randomUUID();
    await send({
        id,
        title: "Original title",
        body: "Will be updated in 2 seconds...",
        data: { source: "update" },
    });
    status(`Sent id=${id}. Updating in 2s...`);
    await new Promise((r) => setTimeout(r, 2000));
    await update({
        id,
        title: "Updated title",
        body: "macOS replaces in place; Linux uses replaces_id; Windows redelivers.",
        data: { source: "update" },
    });
    status(`Updated id=${id}.`);
});

document.querySelector("#cancel")?.addEventListener("click", async () => {
    if (!lastScheduledID) {
        status("No scheduled notification to cancel.");
        return;
    }
    await NotificationService.RemovePendingNotification(lastScheduledID);
    status(`Cancelled scheduled id=${lastScheduledID}.`);
    lastScheduledID = null;
});

// --- Builder form -------------------------------------------------------------

function val(id: string): string {
    return (document.querySelector("#" + id) as HTMLInputElement | HTMLSelectElement | null)?.value ?? "";
}
function checked(id: string): boolean {
    return (document.querySelector("#" + id) as HTMLInputElement | null)?.checked ?? false;
}

async function buildFromForm(id?: string): Promise<Notif> {
    const notif: Notif = {
        id: id ?? crypto.randomUUID(),
        title: val("b-title") || "(no title)",
        body: val("b-body"),
        subtitle: val("b-subtitle") || undefined,
        threadId: val("b-thread") || undefined,
        interruptionLevel: val("b-level") || undefined,
        data: { source: "builder" },
    };

    const soundChoice = val("b-sound");
    if (soundChoice === "silent") {
        notif.sound = { silent: true };
    } else if (soundChoice === "named") {
        const name = val("b-sound-name").trim();
        if (name) notif.sound = { name };
    }

    const delay = Number(val("b-delay"));
    if (Number.isFinite(delay) && delay > 0) {
        notif.schedule = { delaySeconds: Math.floor(delay) };
    }

    if (checked("b-attach")) {
        const path = await getSampleImagePath();
        // No `type` here so macOS auto-infers the UTI from the extension.
        // Windows defaults to inline placement; pass type:"hero" or
        // type:"appLogoOverride" if you want a different placement.
        notif.attachments = [{ path }];
    }

    if (checked("b-actions")) {
        await ensureCategory();
        notif.categoryId = TEST_CATEGORY_ID;
    }

    return notif;
}

document.querySelector("#b-send")?.addEventListener("click", async () => {
    if (!(await ensureAuthorized())) return;
    try {
        const notif = await buildFromForm();
        await send(notif);
        status(`Sent id=${notif.id}.`);
    } catch (error) {
        console.error(error);
        status(`Send failed: ${error}`);
    }
});

document.querySelector("#b-update-by-id")?.addEventListener("click", async () => {
    if (!lastSentID) {
        status("No previous notification to update — send one first.");
        return;
    }
    if (!(await ensureAuthorized())) return;
    try {
        const notif = await buildFromForm(lastSentID);
        await update(notif);
        status(`Updated id=${notif.id}.`);
    } catch (error) {
        console.error(error);
        status(`Update failed: ${error}`);
    }
});

document.querySelector("#b-remove")?.addEventListener("click", async () => {
    if (!lastSentID) {
        status("No previous notification to remove.");
        return;
    }
    try {
        await NotificationService.RemoveNotification(lastSentID);
        await NotificationService.RemovePendingNotification(lastSentID);
        await NotificationService.RemoveDeliveredNotification(lastSentID);
        status(`Remove called for id=${lastSentID} (no-op on platforms that don't track delivered toasts).`);
    } catch (error) {
        console.error(error);
        status(`Remove failed: ${error}`);
    }
});

// --- Action / reply response handler -----------------------------------------

const unlisten = Events.On("notification:action", (response) => {
    console.info(`Received a ${response.name} event`);
    // The current @wailsio/runtime passes the emitted payload directly as
    // event.data; older versions wrapped it in [data]. Handle both shapes.
    const payload =
        response.data && typeof response.data === "object" && "id" in response.data
            ? response.data
            : Array.isArray(response.data)
                ? response.data[0]
                : null;
    if (!payload) {
        console.warn("notification:action received with empty payload", response);
        return;
    }
    const { userInfo, ...base } = payload;
    console.info("Notification Response:");
    console.table(base);
    if (userInfo) {
        console.info("Notification Response Metadata:");
        console.table(userInfo);
    }

    // Build tables via DOM API (not innerHTML) so notification payloads
    // containing HTML/JS cannot execute in the demo page.
    function buildTable(data: Record<string, unknown>): HTMLTableElement {
        const table = document.createElement("table");
        const thead = table.createTHead();
        const tbody = table.createTBody();
        const headerRow = thead.insertRow();
        const dataRow = tbody.insertRow();
        for (const [key, value] of Object.entries(data)) {
            const th = document.createElement("th");
            th.textContent = key;
            headerRow.appendChild(th);
            const td = dataRow.insertCell();
            td.textContent = String(value);
        }
        return table;
    }

    if (footer) {
        footer.textContent = "";
        const h5 = document.createElement("h5");
        h5.textContent = "Notification Response";
        footer.appendChild(h5);
        footer.appendChild(buildTable(base));
        if (userInfo) {
            const h5meta = document.createElement("h5");
            h5meta.textContent = "Notification Metadata";
            footer.appendChild(h5meta);
            footer.appendChild(buildTable(userInfo as Record<string, unknown>));
        }
    }
});

window.onbeforeunload = () => unlisten();
