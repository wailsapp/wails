
export function UniqueID(text) {
    return text + "-" + Date.now().toString() + Math.random().toString();
}