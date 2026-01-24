import {Call} from "@wailsio/runtime";

const mainView = document.getElementById('mainView');
const childView = document.getElementById('childView');

const params = new URLSearchParams(window.location.search);
const isChild = params.get('child') === '1';
const childName = params.get('name') || '';

if (mainView && childView) {
    mainView.style.display = isChild ? 'none' : '';
    childView.style.display = isChild ? '' : 'none';
}

const windowsListEl = document.getElementById('windowsList');
const lastChildNameEl = document.getElementById('lastChildName');

window.refreshWindowList = async () => {
    try {
        const windows = await Call.ByName('main.GreetService.ListWindows');
        if (windowsListEl) {
            windowsListEl.textContent = JSON.stringify(windows, null, 2);
        }
    } catch (err) {
        console.error(err);
        if (windowsListEl) {
            windowsListEl.textContent = String(err);
        }
    }
};

window.openChildWindow = async () => {
    try {
        const name = await Call.ByName('main.GreetService.OpenChildWindow');
        if (lastChildNameEl) {
            lastChildNameEl.textContent = name || '-';
        }
        await window.refreshWindowList();
    } catch (err) {
        console.error(err);
        if (lastChildNameEl) {
            lastChildNameEl.textContent = String(err);
        }
    }
};

// Child window helpers
const childNameEl = document.getElementById('childName');
const childLogEl = document.getElementById('childLog');
if (isChild && childNameEl) {
    childNameEl.textContent = childName || '(missing ?name=...)';
}

function setChildLog(obj) {
    if (!childLogEl) return;
    if (typeof obj === 'string') {
        childLogEl.textContent = obj;
        return;
    }
    try {
        childLogEl.textContent = JSON.stringify(obj, null, 2);
    } catch (e) {
        childLogEl.textContent = String(obj);
    }
}

window.closeNoCurrent = async () => {
    try {
        const ok = await Call.ByName('main.GreetService.CloseByName', childName);
        setChildLog({action: 'CloseByName', childName, ok});
    } catch (err) {
        console.error(err);
        setChildLog({action: 'CloseByName', childName, error: String(err)});
    }
};

window.closeAfterCurrent = async () => {
    try {
        const ok = await Call.ByName('main.GreetService.CloseAfterCurrentByName', childName);
        setChildLog({action: 'CloseAfterCurrentByName (calls App.Window.Current first)', childName, ok});
    } catch (err) {
        console.error(err);
        setChildLog({action: 'CloseAfterCurrentByName', childName, error: String(err)});
    }
};

window.closeUsingCurrent = async () => {
    try {
        const info = await Call.ByName('main.GreetService.CloseUsingCurrent');
        setChildLog({action: 'CloseUsingCurrent (closes App.Window.Current())', info});
    } catch (err) {
        console.error(err);
        setChildLog({action: 'CloseUsingCurrent', error: String(err)});
    }
};

// Nice default: keep the window list updated on the main window.
if (!isChild) {
    window.refreshWindowList();
}
